# go-url-shortener
A simple url shortener written in go.

## Introduction
go-url-shortener is a lightweight URL shortening service written in Go. It allows users to generate short URLs for long URLs and provides basic analytics (such as the number of clicks).

## Features
- Generate short URLs for long URLs
- Allow users to use the short URLs to redirect to original URL

## Installation

1. **Clone the repository:**
    ```bash
    git clone https://github.com/mky-cyber/go-url-shortener.git
    cd go-url-shortener
    ```

2. **Install dependencies:**
    ```bash
    go mod tidy
    ```

3. **Build the binary:**
    ```bash
    go build -o url-shortener ./cmd/main.go
    ```

4. **Run the project:**
    ```bash
    ./url-shortener
    ```

## Running the Project (Docker)

### To build the Docker image:

```bash
docker build -t url-shortener .
```
### To run the container:
```bash
docker run -d -p 8080:8080 --name go-url-shortener url-shortener
```
The service will be available at http://localhost:8080.


## API Documentation

This project provides a Postman collection to document and interact with the API.

### Using Postman to Explore the API

1. **Install Postman**  
   If you don't have Postman installed, you can download it from [here](https://www.postman.com/downloads/).

2. **Import the Postman Collection**  
   You can import the Postman collection into your Postman application by following these steps:
   - Download the Postman collection file from the repository: [go-url-shortener.postman_collection.json](./docs/go_url_shortener_postman_collection.json).
   - Open Postman, click the **Import** button, and select the collection file (`go-url-shortener.postman_collection.json`).

3. **Run API Requests**  
   Once the collection is imported, you'll see the list of available API endpoints. You can run the requests directly from Postman and interact with the API. Just make sure you have app running in local or in docker.

## Configuration
Environment variables:

- `PORT`: The port number on which the server will run. Default is `8080`.
- `DATABASE_PATH`: The path to the SQLite database. Default is `./db/migrations/database.db`.

## Testing

Run the tests using Go's built-in testing tool:

```bash
go test ./...
```