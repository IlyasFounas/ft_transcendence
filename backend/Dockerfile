FROM golang:1.25-alpine

# Installs build essentials: git (for go mod), make, gcc and musl-dev (for CGO)
RUN --mount=type=cache,target=/var/cache/apk \
    apk add git make gcc musl-dev

RUN --mount=type=cache,target=/root/.cache/go-build \
    go install github.com/air-verse/air@latest

# Adds /go/bin to PATH to access Go-installed binaries like air
ENV PATH=$PATH:/go/bin

WORKDIR /app

COPY go.mod go.sum ./

# bcrypt
RUN --mount=type=cache,target=/root/.cache/go-build \
    go get golang.org/x/crypto/bcrypt

# Downloads Go modules dependencies to /go/pkg/mod/
RUN --mount=type=cache,target=/root/.cache/go-build \
    go mod download

RUN --mount=type=cache,target=/root/.cache/go-build \
    go get github.com/wneessen/go-mail

EXPOSE 8080

COPY . .