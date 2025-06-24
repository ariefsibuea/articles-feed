package test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type ArticlesFeedTestSuite struct {
	suite.Suite
	dbpool *pgxpool.Pool
	sqlDB  *sql.DB
	echo   *echo.Echo
	ctx    context.Context
}

func (suite *ArticlesFeedTestSuite) SetupSuite() {
	dbHost := getEnv("TEST_DB_HOST", "localhost")
	dbPort := getEnv("TEST_DB_PORT", "5433")
	dbUser := getEnv("TEST_DB_USER", "user_articles_feed_test")
	dbPassword := getEnv("TEST_DB_PASSWORD", "pass_articles_feed_test")
	dbName := getEnv("TEST_DB_NAME", "articles_feed_test")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	suite.Require().NoError(err)

	poolConfig.MaxConns = 5
	poolConfig.MinConns = 1

	dbpool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	suite.Require().NoError(err)

	suite.dbpool = dbpool
	suite.ctx = context.Background()

	suite.sqlDB = stdlib.OpenDBFromPool(dbpool)

	suite.migrateDatabase()
}

func (suite *ArticlesFeedTestSuite) TearDownSuite() {
	suite.Require().NoError(suite.sqlDB.Close(), "failed close sqlDB")
	suite.dbpool.Close()
}

func (suite *ArticlesFeedTestSuite) migrateDatabase() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(filepath.Dir(b))
	migrationPath := "file://" + filepath.Join(basepath, "migrations")

	driver, err := postgres.WithInstance(suite.sqlDB, &postgres.Config{})
	suite.Require().NoError(err, "failed to create postgres driver")

	m, err := migrate.NewWithDatabaseInstance(migrationPath, "postgres", driver)
	suite.Require().NoError(err, "failed to create database instance")

	if err := m.Up(); err != migrate.ErrNoChange {
		suite.Require().NoError(err, "failed to run migrations")
	}

	sourceErr, dbErr := m.Close()
	suite.Require().NoError(sourceErr, "failed to close source")
	suite.Require().NoError(dbErr, "failed to close database")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
