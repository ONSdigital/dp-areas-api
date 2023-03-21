package main

import (
	"context"
	"flag"
	"io"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"

	"github.com/ONSdigital/dp-areas-api/features/steps"
	dplogs "github.com/ONSdigital/log.go/v2/log"
)

var componentFlag = flag.Bool("component", false, "perform component tests")
var quietComponentFlag = flag.Bool("quiet-component", false, "perform component tests with dp logging disabled")

type ComponentTest struct {
	*steps.Component
}

func init() {
	dplogs.Namespace = "dp-areas-api"
}

func (ct *ComponentTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.After(func(ctx context.Context, scenario *godog.Scenario, err error) (context.Context, error) {
		ct.Reset()
		return ctx, nil
	})
}

func (ct *ComponentTest) InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ct.RegisterSteps(ctx.ScenarioContext())

	ctx.AfterSuite(func() {
		ct.Close()
	})
}

func TestComponent(t *testing.T) {
	if *componentFlag || *quietComponentFlag {
		status := 0

		var output io.Writer = os.Stdout

		if *quietComponentFlag {
			dplogs.SetDestination(io.Discard, io.Discard)
		}

		var opts = godog.Options{
			Output:   colors.Colored(output),
			Format:   "pretty",
			Paths:    flag.Args(),
			TestingT: t,
		}

		ct := &ComponentTest{
			Component: steps.NewComponent(t),
		}

		status = godog.TestSuite{
			Name:                 "feature_tests",
			ScenarioInitializer:  ct.InitializeScenario,
			TestSuiteInitializer: ct.InitializeTestSuite,
			Options:              &opts,
		}.Run()

		if status > 0 {
			t.Fail()
		}
	} else {
		t.Skip("component flag required to run component tests")
	}
}
