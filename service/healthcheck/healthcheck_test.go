package healthcheck_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/ONSdigital/dp-areas-api/rds/mock"

	servicerds "github.com/aws/aws-sdk-go/service/rds"

	healthcheck "github.com/ONSdigital/dp-areas-api/service/healthcheck"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"

	"github.com/aws/aws-sdk-go/aws/awserr"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetHealthCheck(t *testing.T) {
	ctx := context.Background()

	m := &mock.ClientMock{}

	Convey("dp-areas-api healthchecker reports healthy", t, func() {

		m.DescribeDBInstancesFunc = func(input *servicerds.DescribeDBInstancesInput) (*servicerds.DescribeDBInstancesOutput, error) {
			testDBName := "testDB1"
			return &servicerds.DescribeDBInstancesOutput{
				DBInstances: []*servicerds.DBInstance{
					{
						DBName: &testDBName,
					},
				},
			}, nil
		}

		checkState := health.NewCheckState("dp-areas-api-test")
		checker := healthcheck.RDSHealthCheck(m)
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
		Convey("When the user pool can't be found", func() {
			m.DescribeDBInstancesFunc = func(input *servicerds.DescribeDBInstancesInput) (*servicerds.DescribeDBInstancesOutput, error) {
				awsErrCode := "ResourceNotFoundException"
				awsErrMessage := "Group not found."
				awsOrigErr := errors.New(awsErrCode)
				awsErr := awserr.New(awsErrCode, awsErrMessage, awsOrigErr)
				return nil, awsErr
			}

			checkState := health.NewCheckState("dp-areas-api-test")

			checker := healthcheck.RDSHealthCheck(m)
			err := checker(ctx, checkState)
			Convey("When GetHealthCheck is called", func() {
				Convey("Then the HealthCheck flag is set to true and HealthCheck is returned", func() {
					So(checkState.StatusCode(), ShouldEqual, http.StatusTooManyRequests)
					So(checkState.Status(), ShouldEqual, health.StatusCritical)
					So(err, ShouldNotBeNil)
				})
			})
		})
	})
}
