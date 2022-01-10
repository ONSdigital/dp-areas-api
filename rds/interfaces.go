package rds

import (
	"github.com/aws/aws-sdk-go/service/rds"
)

//go:generate moq -out mock/rds.go -pkg mock . Client

// Client interface is abstraction of rds sdk package
type Client interface {
	DescribeDBInstances(input *rds.DescribeDBInstancesInput) (*rds.DescribeDBInstancesOutput, error)
}
