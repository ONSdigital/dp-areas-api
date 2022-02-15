package healthcheck_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/ONSdigital/dp-areas-api/api/mock"
	"github.com/ONSdigital/dp-areas-api/config"

	healthcheck "github.com/ONSdigital/dp-areas-api/service/healthcheck"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"

	"github.com/aws/aws-sdk-go/aws/awserr"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetHealthCheck(t *testing.T) {
	ctx := context.Background()

	m := &mock.RDSAreaStoreMock{}

	cfg, _ := config.Get()

	Convey("dp-areas-api healthchecker reports healthy", t, func() {

		m.PingFunc = func(ctx context.Context) error {
			return nil
		}

		checkState := health.NewCheckState("dp-areas-api-test")
		checker := healthcheck.RDSHealthCheck(context.Background(), cfg, m)
		err := checker(ctx, checkState)
		Convey("When GetHealthCheck is called", func() {
			Convey("Then the HealthCheck flag is set to true and HealthCheck is returned", func() {
				So(checkState.StatusCode(), ShouldEqual, http.StatusOK)
				So(checkState.Status(), ShouldEqual, health.StatusOK)
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("dp-areas-api healthchecker reports critical", t, func() {
		Convey("When the rds instance cannot be successfully pinged", func() {
			m.PingFunc = func(ctx context.Context) error {
				awsErrCode := "ResourceNotFoundException"
				awsErrMessage := "Group not found."
				awsOrigErr := errors.New(awsErrCode)
				awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
				return awsErr
			}

			checkState := health.NewCheckState("dp-areas-api-test")

			checker := healthcheck.RDSHealthCheck(context.Background(), cfg, m)
			err := checker(ctx, checkState)
			Convey("When GetHealthCheck is called", func() {
				Convey("Then the HealthCheck flag is set to true and HealthCheck is returned", func() {
					So(checkState.StatusCode(), ShouldEqual, http.StatusBadGateway)
					So(checkState.Status(), ShouldEqual, health.StatusCritical)
					So(err, ShouldNotBeNil)
				})
			})
		})
	})
}
