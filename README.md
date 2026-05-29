## Setup Instructions for Swagger Documentation in a Go Project

### Prerequisites
- **Go Version**: Ensure you have Go version 1.23 or newer installed.

### Installation Steps

1. **Install Swag CLI**: This command installs the Swag CLI tool, which is used to generate Swagger documentation.
    ```sh
    go install github.com/swaggo/swag/cmd/swag@latest
    ```

2. **Get Swag Package**: This command fetches the Swag package, which is necessary for generating Swagger documentation.
    ```sh
    go get -u github.com/swaggo/swag
    ```

3. **Get Gin-Swagger Middleware**: This command fetches the Gin-Swagger middleware, which integrates Swagger with the Gin web framework.
    ```sh
    go get -u github.com/swaggo/gin-swagger
    ```

4. **Get Swag Files**: This command fetches the Swag files package, which contains the Swagger UI files.
    ```sh
    go get -u github.com/swaggo/files
    ```

5. **Get Gin Web Framework**: This command fetches the Gin web framework, which is used to build the Go web application.
    ```sh
    go get -u github.com/gin-gonic/gin
    ```

6. **Clean Up Dependencies**: This command cleans up the module dependencies, ensuring that only the necessary packages are included.
    ```sh
    go mod tidy
    ```

### Initialize Swagger Documentation

- **Generate Documentation**: This command initializes Swagger documentation for the Go project. The `-g` flag specifies the main Go file to parse for generating the documentation. In this case, `cmd/main.go` is the entry point of the application.
    ```sh
    swag init -g cmd/main.go
    ```

    ```sh
    go get github.com/gorilla/websocket
    ```


    ```sh
    go get github.com/google/uuid
    ```
    