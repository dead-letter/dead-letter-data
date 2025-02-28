FROM golang:1.24.0-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN go build -ldflags "-w -s" -o main

ENTRYPOINT ["/app/main"]
