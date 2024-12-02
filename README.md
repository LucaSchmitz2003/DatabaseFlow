# DatabaseHelper

DatabaseHelper is an abstraction layer for the database connection, standardizing interaction across the company's ecosystem programs. It provides a singleton instance of a `gorm.DB` connection, supports automatic migration, and integrates with OpenTelemetry for tracing and logging. Its goal is to simplify implementation and ensure consistent usage.

---

## 1. Database Initialization

### Model Registration
To enable migration for your database models, register them globally using `RegisterModels`:
```go
ctx := context.Background()
databaseHelper.RegisterModels(ctx, &YourModel{}, &AnotherModel{})
```

> **Note:** Models must be registered **before** `GetDB()` is called. If no models are registered, the application terminates with an error.

---

### Setup

To initialize the database, use the following method:
```go
ctx := context.Background()
db := databaseHelper.GetDB(ctx)
```

The database connection is configured using environment variables (see section **Environment Variables**). If the connection fails, the application logs the error **and terminates**.

> **Note:** Ensure all required models are registered before calling `GetDB()`.

---

### Automatic Migration
When `GetDB(ctx)` is called, all registered models are automatically migrated to the database:
```go
err := db.AutoMigrate(models...)
if err != nil {
    log.Fatal(err) // Logs and terminates the application on failure
}
```

---

## 2. Environment Variables
The library relies on the following environment variables for database configuration:

| Variable         | Default Value   | Description                                    |
|-------------------|-----------------|------------------------------------------------|
| `DB_HOST`        | `db`            | Hostname of the database server.               |
| `DB_USERNAME`    | `test`          | Database username.                             |
| `DB_PASSWORD`    | `test`          | Database password.                             |
| `DB_NAME`        | `test`          | Name of the database.                          |
| `DB_PORT`        | `5432`          | Port of the database server.                   |
| `DB_SSL_MODE`    | `true`          | Use SSL for database connections (boolean).    |
| `TZ`             | `Europe/Berlin` | Timezone for the database connection.          |

---

## 3. Error Handling
- If no models are registered via `RegisterModels`, `GetDB(ctx)` terminates the application with an appropriate error message.
- Errors during database connection or migration are logged using the logging helper.

---

## 4. OpenTelemetry Integration
DatabaseHelper integrates with OpenTelemetry to provide tracing spans for database initialization and operations.

### Example Tracing
The library creates spans automatically during database initialization. However, you can also use custom spans in your code:
```go
ctx, span := tracer.Start(ctx, "Custom Operation")
defer span.End()
```

---

## 5. Example
```go
package main

import (
	"context"
	"github.com/LucaSchmitz2003/DatabaseFlow/databaseHelper"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name  string
	Email string
}

type Product struct {
	gorm.Model
	Name  string
	Price float64
}

func main() {
	// Register models
	ctx := context.Background()
	databaseHelper.RegisterModels(ctx, &User{}, &Product{})

	// Initialize database
	db := databaseHelper.GetDB(ctx)

	// Example query
	user := &User{Name: "John Doe", Email: "john@example.com"}
	db.Create(user)
}
```

---

## 6. Import in Other Projects
To use DatabaseHelper in your project, import it and set up the required models:
```bash
export GOPRIVATE=github.com/LucaSchmitz2003/*
GIT_SSH_COMMAND="ssh -v" go get github.com/LucaSchmitz2003/DatabaseFlow@main
```
