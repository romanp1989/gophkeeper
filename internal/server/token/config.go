package token

import "time"

type Config struct {
	Secret string
	Name   string
	Expire time.Duration
}
