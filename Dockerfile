# Start from a specific version of the golang base image
FROM golang:1.22.2 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Start a new stage from a specific version of the alpine base image
FROM alpine:3.15

# Add ca-certificates in case you need to make calls to HTTPS endpoints
RUN apk --no-cache add ca-certificates

# Create a new user to run the application
RUN adduser -D appuser
USER appuser

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# accept envrionment variables and run using it
ENV PORT=8080
ENV NAME=World
# Command to run the executable
CMD ["./main"]
