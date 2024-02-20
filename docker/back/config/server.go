package config

import "time"

const (
	ReadTimeout                   = 5 * time.Second
	WriteTimeout                  = 10 * time.Second
	IdleTimeout                   = 15 * time.Second
	GracefulShutdownTimeout       = 5 * time.Second
	PreflightCacheDurationSeconds = 300
)

type ContextKey string

const ContextUserIDKey ContextKey = "userID"
