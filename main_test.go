package main

import (
	"flag"
	"fmt"
	"os"
	"testing"

	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-topic-api/features/steps"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

var componentFlag = flag.Bool("component", false, "perform component tests")
var allFlag = flag.Bool("all", false, "perform all tests")

func InitializeScenario(ctx *godog.ScenarioContext) {
	server := NewServer()
	feature := steps.NewMyAppFeature(server.Handler)

	apiFeature := componenttest.NewAPIFeature(feature.Initialiser(server.Handler))

	ctx.BeforeScenario(func(*godog.Scenario) {
		apiFeature.Reset()
	})

	apiFeature.RegisterSteps(ctx)
}

func TestMain(m *testing.M) {
	flag.Parse()
	status := 0
	if *componentFlag || *allFlag {
		var opts = godog.Options{
			Output: colors.Colored(os.Stdout),
			Format: "pretty",
		}

		status = godog.TestSuite{
			Name:                "feature_tests",
			ScenarioInitializer: InitializeScenario,
			Options:             &opts,
		}.Run()
	}

	if !*componentFlag || *allFlag {
		if st := m.Run(); st > status {
			status = st
		}
	}

	if *componentFlag {
		fmt.Printf("coverage: %.1f%s\n", testing.Coverage()*100, "% of all statements")
	}

	os.Exit(status)
}
