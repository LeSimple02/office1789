# Build stage
FROM golang:1.23-alpine AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server .

# Run stage: alpine + docker CLI for docker socket interactions
FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache docker-cli ca-certificates bash

COPY --from=build /app/server .

ENV GIN_MODE=release
EXPOSE 8080

CMD ["./server"]
