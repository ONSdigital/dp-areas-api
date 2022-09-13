package config

import (
	"github.com/ONSdigital/dp-topic-api/config"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	CollectionId string `envconfig:"COLLECTION_ID"`
	ReleaseDate  string `envconfig:"RELEASE_DATE"`
	MongoConfig  config.MongoConfig
}

var cfg *Config

func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}
	return cfg, envconfig.Process("", cfg)
}
