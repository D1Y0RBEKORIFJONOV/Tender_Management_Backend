FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main cmd/app/main.go

FROM alpine:3.18


WORKDIR /app

COPY --from=builder /app/main .

COPY auth.conf /app/auth.conf
COPY auth.csv /app/auth.csv


CMD ["./main"]
