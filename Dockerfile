# syntax=docker/dockerfile:1

# TO BUILD:
# docker build --tag tiny-api:v1.0 . 
# TO RUN:
# docker run --publish 9001:9001 -e TINY_API_HOST=0.0.0.0 -e TINY_API_INSTANCE_NAME=instance1 -t tiny-api:v1.0

# image available at:
# https://hub.docker.com/repository/docker/stubin87/tiny-api

FROM golang:1.19.1-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg

RUN go build -o /bin/api ./cmd/api/main.go

EXPOSE 9001

CMD [ "/bin/api" ]
