FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY .env /app/.env

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/

FROM scratch

COPY --from=builder /app/main /app/main
COPY --from=builder /app/.env /app/.env

WORKDIR /app

CMD ["./main"]
