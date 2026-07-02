FROM golang:1.25-alpine

# Installs build essentials: git (for go mod), make, gcc and musl-dev (for CGO)
RUN apk add --no-cache git make gcc musl-dev

RUN go install github.com/air-verse/air@latest

# Adds /go/bin to PATH to access Go-installed binaries like air
ENV PATH=$PATH:/go/bin

WORKDIR /app

COPY go.mod go.sum ./

# Downloads Go modules dependencies to /go/pkg/mod/
RUN go mod download

EXPOSE 8080

COPY . .