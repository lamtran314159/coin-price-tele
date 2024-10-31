# Use the official Golang image as a parent image
FROM golang:1.23.1-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy only necessary files to leverage Docker cache more effectively
COPY src/go.mod src/go.sum ./

# Download all dependencies first, leveraging Docker caching
RUN go mod download

# Copy the rest of the source code into the container
COPY src/ /app/src

# Set the working directory to where main.go is located
WORKDIR /app/src

# Build the application (output as "main" instead of "main.go")
RUN go build -o main .

# Expose the port that the application listens on
EXPOSE 8443

# Command to run the application
CMD ["./main"]
