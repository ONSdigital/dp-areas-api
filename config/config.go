package config

import (
	"time"

	"github.com/ONSdigital/dp-mongodb/v3/mongodb"

	"github.com/kelseyhightower/envconfig"
)

type MongoConfig = mongodb.MongoConnectionConfig

// Config represents service configuration for dp-areas-api
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	DefaultMaxLimit            int           `envconfig:"DEFAULT_MAXIMUM_LIMIT"`
	DefaultLimit               int           `envconfig:"DEFAULT_LIMIT"`
	DefaultOffset              int           `envconfig:"DEFAULT_OFFSET"`
	MongoConfig
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		BindAddr:                   "localhost:25500",
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		DefaultMaxLimit:            1000,
		DefaultLimit:               20,
		DefaultOffset:              0,
		MongoConfig: MongoConfig{
			ClusterEndpoint:               "localhost:27017",
			Username:                      "",
			Password:                      "",
			Database:                      "areas",
			Collection:                    "areas",
			ReplicaSet:                    "",
			IsStrongReadConcernEnabled:    false,
			IsWriteConcernMajorityEnabled: true,
			ConnectTimeoutInSeconds:       5 * time.Second,
			QueryTimeoutInSeconds:         15 * time.Second,
			TLSConnectionConfig: mongodb.TLSConnectionConfig{
				IsSSL: false,
			},
		},
	}

	return cfg, envconfig.Process("", cfg)
}
