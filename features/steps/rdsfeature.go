package steps

import (
	"context"
	"testing"

	"github.com/ONSdigital/dp-areas-api/config"
	"github.com/ONSdigital/dp-areas-api/rds"
	componenttest "github.com/ONSdigital/dp-component-test"

	"github.com/cucumber/godog"
)

type RDSFeature struct {
	componenttest.ErrorFeature
	Client rds.RDS
}

func NewRDSFeature(t *testing.T, cfg *config.Config) *RDSFeature {
	rf := &RDSFeature{
		ErrorFeature: componenttest.ErrorFeature{TB: t},
	}

	if err := rf.Client.Init(context.Background(), cfg); err != nil {
		panic("couldn't start RDSFeature: " + err.Error())
	}

	return rf
}

func (rf *RDSFeature) Reset() {}

func (rf *RDSFeature) Close() { rf.Client.Close() }

func (rf *RDSFeature) RegisterSteps(ctx *godog.ScenarioContext) {
	// Use the steps to inject your setup into the RDS tables
	// Or to check results correspond to what is in the RDS tables
}
