package config

import (
	"time"

	"github.com/ONSdigital/dp-mongodb/v3/mongodb"

	"github.com/kelseyhightower/envconfig"
)

// Config represents service config for dp-topic-api
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	ZebedeeURL                 string        `envconfig:"ZEBEDEE_URL"`
	EnablePrivateEndpoints     bool          `envconfig:"ENABLE_PRIVATE_ENDPOINTS"`
	EnablePermissionsAuth      bool          `envconfig:"ENABLE_PERMISSIONS_AUTHZ"`
	MongoConfig                MongoConfig
}

// MongoConfig contains the config required to connect to MongoDB.
type MongoConfig struct {
	BindAddr                      string `envconfig:"MONGODB_BIND_ADDR"   json:"-"`
	Username                      string `envconfig:"MONGODB_USERNAME"    json:"-"`
	Password                      string `envconfig:"MONGODB_PASSWORD"    json:"-"`
	Database                      string `envconfig:"MONGODB_TOPICS_DATABASE"`
	TopicsCollection              string `envconfig:"MONGODB_TOPICS_COLLECTION"`
	ContentCollection             string `envconfig:"MONGODB_CONTENT_COLLECTION"`
	ReplicaSet                    string `envconfig:"MONGODB_REPLICA_SET"`
	IsStrongReadConcernEnabled    bool   `envconfig:"MONGODB_ENABLE_READ_CONCERN"`
	IsWriteConcernMajorityEnabled bool   `envconfig:"MONGODB_ENABLE_WRITE_CONCERN"`

	ConnectTimeoutInSeconds time.Duration `envconfig:"MONGODB_CONNECT_TIMEOUT"`
	QueryTimeoutInSeconds   time.Duration `envconfig:"MONGODB_QUERY_TIMEOUT"`

	mongodb.TLSConnectionConfig
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
		ZebedeeURL:                 "http://localhost:8082",
		EnablePrivateEndpoints:     true,
		EnablePermissionsAuth:      false,
		MongoConfig: MongoConfig{
			BindAddr:                      "localhost:27017",
			Username:                      "",
			Password:                      "",
			Database:                      "topics",
			TopicsCollection:              "topics",
			ContentCollection:             "content",
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
