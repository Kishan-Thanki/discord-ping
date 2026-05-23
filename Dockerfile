# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o discord-ping .

# Run stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates && \
    adduser -D -g '' botuser

WORKDIR /app
COPY --from=builder /app/discord-ping .

USER botuser
CMD ["./discord-ping"]
