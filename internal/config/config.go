// Package config handles API server configuration binding and loading.
package config

import (
	"time"
)

// Config defines configuration options structure for Fantom API server.
type Config struct {
	// AppName holds the name of the application
	AppName string `mapstructure:"app_name"`

	// Server configuration
	Server Server `mapstructure:"server"`

	// Logger configuration
	Log Log `mapstructure:"log"`

	// Node represents the node structure
	Node Node `mapstructure:"node"`

	// IPFS represents the node structure
	Ipfs Ipfs `mapstructure:"ipfs"`

	// Database configuration
	Db Database `mapstructure:"db"`

	// Cache configuration
	Cache Cache `mapstructure:"cache"`
}

// Server represents the GraphQL server configuration
type Server struct {
	BindAddress     string   `mapstructure:"bind"`
	DomainAddress   string   `mapstructure:"domain"`
	CorsOrigin      []string `mapstructure:"cors"`
	ReadTimeout     int64    `mapstructure:"read_timeout"`
	WriteTimeout    int64    `mapstructure:"write_timeout"`
	IdleTimeout     int64    `mapstructure:"idle_timeout"`
	HeaderTimeout   int64    `mapstructure:"header_timeout"`
	ResolverTimeout int64    `mapstructure:"resolver_timeout"`
}

// Log represents the logger configuration
type Log struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// Node represents the Fantom Opera node access configuration
type Node struct {
	Url string `mapstructure:"url"`
}

// Ipfs represents the IPFS node access configuration
type Ipfs struct {
	Url string `mapstructure:"url"`
}

// Database represents the database access configuration.
type Database struct {
	Url    string `mapstructure:"url"`
	DbName string `mapstructure:"db"`
}

// Cache represents the cache sub-system configuration.
type Cache struct {
	Eviction time.Duration `mapstructure:"eviction"`
	MaxSize  int           `mapstructure:"size"`
}
