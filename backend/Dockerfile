FROM golang:1.25-alpine AS base
# Adds /go/bin to PATH to access Go-installed binaries like air
ENV PATH=$PATH:/go/bin
WORKDIR /app
COPY go.mod go.sum ./
COPY entrypoint.sh ./
COPY . .

FROM base AS dev
# Installs build essentials: git (for go mod), make, gcc and musl-dev (for CGO)
RUN --mount=type=cache,target=/var/cache/apk \
    apk add git make gcc musl-dev
RUN --mount=type=cache,target=/root/.cache/go-build \
    go install github.com/air-verse/air@v1.52.3
# Downloads Go modules dependencies to /go/pkg/mod/
RUN --mount=type=cache,target=/root/.cache/go-build \
    go mod download
EXPOSE 8080
CMD ["sh", "-c", "air"]

FROM base AS build
RUN go mod download
RUN go build -o /golang_api
RUN go build -o /bun_migrator ./cmd/bun

FROM golang:1.25-alpine AS prod
WORKDIR /app
EXPOSE 8080
COPY --from=build /golang_api .
COPY --from=build /bun_migrator .
COPY --from=build /app/entrypoint.sh .
COPY --from=base /app/migrations ./migrations
RUN chmod +x /app/entrypoint.sh
ENTRYPOINT ["/app/entrypoint.sh"]