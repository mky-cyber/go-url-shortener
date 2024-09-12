ARG GO_VERSION=1
ARG ALPINE_VERSION=3.17

# Stage one - build the binary
FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod tidy && go mod verify

COPY . .
RUN GOOS=linux go build -v -ldflags="-s -w" -o /gourlshortener ./cmd

# Stage two - deploy the binary
FROM alpine:${ALPINE_VERSION}

WORKDIR /opt

# Copy over the Go binary and set it as the command to run on boot
COPY --from=builder /gourlshortener /usr/local/bin/gourlshortener
ENTRYPOINT ["/bin/sh", "-c", "/usr/local/bin/gourlshortener"]