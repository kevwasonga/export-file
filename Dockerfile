# Use the official Go image (version 1.22.2) as the build environment
FROM golang:1.22.2-alpine As jose

# Set the working directory inside the Docker container
WORKDIR /dockerContainer

# Copy the entire project directory from the host to the container
COPY . .

# Build the Go project with the target OS set to Linux
# The resulting binary is named 'ascii'
RUN GOOS=linux go build -o ascii

# Start of the second stage, which creates a smaller image
# using Alpine Linux without Go (multi-stage build)
FROM alpine:latest

# Set the working directory inside the Docker container
WORKDIR /dockerContainer

# Copy the built 'ascii' binary from the first stage
COPY --from=jose /dockerContainer/ascii .

# Copy the 'banners' directory from the first stage to the new container
COPY --from=jose /dockerContainer/banners /dockerContainer/banners

# Copy the 'static' directory from the first stage to the new container
COPY --from=jose /dockerContainer/static /dockerContainer/static

# Specify the command to run the 'ascii' binary when the container starts
CMD ["/dockerContainer/ascii"]

# Add metadata labels to the Docker image
LABEL docker-version="27.1.1"
LABEL golang version="1.22.2"
LABEL Contributers="<joseowino> <kewasonga> <vomolo>"

# Docker permission setup (optional)
# Commands to install Docker in rootless mode
# curl -fsSL https://get.docker.com/rootless | sh
# export PATH=/home/docker/bin:$PATH
# export DOCKER_HOST=unix:///run/user/10531/docker.sock
