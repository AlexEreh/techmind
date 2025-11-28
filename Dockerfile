FROM golang:1.24.4-alpine3.22 AS builder

RUN apk update && apk add --no-cache ca-certificates

WORKDIR /opt/build

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy application source.
COPY . .

# Build the application.
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o /opt/bin/application ./cmd/main.go

# Prepare executor image.
FROM alpine:3.21 AS runner

RUN apk update && apk add --no-cache tzdata ca-certificates bash && \
    cp /usr/share/zoneinfo/Europe/Moscow /etc/localtime && \
    echo "Europe/Moscow" > /etc/timezone && \
    apk del tzdata && \
    rm -rf /var/cache/apk/*

ENV TZ=/etc/localtime

WORKDIR /app

COPY ./migrations ./migrations
COPY --from=builder /opt/bin/application ./

CMD ["./application"]
