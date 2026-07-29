FROM golang:1.25-alpine AS base
# Adds /go/bin to PATH to access Go-installed binaries like air
ENV PATH=$PATH:/go/bin
WORKDIR /app
COPY go.mod go.sum ./
COPY . .


FROM base AS dev
# Installs build essentials: git (for go mod), make, gcc and musl-dev (for CGO)
RUN --mount=type=cache,target=/var/cache/apk \
    apk add git make gcc musl-dev
RUN --mount=type=cache,target=/root/.cache/go-build \
    go install github.com/air-verse/air@v1.52.3
# bcrypt
RUN --mount=type=cache,target=/root/.cache/go-build \
    go get golang.org/x/crypto/bcrypt
# Downloads Go modules dependencies to /go/pkg/mod/
RUN --mount=type=cache,target=/root/.cache/go-build \
    go mod download
RUN --mount=type=cache,target=/root/.cache/go-build \
    go get github.com/wneessen/go-mail
EXPOSE 8080
CMD ["sh", "-c", "go mod tidy && air"]

FROM base AS build
RUN go get golang.org/x/crypto/bcrypt
RUN go mod download
RUN go get github.com/wneessen/go-mail
RUN go run 

FROM build AS prod
EXPOSE 8080
CMD ["sh", "-c", "go mod tidy && air"]
