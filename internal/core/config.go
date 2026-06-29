package core

import (
	"time"
)

type ProxySpec struct {
	http string `validate:"required,url"`
}

type Config struct {
	Language  string        `validate:"required"`
	Country   string        `validate:"required"`
	Proxy     ProxySpec     `validate:"required"`
	UserAgent string        `validate:"required"`
	Timeout   time.Duration `validate:"required"`
	AutoSkip  bool          `validate:"required"`
}
