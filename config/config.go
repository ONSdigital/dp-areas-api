package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Configuration represents service configuration for dp-topic-api
type Configuration struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	ZebedeeURL                 string        `envconfig:"ZEBEDEE_URL"`
	EnablePrivateEndpoints     bool          `envconfig:"ENABLE_PRIVATE_ENDPOINTS"`
	EnablePermissionsAuth      bool          `envconfig:"ENABLE_PERMISSIONS_AUTHZ"`
	MongoConfig                MongoConfiguration
}

// MongoConfiguration contains the config required to connect to MongoDB.
type MongoConfiguration struct {
	BindAddr          string `envconfig:"MONGODB_BIND_ADDR"           json:"-"` // This line contains sensitive data and the json:"-" tells the json marshaller to skip serialising it.
	Database          string `envconfig:"MONGODB_TOPICS_DATABASE"`
	TopicsCollection  string `envconfig:"MONGODB_TOPICS_COLLECTION"`
	ContentCollection string `envconfig:"MONGODB_CONTENT_COLLECTION"`
}

var cfg *Configuration

// Get returns the default config with any modifications through environment
// variables
func Get() (*Configuration, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Configuration{
		BindAddr:                   "localhost:25300",
		GracefulShutdownTimeout:    10 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		ZebedeeURL:                 "http://localhost:8082",
		EnablePrivateEndpoints:     true,
		EnablePermissionsAuth:      false,
		MongoConfig: MongoConfiguration{
			BindAddr:          "localhost:27017",
			Database:          "topics",
			TopicsCollection:  "topics",
			ContentCollection: "content",
		},
	}

	return cfg, envconfig.Process("", cfg)
}
