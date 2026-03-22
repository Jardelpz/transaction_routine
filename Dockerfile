FROM golang:1.25.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY .env .env

EXPOSE 8080

CMD ["./main"]