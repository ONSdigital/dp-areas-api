package config

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfig(t *testing.T) {
	os.Clearenv()
	var err error
	var configuration *Config

	Convey("Given an environment with no environment variables set", t, func() {
		Convey("Then cfg should be nil", func() {
			So(cfg, ShouldBeNil)
		})

		Convey("When the config values are retrieved", func() {

			Convey("Then there should be no error returned, and values are as expected", func() {
				configuration, err = Get() // This Get() is only called once, when inside this function
				So(err, ShouldBeNil)
				So(configuration, ShouldResemble, &Config{
					BindAddr:                   "localhost:25500",
					GracefulShutdownTimeout:    5 * time.Second,
					HealthCheckInterval:        30 * time.Second,
					HealthCheckCriticalTimeout: 90 * time.Second,
					DPPostgresLocal:            true,
					DPPostgresUserName:         "postgres",
					DPPostgresLocalPort:        "5432",
					DPPostgresLocalDB:          "dp-areas-api",
					RDSDBConnectionTTL:         24 * time.Hour,
					RDSDBMaxConnections:        4,
					RDSDBMinConnections:        1,
					EnablePrivateEndpoints:     true,
					S3Bucket:                   "ons-dp-area-boundaries",
					LoadSampleData:             false,
				})
			})

			Convey("Then a second call to config should return the same config", func() {
				// This achieves code coverage of the first return in the Get() function.
				newCfg, newErr := Get()
				So(newErr, ShouldBeNil)
				So(newCfg, ShouldResemble, cfg)
			})
		})
	})

	Convey("Given the app config contains sensitive secrets", t, func() {
		cfg := &Config{
			BindAddr:                   "localhost:25500",
			GracefulShutdownTimeout:    5 * time.Second,
			HealthCheckInterval:        30 * time.Second,
			HealthCheckCriticalTimeout: 90 * time.Second,
			DPPostgresLocal:            true,
			DPPostgresUserName:         "postgres",
			DPPostgresLocalPort:        "5432",
			DPPostgresLocalDB:          "dp-areas-api",
			RDSDBConnectionTTL:         24 * time.Hour,
			RDSDBMaxConnections:        4,
			RDSDBMinConnections:        1,
			EnablePrivateEndpoints:     true,
			S3Bucket:                   "ons-dp-area-boundaries",
			LoadSampleData:             false,
			AWSAccessKey:               "awsAccessKeyID",     // Sensitive field.
			AWSSecretKey:               "awsSecretAccessKey", // Sensitive field.
		}

		Convey("When the config struct is marshalled to JSON", func() {
			b, err := json.Marshal(cfg)
			So(err, ShouldBeNil)

			out := string(b)

			Convey("Then the output does not contain the AWS Access Key ID", func() {
				So(out, ShouldNotContainSubstring, "AWS_ACCESS_KEY_ID")
				So(out, ShouldNotContainSubstring, "awsAccessKeyID")
			})

			Convey("And the output does not contain the AWS Secret Access Key", func() {
				So(out, ShouldNotContainSubstring, "AWS_SECRET_ACCESS_KEY")
				So(out, ShouldNotContainSubstring, "awsSecretAccessKey")
			})
		})
	})
}
