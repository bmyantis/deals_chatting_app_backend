# Use the official Golang image as a base image
FROM golang:1.22.4

# Set the current working directory inside the container
WORKDIR /app

# Install Air globally
RUN go install github.com/cosmtrek/air@v1.27.3

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Install dependencies
RUN go mod tidy

# Expose port 8080 to the outside world
EXPOSE 8080


# Command to run Air
CMD ["air"]
