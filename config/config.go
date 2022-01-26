package config

import (
	"fmt"
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
	RDSDBName                  string `envconfig:"DBNAME"`
	RDSDBUser                  string `envconfig:"DBUSER"`
	RDSDBHost                  string `envconfig:"DBHOST"`
	RDSDBPort                  string `envconfig:"DBPORT"`
	AWSRegion                  string `envconfig:"AWSREGION"`
	RDSDBInstance1             string `envconfig:"RDSINSTANCE1"`
	RDSDBInstance2             string `envconfig:"RDSINSTANCE2"`
	RDSDBInstance3             string `envconfig:"RDSINSTANCE3"`
	// flag to use local postres instace provided by dp-compose
	DPPostgresLocal            bool   `envconfig:"USEPOSTGRESLOCAL"`
	DPPostgresUserName         string `envconfig:"DPPOSTGRESUSERNAME"`
	DPPostgresUserPassword     string `envconfig:"DPPOSTGRESPASSWORD"`
	DPPostgresLocalPort        string `envconfig:"DPPOSTGRESPORT"`
	DPPostgresLocalDB          string `envconfig:"USEPOSTGRESDB"`

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
	}

	return cfg, envconfig.Process("", cfg)
}

// GetDBEndpoint get sql endpoint
func (c Config) GetDBEndpoint () string {
	return fmt.Sprintf("%s:%s", c.RDSDBHost, c.RDSDBPort)
}

// GetLocalDBConnectionString returns local connection string
func (c Config) GetLocalDBConnectionString () string {
	return fmt.Sprintf("postgres://%s:%s@localhost:%s/%s", c.DPPostgresUserName, c.DPPostgresUserPassword, c.DPPostgresLocalPort, c.DPPostgresLocalDB)
}

// GetRemoteDBConnectionString returns remote connection string
func (c Config) GetRemoteDBConnectionString (authToken string) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", c.RDSDBHost, c.RDSDBPort, c.RDSDBUser, authToken, c.RDSDBName)
}
