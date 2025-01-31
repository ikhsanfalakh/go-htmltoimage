# Gunakan image resmi Golang sebagai base image
FROM golang:1.23-alpine

# Set working directory di dalam container
WORKDIR /app

# Copy go.mod dan go.sum untuk mengunduh dependencies
COPY go.mod go.sum ./

# Unduh semua dependencies
RUN go mod tidy
RUN go mod download

# Copy seluruh kode sumber ke working directory
COPY . .

# Build aplikasi Golang
#RUN go build -o main .
#RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/htmlcsstoimage
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main

# Expose port yang digunakan oleh aplikasi
EXPOSE 8080

# Command untuk menjalankan aplikasi
CMD ["./main"]