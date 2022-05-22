package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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
	lambda.Start(handler)
}
