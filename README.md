# Counter Request API

This API tracks and returns the total number of requests received in the past 60 seconds. It is implemented using Go's standard libraries.

## API URL

- Base URL: `http://localhost:8080/`

## Running the API

To start the API server:

```sh
go run ./cmd/api/

```

## (Optional) Testing with HTTP Requests

You can use the requests.http file to test the API endpoints. This may require installing the HTTP extension in VS Code.
[Rest Client extension](https://marketplace.visualstudio.com/items?itemName=humao.rest-client)


## Running Tests

To execute all tests with detailed output, run:
```sh
go test ./... -v
```