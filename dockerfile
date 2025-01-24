# Stage 1: Build the Go application
FROM golang:alpine AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o bin ./cmd/api/main.go

# Stage 2: Run the Go application
FROM alpine:latest

# Install any required dependencies
RUN apk --no-cache add ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=build /app/bin .

# Copy the .env file
COPY .env .

EXPOSE 8080

CMD ["./bin"]