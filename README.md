# 🏗️ Go Service Boilerplate

Welcome to the GoServiceBoilerplate repository, the perfect starting point for building robust and scalable microservices using the Go programming language within our organization. Customize it to meet your specific requirements.

## Structure
`cmd/`

This directory is home to our application's entry points and data seeding.

- `api/`: This is where the application begins its execution.
- `seeder/`: Here, you'll find tools for injecting data into the system, often used for populating the database during development.

`config/`

All things configuration reside here, both in files and structures.

- `config-local.yaml`: This YAML file contains local configuration settings.
- `config.go`: In this Go file, you'll find the configuration structure that the application uses.

`database/migrations/`

This section deals with database migrations, including schema definitions and versioning.

- `schema.py`: This Python file holds the schema for database migrations.
- `versions/`: Within this subdirectory, you'll find organized versions of the database schema for easy tracking and management.

`pkg/`

Here, you'll discover reusable packages and utilities.

- `apperror/`: This is where we handle application errors and provide a global error handler.
- `database/`: Utilities for interacting with the database, such as obtaining connections and managing transactions.
- `httpclient/`: A utility package for making HTTP API requests.
- `observability/`: Contains utilities for observability, including metrics and tracing.
- `resourceful/`: Utilities for handling server-side tasks like pagination, sorting, and filtering.
- `optional/`: Handling patch processes goes here.

`internal/`

This is the core of our application's domain logic.

- `middleware/`: Examples of middleware, including CORS (Cross-Origin Resource Sharing) and guards.
- `server/`: The domain server that initializes the application, often using the Fiber web framework.
- `{sub_domain}/`: This is where we organize subdomains, vendors, and contracts.

  - Domain Layer:
    - `dtos/`: Contains Data Transfer Objects (DTOs) for request and response.
    - `entities/`: Defines domain models.
    - `tabledefinition/`: Definitions for the resourceful package.
  
  - App Layer:
    - `delivery/http/external/`: Handles incoming requests from the front-end to our service.
    - `delivery/http/internal/`: Manages requests between different internal services.
    - `usecase/`: The logic layer.
    - `repository/`: Handles database interactions.
    - `mock/`: Provides mocks for testing purposes.
  
  - `delivery.go`: Defines delivery interfaces.
  - `usecase.go`: Specifies usecase interfaces.
  - `repository.go`: Contains repository interfaces.

`infrastructre/`

- `docker-compose.dev.yml`: You can use this Docker Compose file during development to run the necessary dependencies for the application.

- `Dockerfile`: This Dockerfile serves as the entry point for building the Docker image of our application.


`Makefile`

Our Makefile offers a collection of convenient commands for various tasks, simplifying development and operational processes.

## Contributing
Contributions are always welcome! We encourage contributions from developers to help improve and evolve this boilerplate. Whether it's reaadme.MD, bug fixes, new features, or best practices, your input is invaluable in making this template even better.

## Maintainers
- Fathil Arham (fathil.arham@mceasy.co.id)
- M Dwi Aldi (dwi.aldi@mceasy.co.id)
