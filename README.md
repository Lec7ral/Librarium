# Librarium API

![Go Version](https://img.shields.io/badge/go-1.18+-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)

Librarium is a comprehensive, production-ready RESTful API for managing a digital library system. Built with Go, this project demonstrates a professional, layered architecture and a wide range of features expected in a modern web service.

---

## Features

- **Full CRUD Operations:** Manage books, authors, and users.
- **Advanced API Queries:**
  - **Pagination:** Control the size and page of listed results (`?limit=20&page=1`).
  - **Filtering:** Dynamically filter results by fields like title or author (`?title=Dune`).
  - **Sorting:** Order results by any specified field (`?sort=published_date&order=desc`).
- **Authentication & Authorization:**
  - **JWT Authentication:** Secure endpoints using JSON Web Tokens.
  - **Role-Based Access Control (RBAC):** Differentiated permissions for "members" and "librarians".
- **Complex Business Logic:**
  - **Transactional Operations:** Safely handle book loans and returns, ensuring stock is updated atomically.
  - **Inventory Management:** Keep track of book stock.
- **Performance Optimization:**
  - **N+1 Problem Solved:** Efficient data loading strategy to prevent excessive database queries.
  - **Redis Caching:** High-performance caching layer for frequently accessed data.
- **Professional Tooling:**
  - **Interactive API Documentation:** Automatically generated, interactive documentation via Swagger/OpenAPI.
  - **Configuration Management:** Environment-aware configuration for both local development and production.
  - **CLI Tools:** Separate, secure command-line tools for administrative tasks like database seeding and role management.
  - **Graceful Shutdown:** Ensures the server finishes processing current requests before shutting down.

---

## Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) (version 1.18 or higher)
- [Redis](https://redis.io/topics/quickstart) (running on the default port `localhost:6379` for local development)
- A C compiler (like `gcc` or `MinGW` on Windows) for the `go-sqlite3` driver.

### 1. Clone the Repository

```sh
git clone <your-repository-url>
cd fullAPI
```

### 2. Install Dependencies

This command will download all the necessary libraries defined in `go.mod`.

```sh
go mod tidy
```

### 3. Set Up Local Environment File

This project uses a `.env` file for local development configuration. Create a file named `.env` in the root of the project and paste the following content:

```env
# .env
# Environment variables for local development.

# Server port for the API
SERVER_PORT=8080

# Public URL for Swagger UI (for local development)
PUBLIC_HOST=localhost:8080
PUBLIC_SCHEME=http

# Database file path
DB_DSN=./library.db

# Redis connection
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=

# JWT Secret Key (use a long, random string)
JWT_SECRET_KEY=local_development_secret_key
```

### 4. Run the Database Seeder (Optional but Recommended)

To populate the database with a large volume of sample data (50 authors, 500 books), run the seeder tool. This is great for testing pagination and filtering.

**Note:** Make sure the API is not running when you execute this.

```sh
go run ./tools/seed.go
```

### 5. Run the API Server

```sh
go run ./cmd/api/main.go
```

The server will start, and you should see a log message like:
`Starting server on port :8080`

The API is now running and accessible at `http://localhost:8080`.

---

## Usage

### API Documentation

The best way to explore and interact with the API is through the **auto-generated Swagger documentation**.

Once the server is running, open your browser and navigate to:
**[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)**

From the Swagger UI, you can:
- View all available endpoints.
- See the required parameters and request body structures.
- Read descriptions for each endpoint and its responses.
- **Execute API requests directly from your browser**, including authorizing with a JWT token.

### Administrative Tools

#### Creating an Administrator (Librarian)

By default, all users are created with the `member` role. To promote a user to `librarian`, use the `manage_user` CLI tool.

1. **Stop the API server.**
2. **Ensure the user exists.** (Register them via the API if needed).
3. **Run the command:**

   ```sh
   # Replace "username" with the actual user's name
   go run ./tools/manage_user.go --username="username" --role="librarian"
   ```

4. **Restart the API server.** The user will now have admin privileges.

---

## Project Structure

```
.
├── cmd/api/         # Main application entry point
├── configs/         # Configuration loading
├── docs/            # Auto-generated Swagger/OpenAPI files
├── internal/        # All private application logic
│   ├── database/    # Database initialization and schema
│   ├── handlers/    # HTTP handlers
│   ├── middleware/  # HTTP middlewares
│   ├── models/      # Data structures
│   ├── repository/  # Data access layer (database logic)
│   └── web/         # Shared web utilities (e.g., response helpers)
└── tools/           # Standalone CLI tools (seeder, user management)
```
