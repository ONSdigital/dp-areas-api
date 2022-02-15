package config

import (
	"fmt"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-areas-api
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	RDSDBName                  string        `envconfig:"DBNAME"`
	RDSDBUser                  string        `envconfig:"DBUSER"`
	RDSDBHost                  string        `envconfig:"DBHOST"`
	RDSDBPort                  string        `envconfig:"DBPORT"`
	AWSRegion                  string        `envconfig:"AWSREGION"`
	RDSDBConnectionTTL         time.Duration `envconfig:"RDSCONNECTIONTTL"`
	RDSDBMaxConnections        int           `envconfig:"RDSMAXCONNECTIONS"`
	RDSDBMinConnections        int           `envconfig:"RDSMINCONNECTIONS"`
	// flag to use local postres instace provided by dp-compose
	DPPostgresLocal        bool   `envconfig:"USEPOSTGRESLOCAL"`
	DPPostgresUserName     string `envconfig:"DPPOSTGRESUSERNAME"`
	DPPostgresUserPassword string `envconfig:"DPPOSTGRESPASSWORD"`
	DPPostgresLocalPort    string `envconfig:"DPPOSTGRESPORT"`
	DPPostgresLocalDB      string `envconfig:"USEPOSTGRESDB"`
}

func (c Config) GetRDSEndpoint() string {
	return fmt.Sprintf("%s:%s", cfg.RDSDBHost, cfg.RDSDBPort)
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
		DPPostgresLocal:            true,
		DPPostgresUserName:         "postgres",
		DPPostgresUserPassword:     os.Getenv("DPPOSTGRESPASSWORD"),
		DPPostgresLocalPort:        "5432",
		DPPostgresLocalDB:          "dp-areas-api",
		RDSDBConnectionTTL:         24 * time.Hour,
		RDSDBMaxConnections:        4,
		RDSDBMinConnections:        1,
	}

	return cfg, envconfig.Process("", cfg)
}

// GetDBEndpoint get sql endpoint
func (c Config) GetDBEndpoint() string {
	return fmt.Sprintf("%s:%s", c.RDSDBHost, c.RDSDBPort)
}

// GetLocalDBConnectionString returns local connection string
func (c Config) GetLocalDBConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@localhost:%s/%s", c.DPPostgresUserName, c.DPPostgresUserPassword, c.DPPostgresLocalPort, c.DPPostgresLocalDB)
}

// GetRemoteDBConnectionString returns remote connection string
func (c Config) GetRemoteDBConnectionString(authToken string) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s pool_max_conns=%d pool_min_conns=%d pool_max_conn_lifetime=%s", c.RDSDBHost, c.RDSDBPort, c.RDSDBUser, authToken, c.RDSDBName, c.RDSDBMaxConnections, c.RDSDBMinConnections, c.RDSDBConnectionTTL)
}
