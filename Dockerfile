# Use the official Go 1.21 image as a parent image
FROM golang:1.21

# Set the working directory inside the container
WORKDIR /app

# Copy the local source files to the container's working directory
COPY . .

# Build the Go application
RUN go build -o main .

# Command to run the application
CMD ["./main"]