# gymulty

A multi-tenant rest api for managing gym operations.

Note: This api is still in early stages and updates will be added as they are ready to be shipped.

## Getting started locally

### Download the project

```cmd
git clone https://github.com/emanuelquerty/gymulty.git
```

### Download all dependencies

Open the project in your editor of choice and run the following command in the terminal:

```cmd
go get
```

### Environment variables

Create a `.env` file in the root of the project with the following configuration variables:

```cmd
DB_HOST=localhost
DB_PORT=5432
DB_NAME=gymulty
DB_USERNAME=postgres
DB_PASSWORD=YourPostgresPassword
```

### Runn the project

```cmd
go run ./cmd/http/.
```

This will start a server listening on port `8080`

### Run tests

```cmd
go test ./...
```

### Build the project

```cmd
go build -o gymulty ./cmd/http
```

This will build the executable named `gymulty` on the project root directory.

### Start the server

Run the generated executable:

```cmd
./gymulty
```

This will spin up the api server on `port 8080`

### Documentation

Once the server is running, you can access the full API Documentation on [http://localhost:8080/api/swagger/index.html](http://localhost:8080/api/swagger/index.html) to interact with the API via swagger UI or your tool of choice.
