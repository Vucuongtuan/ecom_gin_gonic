# Use the official Golang image
FROM golang:latest

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main ./src/main.go

# Expose port 8080
EXPOSE 8080

# Run the executable
CMD ["./main"]
