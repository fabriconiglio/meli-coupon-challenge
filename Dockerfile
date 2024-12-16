# Build stage
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copiar archivos de dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/api/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copiar el binario compilado
COPY --from=builder /app/main .

# Exponer el puerto
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./main"]