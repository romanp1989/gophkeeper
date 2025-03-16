package db

import "time"

type Config struct {
	Dsn             string
	MaxIdleConns    int
	MaxOpenConns    int
	MaxLifetimeConn time.Duration
}
