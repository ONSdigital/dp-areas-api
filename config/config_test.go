package config

import (
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfig(t *testing.T) {
	os.Clearenv()
	var err error
	var config *Config

	Convey("Given an environment with no environment variables set", t, func() {
		Convey("Then cfg should be nil", func() {
			So(cfg, ShouldBeNil)
		})

		Convey("When the config values are retrieved", func() {

			Convey("Then there should be no error returned, and values are as expected", func() {
				config, err = Get() // This Get() is only called once, when inside this function
				So(err, ShouldBeNil)

				So(config.BindAddr, ShouldEqual, "localhost:25300")
				So(config.GracefulShutdownTimeout, ShouldEqual, 10*time.Second)
				So(config.HealthCheckInterval, ShouldEqual, 30*time.Second)
				So(config.HealthCheckCriticalTimeout, ShouldEqual, 90*time.Second)

				So(config.MongoConfig.BindAddr, ShouldEqual, "localhost:27017")
				So(config.MongoConfig.Database, ShouldEqual, "topics")
				So(config.MongoConfig.TopicsCollection, ShouldEqual, "topics")
				So(config.MongoConfig.ContentCollection, ShouldEqual, "content")
				So(cfg.MongoConfig.Username, ShouldEqual, "test")
				So(cfg.MongoConfig.Password, ShouldEqual, "test")
				So(cfg.MongoConfig.CAFilePath, ShouldEqual, "")
				So(cfg.ZebedeeURL, ShouldEqual, "http://localhost:8082")
				So(cfg.EnablePrivateEndpoints, ShouldEqual, true)
				So(cfg.EnablePermissionsAuth, ShouldBeFalse)
			})

			Convey("Then a second call to config should return the same config", func() {
				// This achieves code coverage of the first return in the Get() function.
				newCfg, newErr := Get()
				So(newErr, ShouldBeNil)
				So(newCfg, ShouldResemble, config)
			})
		})
	})
}
