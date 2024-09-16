# go-url-shortener
A simple url shortener written in go.

# To build the docker image
```bash
docker build -t url-shortener .
```

# To run the container
```bash
docker run -d -p 8080:8080 --name go-url-shortener url-shortener
```
