# Start from a base image containing the Go runtime 
FROM golang:latest

# Set the current working directory inside the container
WORKDIR /go/src/app

# Copy the local package directory into the container's workspace
COPY . .

# Download dependencies using go mod
RUN go mod download

# Build the Go application
RUN go build -o main .

# Expose port 8082 to the outside world
EXPOSE 8082

# Command to run the executable
CMD ["./main"]
