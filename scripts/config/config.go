package config

import (
	"github.com/ONSdigital/dp-mongodb/v3/mongodb"
	"github.com/kelseyhightower/envconfig"
	"time"
)

type MongoConfig = mongodb.MongoDriverConfig

type Config struct {
	MongoConfig
}

var cfg *Config

const (
	TopicsCollection  = "TopicsCollection"
	ContentCollection = "ContentCollection"
)

func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
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
	}
	return cfg, envconfig.Process("", cfg)
}
