ARG GO_VERSION=1.14.4

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*

RUN mkdir -p /app
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o ./app ./main.go

FROM jrottenberg/ffmpeg:4.1-alpine

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

RUN mkdir -p /app
RUN mkdir -p /app/logs
WORKDIR /app
COPY --from=builder /app/app .

ENTRYPOINT ["./app"]
