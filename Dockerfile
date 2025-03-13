FROM golang:1.24.1-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev libpq-dev

COPY . .

RUN go mod tidy && CGO_ENABLED=1 GOOS=linux go build -o test-build-app main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/ .

EXPOSE 8080

CMD ["./test-build-app"]
# CMD ["sh", "-c", "sleep 20 && ./test-build-app"]