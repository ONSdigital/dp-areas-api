package config

import (
	"time"

	"github.com/ONSdigital/dp-mongodb/v3/mongodb"

	"github.com/kelseyhightower/envconfig"
)

type MongoConfig = mongodb.MongoDriverConfig

// Config represents service config for dp-topic-api
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	EnablePermissionsAuth      bool          `envconfig:"ENABLE_PERMISSIONS_AUTHZ"`
	EnablePrivateEndpoints     bool          `envconfig:"ENABLE_PRIVATE_ENDPOINTS"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	MongoConfig
	NavigationCacheMaxAge time.Duration `envconfig:"NAVIGATION_CACHE_MAX_AGE"`
	ZebedeeURL            string        `envconfig:"ZEBEDEE_URL"`
}

var cfg *Config

const (
	TopicsCollection  = "TopicsCollection"
	ContentCollection = "ContentCollection"
)

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		BindAddr:                   "localhost:25300",
		EnablePermissionsAuth:      false,
		EnablePrivateEndpoints:     false,
		GracefulShutdownTimeout:    10 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		MongoConfig: MongoConfig{
			ClusterEndpoint:               "localhost:27017",
			Username:                      "",
			Password:                      "",
			Database:                      "topics",
			Collections:                   map[string]string{TopicsCollection: "topics", ContentCollection: "content"},
			ReplicaSet:                    "",
			IsStrongReadConcernEnabled:    false,
			IsWriteConcernMajorityEnabled: true,
			ConnectTimeout:                5 * time.Second,
			QueryTimeout:                  15 * time.Second,
			TLSConnectionConfig: mongodb.TLSConnectionConfig{
				IsSSL: false,
			},
		},
		NavigationCacheMaxAge: 30 * time.Minute,
		ZebedeeURL:            "http://localhost:8082",
	}

	return cfg, envconfig.Process("", cfg)
}
