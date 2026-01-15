FROM golang:1.24.2 AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o main ./cmd

FROM debian:bookworm-slim

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/configs/config.yaml . 

EXPOSE 8000

CMD ["./main"]