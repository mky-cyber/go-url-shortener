ARG GO_VERSION=1
ARG ALPINE_VERSION=3.17

# Stage one - build the binary
FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod tidy && go mod verify

COPY . .
RUN GOOS=linux go build -v -ldflags="-s -w" -o /gourlshortener ./cmd

# Stage two - deploy the binary and the SQLite DB
FROM alpine:${ALPINE_VERSION}

WORKDIR /opt

# Copy over the database migration
COPY ./db db

# Install sqlite
RUN apk --no-cache add sqlite

# Copy over the Go binary and set it as the command to run on boot
COPY --from=builder /gourlshortener /usr/local/bin/gourlshortener
COPY ./db/migrations/202405191609_create_urls_table.sql /docker-entrypoint-initdb.d/202405191609_create_urls_table.sql

# Ensure the database path is defined
ENV DATABASE_PATH="/opt/database.db"

# Run the SQL script to initialize the DB if it doesn't exist
RUN sqlite3 /opt/database.db < /docker-entrypoint-initdb.d/202405191609_create_urls_table.sql

ENTRYPOINT ["/bin/sh", "-c", "/usr/local/bin/gourlshortener"]