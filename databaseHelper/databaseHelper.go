package databaseHelper

import (
	"context"
	"fmt"
	"github.com/LucaSchmitz2003/FlowWatch/loggingHelper"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strconv"
	"sync"
)

var (
	tracer = otel.Tracer("DatabaseTracer")
	logger = loggingHelper.GetLogHelper()

	db   *gorm.DB
	once sync.Once
)

// initDB initializes a new instance of gorm.DB with the specified database connection parameters.
// It panics if it fails to connect to the database, ensuring that the application does not start with an invalid state.
func initDB(ctx context.Context) (*gorm.DB, error) {
	// Start a new span for the database initialization
	ctx, span := tracer.Start(ctx, "initDB")
	defer span.End()

	// Load the environment variables to make sure that the settings have already been loaded
	_ = godotenv.Load(".env")

	// Read the environment variables for the database connection
	host := os.Getenv("DB_HOST")
	if host == "" {
		logger.Error(ctx, "DB_HOST not set, using default")
		host = "db"
	}
	userName := os.Getenv("DB_USERNAME")
	if userName == "" {
		logger.Error(ctx, "DB_USERNAME not set, using default")
		userName = "test"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		logger.Error(ctx, "DB_PASSWORD not set, using default")
		password = "test"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		logger.Error(ctx, "DB_NAME not set, using default")
		dbName = "test"
	}
	timeZone := os.Getenv("TZ")
	if timeZone == "" {
		logger.Error(ctx, "TZ not set, using default")
		timeZone = "Europe/Berlin"
	}
	connectTimeoutSeconds, err := strconv.Atoi(os.Getenv("CONNECT_TIMEOUT_SECONDS"))
	if err != nil {
		err = errors.Wrap(err, "Failed to parse CONNECT_TIMEOUT_SECONDS, using default")
		logger.Error(ctx, err)
		connectTimeoutSeconds = 10
	}
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		err = errors.Wrap(err, "Failed to parse DB_PORT, using default")
		logger.Error(ctx, err)
		port = 5432
	}
	sslMode, err := strconv.ParseBool(os.Getenv("DB_SSL_MODE"))
	if err != nil {
		err = errors.Wrap(err, "Failed to parse DB_SSL_MODE, using default")
		logger.Error(ctx, err)
		sslMode = true
	}

	// Translate the boolean sslMode to a string.
	sslModeString := "enable"
	if !sslMode {
		sslModeString = "disable"
	}
	// Create a new database connection string.
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s connect_timeout=%d",
		host, userName, password, dbName, port, sslModeString, timeZone, connectTimeoutSeconds)

	// Connect to the database.
	db, err1 := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			// Logger: &logging.CustomLogger{},  // TODO: Add custom middleware for logging.
		},
	)
	if err1 != nil {
		err1 = errors.Wrap(err1, "Failed to connect to the database")
		logger.Fatal(ctx, err1)
	}

	// Migrate all models to create the tables:
	err2 := db.AutoMigrate(models...)
	if err2 != nil {
		err2 = errors.Wrap(err2, "Failed to migrate the models")
		logger.Fatal(ctx, err2)
	}

	return db, nil
}

// GetDB creates a new DB instance or returns an already existing instance. It's a singleton pattern.
func GetDB(ctx context.Context) *gorm.DB {
	ctx, span := tracer.Start(ctx, "GetDB")
	defer span.End()

	// Check if models have been defined
	if !modelsSet {
		logger.Fatal(context.Background(), "Models must be registered before calling GetDB. Call RegisterModels first.")
	}

	// Create a new DB instance if it does not exist
	once.Do(func() { // ToDo: Add timeout to prevent deadlock
		var err error = nil
		db, err = initDB(ctx)
		if err != nil {
			err = errors.Wrap(err, "Failed to initialize the database")
			logger.Fatal(context.Background(), err)
		}
	})

	return db
}
