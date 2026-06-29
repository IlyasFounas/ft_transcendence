FROM golang:1.25-alpine

RUN go install github.com/air-verse/air@latest

ENV PATH=$PATH:/go/bin

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .