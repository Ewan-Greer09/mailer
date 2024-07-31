# Stage 1: Build the Go binary
FROM golang:1.22-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Cache go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app (assuming 'emailer' is the entry point)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /emailer ./cmd/emailer

# Stage 2: Create a small image with the built binary
FROM alpine:3.18

# Set up certificates (if your app needs to make HTTPS requests)
RUN apk --no-cache add ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /emailer .

# Expose the port the app runs on (if needed)
EXPOSE 8080

# Command to run the binary
CMD ["./emailer"]

