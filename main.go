package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/ONSdigital/dp-topic-api/config"
	"github.com/ONSdigital/dp-topic-api/service"
	"github.com/ONSdigital/log.go/log"
	"github.com/pkg/errors"
)

const serviceName = "dp-topic-api"

var (
	// BuildTime represents the time in which the service was built
	BuildTime string
	// GitCommit represents the commit (SHA-1) hash of the service that is running
	GitCommit string
	// Version represents the version of the service that is running
	Version string

	/* NOTE: replace the above with the below to run code with for example vscode debugger.
	BuildTime string = "1601119818"
	GitCommit string = "6584b786caac36b6214ffe04bf62f058d4021538"
	Version   string = "v0.1.0"

	*/
)


func main() {
	log.Namespace = serviceName
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Event(ctx, "fatal runtime error", log.Error(err), log.FATAL)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)

	// Run the service, providing an error channel for fatal errors
	svcErrors := make(chan error, 1)
	svcList := service.NewServiceList(&service.Init{})

	log.Event(ctx, "dp-topic-api version", log.INFO, log.Data{"version": Version})

	// Read config
	cfg, err := config.Get()
	if err != nil {
		return errors.Wrap(err, "error getting config")
	}
	log.Event(ctx, "config on startup", log.INFO, log.Data{"config": cfg, "build_time": BuildTime, "git-commit": GitCommit})

	// Run the service
	svc := service.New(cfg, svcList)
	if err := svc.Run(ctx, BuildTime, GitCommit, Version, svcErrors); err != nil {
		return errors.Wrap(err, "running service failed")
	}

	// Blocks until an os interrupt or a fatal error occurs
	select {
	case err := <-svcErrors:
		log.Event(ctx, "service error received", log.ERROR, log.Error(err))
	case sig := <-signals:
		log.Event(ctx, "os signal received", log.Data{"signal": sig}, log.INFO)
	}
	return svc.Close(ctx)
}
