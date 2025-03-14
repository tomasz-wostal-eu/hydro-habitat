FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install necessary dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=darwin go build -a -installsuffix cgo -o api ./cmd/api

# Use a smaller image for the final build
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/api .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./api"]