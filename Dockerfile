# Compile stage
FROM golang:1.18.1 AS builder
ARG DEBIAN_FRONTEND=noninteractive

# Add required packages
RUN apt -y update && apt -y install git curl bash

WORKDIR /go/src/GIG

ADD go.mod go.sum ./
RUN go mod download

RUN go get -u github.com/revel/revel
RUN go install github.com/revel/cmd/revel@latest
RUN go get -u github.com/lsflk/gig-sdk

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
ADD . .
RUN revel build "" build -m dev

# Run stage
FROM golang:1.18.1
EXPOSE 9000
ARG DEBIAN_FRONTEND=noninteractive
RUN apt -y update && useradd -m -s /bin/bash 10001

WORKDIR /app
COPY --from=builder /go/src/GIG/build .

RUN chmod +x /app/run.sh

# Create necessary directories
RUN mkdir -p /app/app/cache

# Set proper permissions for the app directory and its contents
RUN chown -R 10001:10001 /app
RUN chmod -R 755 /app
RUN chmod -R 777 /app/app/cache

# Switch to the non-root user
USER 10001

ENTRYPOINT ["/app/run.sh"]