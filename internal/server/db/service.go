package db

import (
	"database/sql"
	"github.com/romanp1989/gophkeeper/internal/server/migrate"
	"go.uber.org/zap"
)

// NewDB создает новый экзепляр подключения к БД
func NewDB(cfg *Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.Dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxLifetime(cfg.MaxLifetimeConn)

	return db, nil
}

// initDB Инициализация БД
func InitDB(cfg *Config, logger *zap.Logger) (*sql.DB, error) {
	mCmd, err := migrate.NewMigrateCmd(&migrate.Config{Dsn: cfg.Dsn})
	if err != nil {
		return nil, err
	}

	if err = mCmd.Up(); err != nil {
		return nil, err
	}

	db, err := NewDB(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	return db, nil
}
