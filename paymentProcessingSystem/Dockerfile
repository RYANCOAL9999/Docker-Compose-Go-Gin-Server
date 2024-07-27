# Use an official Go runtime as a parent image
FROM golang:1.20 as builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# download dependencies
RUN GOSUMDB=off go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOPROXY=direct GOSUMDB=off go build -o paymentProcessingSystem .

# Final stage
FROM alpine:3.12 as production

# Set working directory in the final image
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/paymentProcessingSystem .

# Copy the .env.example to .env in the final image
COPY .env.example .env

# Ensure the .env file is used by the application
ENV PORT=${PORT}

# Run the application
CMD ["sh", "-c", "./paymentProcessingSystem"]