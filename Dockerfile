FROM golang:1.18-alpine3.14 as dev

ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV GO111MODULE="on"
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin
RUN apk add --no-cache git
ENV CGO_ENABLED 0 


# Production container
FROM golang:1.18-alpine3.14 AS prod
RUN apk add --update docker openrc
RUN rc-update add docker boot
WORKDIR /app
COPY --from=dev /app/main .
CMD [ "./main" ]