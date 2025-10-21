package integration

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	postgrestc "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	// Database configuration
	postgresImage    = "postgres:16-alpine"
	postgresDatabase = "testdb"
	postgresUsername = "postgres"
	postgresPassword = "postgres"
	postgresSSLMode  = "sslmode=disable"

	// Container configuration
	containerLogMessage    = "database system is ready to accept connections"
	containerLogOccurrence = 2

	// Connection pool configuration
	maxOpenConns    = 50
	maxIdleConns    = 10
	connMaxLifetime = 5 * time.Minute

	// Timeouts
	containerStartupTimeout = 1 * time.Minute

	// Paths
	migrationsPath = "liquibase/changelog/migrations"
	rootPathLevels = "../../.."

	// messages
	errFailedToCreateFixtureLoader = "failed to create fixture loader"
	errFailedToLoadFixtures        = "failed to load fixtures"
	errFailedToReadMigrationsDir   = "failed to read migrations directory"
	errFailedToReadMigration       = "failed to read migration: %s"
	errFailedToExecuteMigration    = "failed to execute migration: %s"
	shortModeSkipMessage           = "skipping integration test in short mode"
)

type Suite struct {
	suite.Suite
	DB             *gorm.DB
	container      testcontainers.Container
	cleanup        func()
	fixturesPath   string
	skipShortTests bool
}

// WithFixtures configures the suite to load fixtures from the specified directory
func (s *Suite) WithFixtures(path string) *Suite {
	s.fixturesPath = path
	return s
}

// WithSkipShortTests configures the suite to automatically skip tests when running in short mode
func (s *Suite) WithSkipShortTests() *Suite {
	s.skipShortTests = true
	return s
}

// BeforeTest is called before each test method and automatically skips if configured
func (s *Suite) BeforeTest(suiteName, testName string) {
	if s.skipShortTests && testing.Short() {
		s.T().Skip(shortModeSkipMessage)
	}
}

// SetupTest initializes the database before each test
func (s *Suite) SetupTest(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := postgrestc.Run(ctx,
		postgresImage,
		postgrestc.WithDatabase(postgresDatabase),
		postgrestc.WithUsername(postgresUsername),
		postgrestc.WithPassword(postgresPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog(containerLogMessage).
				WithOccurrence(containerLogOccurrence).
				WithStartupTimeout(containerStartupTimeout)),
	)
	require.NoError(t, err)

	connStr, err := pgContainer.ConnectionString(ctx, postgresSSLMode)
	require.NoError(t, err)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	sqlDB, err := db.DB()
	require.NoError(t, err)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	migrationsDir := filepath.Join(rootPath(), migrationsPath)
	applyMigrations(t, sqlDB, migrationsDir)

	if s.fixturesPath != "" {
		loadFixtures(t, sqlDB, s.fixturesPath)
	}

	s.DB = db
	s.container = pgContainer
	s.cleanup = func() {
		_ = pgContainer.Terminate(ctx)
	}
}

// TearDownTest cleans up the database after each test
func (s *Suite) TearDownTest() {
	if s.cleanup != nil {
		s.cleanup()
	}
}

// SkipIfShort skips the test if running in short mode
func (s *Suite) SkipIfShort() {
	if testing.Short() {
		s.T().Skip("Skipping integration test in short mode")
	}
}

// NewUUID generates a new UUID string for idempotency keys
func (s *Suite) NewUUID() string {
	return uuid.New().String()
}

// AssertNoErrors asserts that no errors were received from the error channel
func (s *Suite) AssertNoErrors(errors <-chan error) {
	for err := range errors {
		s.T().Errorf("Unexpected error: %v", err)
	}
}

// ExecuteConcurrent executes a function concurrently the specified number of times
// Returns a channel that will receive any errors from the concurrent executions
func (s *Suite) ExecuteConcurrent(count int, fn func()) <-chan error {
	var wg sync.WaitGroup
	errors := make(chan error, count)

	for i := 0; i < count; i++ {
		wg.Go(fn)
	}

	go func() {
		wg.Wait()
		close(errors)
	}()

	return errors
}

// loadFixtures loads YAML fixtures from the specified directory
func loadFixtures(t *testing.T, sqlDB *sql.DB, fixturesDir string) {
	fixtures, err := testfixtures.New(
		testfixtures.Database(sqlDB),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory(fixturesDir),
	)
	require.NoError(t, err, errFailedToCreateFixtureLoader)

	err = fixtures.Load()
	require.NoError(t, err, errFailedToLoadFixtures)
}

func rootPath() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), rootPathLevels)
}

// applyMigrations executes all SQL files found in the specified directory
func applyMigrations(t *testing.T, sqlDB *sql.DB, migrationsDir string) {
	entries, err := os.ReadDir(migrationsDir)
	require.NoError(t, err, errFailedToReadMigrationsDir)

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != (".sql") {
			continue
		}

		path := filepath.Join(migrationsDir, entry.Name())
		content, err := os.ReadFile(path)
		require.NoError(t, err, errFailedToReadMigration, entry.Name())

		_, err = sqlDB.Exec(string(content))
		require.NoError(t, err, errFailedToExecuteMigration, entry.Name())
	}
}
