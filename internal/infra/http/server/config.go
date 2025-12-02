package httpserver

import (
	"net"
	"time"
)

// Config represents the configuration for the HTTP server.
type Config struct {
	addr            string
	shutdownTimeout time.Duration
	ln              net.Listener
	allowedOrigins  string
}

// Option represents an option for the HTTP server.
type Option func(*Config)

// defaultConfig returns the default configuration for the HTTP server.
func defaultConfig() Config {
	return Config{
		addr:            "0.0.0.0:3000",
		shutdownTimeout: 10 * time.Second,
		allowedOrigins:  "*",
	}
}

// WithAddress sets the address for the HTTP server.
func WithAddress(addr string) Option {
	return func(c *Config) {
		c.addr = addr
	}
}

// WithShutdownTimeout sets the shutdown timeout for the HTTP server.
func WithShutdownTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.shutdownTimeout = timeout
	}
}

// WithAllowedOrigins sets the allowed origins for CORS.
func WithAllowedOrigins(allowedOrigins string) Option {
	return func(c *Config) {
		c.allowedOrigins = allowedOrigins
	}
}

// WithListener sets the listener for the HTTP server.
func WithListener(ln net.Listener) Option {
	return func(c *Config) {
		c.ln = ln
	}
}
