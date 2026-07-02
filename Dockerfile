# Tahap 1: Membangun aplikasi (Builder)
FROM golang:1.22-alpine AS builder

# Install git jika diperlukan untuk download dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go.mod dan go.sum terlebih dahulu agar cache layer berfungsi
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Build aplikasi dengan output bernama 'main'
# Kita arahkan ke file main di dalam folder cmd/api
RUN go build -o main ./api/main.go

# Tahap 2: Menjalankan aplikasi (Runtime)
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Copy hasil build dari tahap pertama
COPY --from=builder /app/main .

# Ekspos port yang digunakan (asumsi port 8080)
EXPOSE 8080

# Jalankan aplikasinya
CMD ["./main"]