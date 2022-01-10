package models

import (
	"github.com/aws/aws-sdk-go/service/rds"
)

// BuildDescibeDBInstancesRequest builds a correctly populated DescribeDBInstancesInput object using the required db instance name
func BuildDescibeDBInstancesRequest(instanceName *string) *rds.DescribeDBInstancesInput {
	return &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: instanceName,
	}
}
