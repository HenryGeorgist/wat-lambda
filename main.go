package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/batch"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/usace/wat-api/utils"
	"github.com/usace/wat-api/wat"
	"gopkg.in/yaml.v2"
)

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {
		task := wat.Task{} //needs to be a Task from batch branch.
		err := yaml.Unmarshal([]byte(string(message.Body)), &task)
		if err != nil {
			fmt.Println("error while parsing message body", err)
		}
		fmt.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)
		fmt.Println("sending task to", task.TaskType)
	}

	return nil
}

func main() {
	fmt.Println("initializing a wat-lambda")
	loader, err := utils.InitLoader("")
	if err != nil {
		log.Fatal(err)
		return
	}
	queue, err := loader.InitQueue()
	if err != nil {
		log.Fatal(err)
		return
	}
	awsBatch, err := loader.InitBatch()
	if err != nil {
		log.Fatal(err)
		return
	}
	lambda.Start(func(ctx context.Context, sqsEvent events.SQSEvent) error {
		for _, message := range sqsEvent.Records {
			task := wat.Task{} //needs to be a Task from batch branch.
			err := yaml.Unmarshal([]byte(string(message.Body)), &task)
			if err != nil {
				fmt.Println("error while parsing message body", err)
			}
			fmt.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)
			fmt.Println("sending task to", task.TaskType)
			if task.TaskType == "Lambda" {
				fmt.Println("sending message: ", task.ModelPayload)
				blob, err := yaml.Marshal(task.ModelPayload)
				payload := string(blob)
				queueName := "lambda-tasks" //sending to a different sqs queue (probably needs to be by plugin type?)
				input := sqs.GetQueueUrlInput{
					QueueName: &queueName,
				}
				queueURL, err := queue.GetQueueUrl(&input)
				if err != nil {
					return err
				}
				fmt.Println("sending message to:", queueURL.QueueUrl)
				_, err = queue.SendMessage(&sqs.SendMessageInput{
					DelaySeconds: aws.Int64(1),
					MessageBody:  aws.String(payload),
					QueueUrl:     queueURL.QueueUrl,
				})
				fmt.Println("message sent")
			} else {
				//default to batch.
				//send task to batch
				path := "some path in s3, need to send payload to s3"
				proptags := true
				batchOutput, err := awsBatch.SubmitJob(&batch.SubmitJobInput{
					//DependsOn: dependsOn,
					ContainerOverrides: &batch.ContainerOverrides{
						Command: []*string{
							aws.String(".\\main -payload=" + path),
						},
					},
					//JobDefinition:              resources[idx].JobARN, //need to verify this.
					//JobName:                    &key,
					//JobQueue:                   resources[idx].QueueARN,
					Parameters:                 nil,       //parameters?
					PropagateTags:              &proptags, //i think.
					RetryStrategy:              nil,
					SchedulingPriorityOverride: nil,
					ShareIdentifier:            nil,
					Tags:                       nil,
					Timeout:                    nil,
				})
				fmt.Println("batchoutput", batchOutput)
				if err != nil {
					fmt.Println("batcherror", err)
					panic(err)
				}
			}
		}

		return nil
	})
}
