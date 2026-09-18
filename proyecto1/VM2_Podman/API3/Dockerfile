# Etapa de compilación
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Etapa final ligera
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]