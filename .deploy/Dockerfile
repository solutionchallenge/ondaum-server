# ------------------------------
# 1) Builder
# ------------------------------
FROM golang:1.24.2-alpine AS builder
WORKDIR /app

RUN apk add --no-cache gcc musl-dev git
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .
RUN swag init
ENV CGO_ENABLED=1
RUN GOOS=linux go build -a -ldflags '-s -w -linkmode external -extldflags "-static"' -o /app/ondaum ./main.go

# ------------------------------
# 2) Runner
# ------------------------------
FROM alpine:3.21 AS runner
WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/ondaum .
COPY --from=builder /app/configs ./configs

EXPOSE 8080
CMD ["/app/ondaum", "http"]