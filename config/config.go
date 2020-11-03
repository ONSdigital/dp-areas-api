package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-topic-api
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	MongoConfig                MongoConfiguration
	ZebedeeURL                 string `envconfig:"ZEBEDEE_URL"`
}

// MongoConfiguration contains the config required to connect to MongoDB.
type MongoConfiguration struct {
	BindAddr          string `envconfig:"MONGODB_BIND_ADDR"           json:"-"` // This line contains sensitive data and the json:"-" tells the json marshaller to skip serialising it.
	Database          string `envconfig:"MONGODB_TOPICS_DATABASE"`
	TopicsCollection  string `envconfig:"MONGODB_TOPICS_COLLECTION"`
	ContentCollection string `envconfig:"MONGODB_CONTENT_COLLECTION"`
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		BindAddr:                   "localhost:25300",
		GracefulShutdownTimeout:    10 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		MongoConfig: MongoConfiguration{
			BindAddr:          "localhost:27017",
			Database:          "topics",
			TopicsCollection:  "topics",
			ContentCollection: "content",
		},
		ZebedeeURL: "http://localhost:8082",
	}

	return cfg, envconfig.Process("", cfg)
}
