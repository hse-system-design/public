FROM golang:1.16-alpine AS builder

WORKDIR /build

COPY ./main.go .
RUN go mod init lockmod && go mod tidy

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o locker .

ENTRYPOINT ["/build/locker"]