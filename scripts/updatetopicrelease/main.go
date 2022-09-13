package main

import (
	"context"
	"fmt"
	"github.com/ONSdigital/dp-topic-api/mongo"
	"github.com/ONSdigital/dp-topic-api/scripts/config"
	"github.com/ONSdigital/log.go/v2/log"
	"os"
	"strings"
)

const serviceName = "dp-topic-api"

func main() {
	log.Namespace = serviceName
	ctx := context.Background()
	config, err := config.Get()
	if err != nil {
		log.Error(ctx, "error getting config", err)
		os.Exit(1)
	}

	if len(strings.TrimSpace(config.CollectionId)) == 0 {
		log.Error(ctx, "collectionId is required", nil)
		os.Exit(1)
	}
	if len(strings.TrimSpace(config.ReleaseDate)) == 0 {
		log.Error(ctx, "releaseDate is required", nil)
		os.Exit(1)
	}

	err = updateTopicReleaseData(ctx, config, config.CollectionId, config.ReleaseDate)
	if err != nil {
		log.Fatal(ctx, "script fatal error while updating topic", err)
		os.Exit(1)
	}
	fmt.Printf("Collection Id %s has been updated with release date %s", config.CollectionId, config.ReleaseDate)
}
func updateTopicReleaseData(ctx context.Context, c *config.Config, collectionId string, releaseDate string) error {
	conn, err := mongo.NewDBConnection(ctx, c.MongoConfig)
	_, err = conn.UpdateReleaseDate(ctx, collectionId, releaseDate)
	if err != nil {

		return err
	}
	return nil
}
