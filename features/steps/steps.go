package steps

import (
	"github.com/cucumber/godog"
)

// RegisterSteps registers the specific steps needed to do component tests for the search api
func (c *Component) RegisterSteps(ctx *godog.ScenarioContext) {
	c.APIFeature.RegisterSteps(ctx)
	c.AuthFeature.RegisterSteps(ctx)
	c.RDSFeature.RegisterSteps(ctx)
}
