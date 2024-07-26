# go-smq-service
this is a simple message broker server with pub sub features through http.

## How to use
Start the server
```bash 
  go run main.go server.go broker.go
```

## Example Usage

- To publish a message to a topic:
```bash 
  curl "http://localhost:8080/publish?topic=example&message=hello"
```

- To subscribe to a topic:
```bash 
  curl "http://localhost:8080/subscribe?topic=example"
```
