package db

import "database/sql"

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
