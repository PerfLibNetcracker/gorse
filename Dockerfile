# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:1.15-alpine

# Add Maintainer Info
LABEL maintainer="PLYSHKA <leruop@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Create data folder for boltDB
RUN mkdir -p /app/bin/data

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# Expose port 8050 to the outside world
EXPOSE 8050

# Command to run the executable
CMD ["./main","serve","-c","bin/config/gorse_docker.toml"]