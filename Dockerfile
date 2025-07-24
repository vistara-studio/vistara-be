FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/vistara-backend cmd/api/main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates curl && update-ca-certificates

WORKDIR /app

COPY --from=builder /app/bin/vistara-backend .

ARG APP_PORT=8081
EXPOSE ${APP_PORT}

ENTRYPOINT ["./vistara-backend"]