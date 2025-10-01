FROM golang:1.25.1-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o backend-main ./cmd/main.go

# ===================================

FROM alpine:latest

WORKDIR /app

COPY --from=builder /build/backend-main .
COPY --from=builder /build/config ./config

EXPOSE 8000

CMD ["./backend-main", "-c", "config/config.yaml"]
