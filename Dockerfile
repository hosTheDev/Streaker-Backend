FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o myapp

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/myapp .
COPY --from=builder /app/.env . 

EXPOSE 8080

CMD ["./myapp"]