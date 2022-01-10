package config

import (
	"os"
	"time"

	"github.com/ONSdigital/dp-mongodb/v3/mongodb"

	"github.com/kelseyhightower/envconfig"
)

type MongoConfig = mongodb.MongoDriverConfig

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
	RDSDBName                  string		 `envconfig:"DBNAME"`
	RDSDBUser                  string		 `envconfig:"DBUSER"`
	RDSDBHost                  string        `envconfig:"DBHOST"`
	RDSDBPort                  string        `envconfig:"DBPORT"`
	AWSRegion                  string        `envconfig:"AWSREGION"`
}

var cfg *Config

const (
	AreasCollection = "AreasCollection"
)

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
			Collections:                   map[string]string{AreasCollection: "areas"},
			ReplicaSet:                    "",
			IsStrongReadConcernEnabled:    false,
			IsWriteConcernMajorityEnabled: true,
			ConnectTimeout:                5 * time.Second,
			QueryTimeout:                  15 * time.Second,
			TLSConnectionConfig: mongodb.TLSConnectionConfig{
				IsSSL: false,
			},
		},
		RDSDBName: os.Getenv("DBNAME"),
		RDSDBUser: os.Getenv("DBUSER"),
		RDSDBHost: os.Getenv("DBHOST"),
		RDSDBPort: os.Getenv("DBPORT"),
		AWSRegion: os.Getenv("AWSREGION"),
	}

	return cfg, envconfig.Process("", cfg)
}
