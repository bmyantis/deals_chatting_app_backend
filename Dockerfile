# Stage 1: Build stage
FROM golang:latest AS build

# Set the current working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Install air globally
RUN go install github.com/cosmtrek/air@v1.27.3

# Copy the rest of the application source code
COPY . .

# Build the Go binary
RUN go build -o main .

# Stage 2: Run stage
FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=build /app/main .

# Copy the air executable from the build stage
COPY --from=build /go/bin/air /usr/local/bin/air

# Copy the rest of the application source code
COPY . .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["air"]
