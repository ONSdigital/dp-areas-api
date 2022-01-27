package config

import (
	"os"
	"testing"
	"time"

	"github.com/ONSdigital/dp-mongodb/v3/mongodb"

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
}
