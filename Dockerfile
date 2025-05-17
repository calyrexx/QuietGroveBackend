FROM golang:1.24.2-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /quietgrove ./cmd/main.go

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /quietgrove .
COPY configuration.yaml credentials.yaml ./
COPY deploy/postgres.sql ./deploy/postgres.sql

EXPOSE 8080

CMD ["./quietgrove"]