# Go REST API Boilerplate with Echo

A modern, feature-rich RESTful API boilerplate built with Go and Echo framework. This project provides a solid foundation for building scalable and maintainable web services with best practices baked in.

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Echo Framework](https://img.shields.io/badge/Echo-v4-00ADD8?style=flat&logo=go)](https://echo.labstack.com/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## 🚀 Features

- **[Echo Framework](https://echo.labstack.com/)**: High performance, extensible, minimalist Go web framework
- **[GORM](https://gorm.io/)**: The fantastic ORM library for Golang
- **[Casbin](https://casbin.org/)**: Powerful and efficient open-source access control library
- **[RabbitMQ](https://www.rabbitmq.com/)**: Message broker for asynchronous processing
- **[WebSockets](https://github.com/gorilla/websocket)**: Real-time communication between server and clients
- **[Swagger](https://github.com/swaggo/swag)**: Automated API documentation
- **[Viper](https://github.com/spf13/viper)**: Complete configuration solution
- **[Zerolog](https://github.com/rs/zerolog)**: Zero allocation JSON logger
- **[Validator](https://github.com/go-playground/validator)**: Request validation
- **Rate Limiting**: Built-in protection against DoS attacks
- **Graceful Shutdown**: Proper handling of server shutdown

## 📍 Prerequisites

- Go 1.24+ ([installation instructions](https://golang.org/doc/install))
- PostgreSQL ([installation instructions](https://www.postgresql.org/download/))
- RabbitMQ ([installation instructions](https://www.rabbitmq.com/download.html))
- Docker and Docker Compose (optional, for containerized development)

## 🛴️ Quick Start

### Clone the repository

```sh
git clone <repository-url>
cd go-rest-boilerplate-echo
```

### Install dependencies

```sh
make dep
# or manually:'
go mod download
```

### Configure the application

```sh
# Copy the example config file
cp config.yaml.example config.yaml

# Edit the config file with your database and RabbitMQ credentials
# using your favorite text editor
vim config.yaml
```

### Start services with Docker (optional)

The project includes a Docker Compose configuration for local development:

```sh
make start-services
```

This will start PostgreSQL and RabbitMQ containers.

### Run database migrations

```sh
go run ./database/migrate/migrate.go
```

### Launch the server

```sh
make serve
# or manually:
go run main.go
```

The API will be available at http://localhost:8080

### API Documentation

Once the server is running, access the Swagger documentation at:

```
http://localhost:8080/swagger/index.html
```

## 👝 Project Structure

```
┋── config/               # Configuration files and functionality
┋── database/             # Database connection and models
│── ┋── migrate/           # Database migration scripts
│── ┋── models/            # Database models (GORM)
┋── docs/                 # API documentation (Swagger)
┋── handlers/             # HTTP request handlers
┋── policy/               # Authorization policies (Casbin)
┋── rabbitmq/             # Message queue clients and task management
┋── websocket/            # WebSocket server implementation
┋── main.go               # Application entry point
┋── config.yaml.example   # Example configuration file
```

## 🔌 API Endpoints

The boilerplate includes a fully-functional tenant management API:

| Method | Endpoint          | Description                     |
|--------|-------------------|--------------------------------|
| GET    | /tenants          | Get a list of all tenants       |
| GET    | /tenants/:id      | Get a specific tenant by ID     |
| POST   | /tenants          | Create a new tenant             |
| PUT    | /tenants          | Update an existing tenant       |
| DELETE | /tenants/:id      | Delete a tenant                 |
| GET    | /swagger/*        | Swagger API documentation       |

## 🔠 Authentication and Authorization

The boilerplate uses Casbin for access control. Policies are defined in the `policy/policy.go` file and can be customized according to your requirements.

## 🔨 Asynchronous Processing

Tasks are processed asynchronously using RabbitMQ. The system includes:

- Push and listen queues
- Task status tracking
- Real-time updates via WebSockets

## 🧪 Testing

Run tests using:

```sh
make test
# or manually:
go test ./...
```

For test coverage:

```sh
make coverage
```

## 📧 Docker Support

Build the Docker image:

```sh
docker build -t go-rest-boilerplate-echo .
```

Run the container:

```sh
docker run -p 8080:8080 go-rest-boilerplate-echo
```

## 👌 Development Guidelines

### Code Style

This project follows the standard Go code style conventions:

```sh
# Format code
go fmt ./...

# Check for code issues
go vet ./...
```

### Adding a New Model

1. Create a new model file in the `database/models/` directory
2. Add the model to the migration process in `database/migrate/migrate.go`
3. Create a handler in the `handlers/` directory
4. Register routes in `main.go`
5. Add authorization policies in `policy/policy.go`

## 🥁 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📏 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👗 Acknowledgements

- [Echo Framework](https://echo.labstack.com/)
- [GORM](https://gorm.io/)
- [Casbin](https://casbin.org/)
- [RabbitMQ](https://www.rabbitmq.com/)
- [Swagger](https://swagger.io/)
- [Viper](https://github.com/spf13/viper)
