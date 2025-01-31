
# Stage 1: Build
FROM golang:1.23-alpine AS builder

# Install build dependencies in a single command
RUN apk add --no-cache gcc musl-dev libwebp-dev

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./

# Download all dependencies and clean up after
RUN go mod tidy && go mod download

# Copy all source code to the working directory
COPY . .

# Build the Golang application
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /htmlcsstoimage

# Stage 2: Package
FROM alpine:latest

# Install minimal dependencies in one command
RUN apk add --no-cache libwebp

# Set up the working directory in the new container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /htmlcsstoimage /htmlcsstoimage

# Copy only necessary files (like .env) instead of all files
COPY .env .env

# Ensure the storage/images directory exists
RUN mkdir -p storage/images/

# Expose the port used by the application
EXPOSE 8080

# Command to run the application
CMD ["/htmlcsstoimage"]
