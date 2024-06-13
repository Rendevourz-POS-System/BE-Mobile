# Builder stage
FROM golang:alpine as builder

# Install git. Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Build the Go app
WORKDIR /app/src
RUN go build -o shelter-apps .

# Ensure your binary is executable (if necessary)
#RUN chmod +x /app/src/shelter-apps
CMD ["/shelter-apps"]

# Start a new stage from scratch
FROM alpine:latest

WORKDIR /app/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/src/shelter-apps /app/src/
COPY --from=builder /app/config /app/config/
COPY --from=builder /app/src/local.env /app/src/

# Expose port 8008 to the outside world
EXPOSE 8080

WORKDIR /app/src
# Command to run the executable
CMD ["./shelter-apps"]
