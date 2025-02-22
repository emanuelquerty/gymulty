# gymulty

A multi-tenant rest api for managing gym operations.

Note: This api is still in early stages and updates will be added as they are ready to be shipped.

## Getting started locally

### Download all dependencies

```cmd
go get
```

### Env

Create a .env file in the root of the project with the following configs:

```cmd
DB_NAME=gymulty
DB_USERNAME=postgres
DB_PASSWORD=YourPostgresPassword
```

### Run

```cmd
go run cmd/http/server.go
```

You can also navigate to `cmd/http` directory and run `go run server.go`

### Run tests

```cmd
go test ./...
```

### Build the project

```cmd
go build cmd/http/server.go
```
