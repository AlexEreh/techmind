FROM golang:1.25-alpine AS builder

WORKDIR /build

# Install dependencies
# RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

# Copy source code
COPY ../.. .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go

FROM alpine:latest

RUN apk update && apk add --no-cache tzdata ca-certificates bash && \
    cp /usr/share/zoneinfo/Europe/Moscow /etc/localtime && \
    echo "Europe/Moscow" > /etc/timezone && \
    apk del tzdata && \
    rm -rf /var/cache/apk/*

ENV TZ=/etc/localtime

WORKDIR /app

# Copy the binary from builder

COPY ./migrations ./migrations
COPY --from=builder /build/app .

CMD ["./app"]

