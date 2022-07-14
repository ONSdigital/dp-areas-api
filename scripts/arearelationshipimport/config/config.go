package config

import (
	"github.com/kelseyhightower/envconfig"
	"os"
)

type ImportConfig struct {
	DatabaseURIString     string `envconfig:"DATABASE_URI_STRING" required:"true" json:"DatabaseURIString,omitempty"`
	AreaRelationshipDir   string `envconfig:"AREA_RELATIONSHIP_DIR" required:"true" json:"AreaRelationshipDir,omitempty"`
	S3AreaRelationshipDir string `envconfig:"S3_AREA_RELATIONSHIP_DIR" required:"true" json:"S3AreaRelationshipDir,omitempty"`
	AWSAccessKey          string `envconfig:"AWS_ACCESS_KEY" required:"true" json:"AWSAccessKey,omitempty"`
	AWSSecretKey          string `envconfig:"AWS_SECRET_KEY" required:"true" json:"AWSSecretKey,omitempty"`
	S3Bucket              string `envconfig:"S3_BUCKET" required:"true" json:"S3Bucket,omitempty"`
	S3Region              string `envconfig:"S3_REGION" required:"true" json:"S3Region,omitempty"`
	ShouldUseS3Source     bool   `envconfig:"SHOULD_USE_S3_SOURCE" required:"true" json:"ShouldUseS3Source"`
}

func GetImportConfig() *ImportConfig {
	conf := &ImportConfig{}
	envconfig.Process("", conf)
	return conf
}

func (ic *ImportConfig) GetAreaRelationshipSchema() []string {
	return []string{"area", "area_name", "area_relationship", "area_type", "relationship_type"}
}

func (ic *ImportConfig) GetAreaRelationshipDir() string {
	return ic.AreaRelationshipDir
}

func (ic *ImportConfig) GetDatabaseURI() string {
	return ic.DatabaseURIString
}

func (ic *ImportConfig) GetTempPath() string {
	return os.TempDir()
}

func (ic *ImportConfig) GetS3AreaRelationshipDir() string {
	return ic.S3AreaRelationshipDir
}

func (ic *ImportConfig) GetS3Region() string {
	return ic.S3Region
}

func (ic *ImportConfig) GetS3Bucket() string {
	return ic.S3Bucket
}

func (ic *ImportConfig) GetShouldUseS3Source() bool {
	return ic.ShouldUseS3Source
}
