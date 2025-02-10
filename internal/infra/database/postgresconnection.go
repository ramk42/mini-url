// Package database handles PostgreSQL connection initialization and management
package database

import (
	"database/sql"
	"fmt"
	"github.com/ramk42/mini-url/internal/infra/env"
	"github.com/ramk42/mini-url/pkg/logger"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

var (
	dbInstance *sql.DB
	once       sync.Once
)

const (
	defaultDBPort          = 5432
	defaultMaxOpenConns    = 10
	defaultMaxIdleConns    = 5
	defaultConnMaxLifetime = 30 * time.Minute
)

type config struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func GetInstance() *sql.DB {
	if dbInstance != nil {
		return dbInstance
	}

	dbConfig := config{
		Host:            os.Getenv("DB_HOST"),
		Port:            env.GetEnvAsInt("DB_PORT", defaultDBPort),
		User:            os.Getenv("DB_USER"),
		Password:        os.Getenv("DB_PASSWORD"),
		DBName:          os.Getenv("DB_NAME"),
		SSLMode:         os.Getenv("DB_SSLMODE"),
		MaxOpenConns:    env.GetEnvAsInt("DB_MAX_OPEN_CONNS", defaultMaxOpenConns),
		MaxIdleConns:    env.GetEnvAsInt("DB_MAX_IDLE_CONNS", defaultMaxIdleConns),
		ConnMaxLifetime: env.GetEnvAsDuration("DB_CONN_MAX_LIFETIME", defaultConnMaxLifetime),
	}
	var err error
	if dbInstance, err = initialize(dbConfig); err != nil {
		log := logger.Instance()
		log.Fatal().Err(err).Msg("failed to initialize database")
	}

	return dbInstance
}

// Initialize sets up the database connection pool
func initialize(cfg config) (*sql.DB, error) {
	var initErr error

	once.Do(func() {
		connStr := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host,
			cfg.Port,
			cfg.User,
			cfg.Password,
			cfg.DBName,
			cfg.SSLMode,
		)

		db, err := sql.Open("postgres", connStr)
		if err != nil {
			initErr = fmt.Errorf("failed to open database: %w", err)
			return
		}

		// Configure connection pool
		db.SetMaxOpenConns(cfg.MaxOpenConns)
		db.SetMaxIdleConns(cfg.MaxIdleConns)
		db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

		// Verify connection
		if err = db.Ping(); err != nil {
			initErr = fmt.Errorf("database ping failed: %w", err)
			return
		}
		log := logger.Instance()
		log.Println("Successfully connected to PostgreSQL database")
		dbInstance = db
	})

	return dbInstance, initErr
}

func Close() error {
	return dbInstance.Close()
}
