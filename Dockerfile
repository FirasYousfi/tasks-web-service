# Start from the latest golang base image
FROM golang:1.19.3-alpine

# Add Maintainer Info
LABEL maintainer="Firas Yousfi <firas.yousfi144@gmail.com>"

# The latest alpine images don't have some tools like (`git` and `bash`).
# Adding git, bash and openssh to the image.
# Git is required to fetch the dependences
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh && \
    apk add curl

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy everything from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app, you need to give it the path to the executable
RUN go build -o main cmd/main.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]