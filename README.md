# ascii-art-web-dockerize

This project aims to Dockerize the ascii-art-web application. The application, initially developed to provide a web GUI for generating ASCII art, is now containerized using Docker to enhance its deployment, scalability, and ease of use.

## Project Overview

The ascii-art-web-dockerize project involves creating a Docker image and container for the ascii-art-web application. The web server is implemented in Go, following best practices for both Go programming and Docker containerization.

## Objectives
### Dockerize the Application:
```
* Create a Dockerfile.
* Build a Docker image.
* Run the application within a Docker container.
```
### Apply Metadata:
```
* Add relevant metadata to Docker objects to ensure clarity and maintainability.
* Garbage Collection:
* Implement practices to clean up unused Docker objects to save space and resources.
```

## Features
### Web Server in Go:

The web server is implemented using the Go programming language.
The server adheres to best practices for performance, security, and code quality.

### Docker Integration:
The project utilizes Docker to create a portable and easily deployable application.
A Dockerfile is provided, adhering to Docker's best practices.

## Usage
### Prerequisites

Docker: Ensure Docker is installed on your system. You can find installation instructions on the Docker documentation.

### Getting Started

1. Clone the Repository:
```
git clone https://learn.zone01kisumu.ke/git/joseowino/ascii-art-web-dockerize
cd ascii-art-web-dockerize
```
2. Run the script:
```
./ascii.sh
```

## Metadata and Garbage Collection

Metadata: Docker objects (images, containers, etc.) include metadata for easy identification and management.
Garbage Collection: Unused Docker objects are periodically cleaned up to prevent resource wastage.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Authors

[kevwasonga](https://learn.zone01kisumu.ke/git/kevwasonga)

[vomolo](https://learn.zone01kisumu.ke/git/vomolo)

[joseowino](https://learn.zone01kisumu.ke/git/joseowino)
