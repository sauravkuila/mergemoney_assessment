FROM golang:1.24.7-alpine AS builder
WORKDIR /src

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the server binary (main is in cmd/api) and place it in /src/cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./cmd/api/server ./cmd/api

FROM alpine:3.18
RUN apk add --no-cache ca-certificates

# install psql client for health checks
RUN apk add --no-cache ca-certificates postgresql-client

# Copy the built app and all relevant files so runtime can find config files
COPY --from=builder /src /app
COPY --from=builder /src/tools/wait-for-db.sh /usr/local/bin/wait-for-db.sh

RUN chmod +x /usr/local/bin/wait-for-db.sh

WORKDIR /app/cmd/api

EXPOSE 8080

# Use the wait-for-db wrapper as entry; it will exec the server when DB is ready
ENTRYPOINT ["/usr/local/bin/wait-for-db.sh", "./server"]
