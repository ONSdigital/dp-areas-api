package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ONSdigital/dp-topic-api/config"
	"github.com/ONSdigital/dp-topic-api/mongo"
)

func main() {
	ctx := context.Background()
	conf, err := config.Get()
	if err != nil {
		_ = fmt.Errorf("error getting config: %v", err)
		os.Exit(1)
	}

	for {
		var collectionID string
		var releaseDate string
		fmt.Print("Enter the collection id: ")
		fmt.Scanf("%s", &collectionID)
		validateInputValue(collectionID)

		fmt.Print("Enter the release date: ")
		fmt.Scanf("%s", &releaseDate)
		validateInputValue(releaseDate)
		rd, err := convertToTime(releaseDate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error while converting release date to time.Time type")
			os.Exit(1)
		}

		err = updateTopicReleaseData(ctx, conf.MongoConfig, collectionID, rd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error while updating topic release date")
			os.Exit(1)
		}
		fmt.Printf("Collection Id %s has been updated with release date %v \n", collectionID, rd)
	}
}
func updateTopicReleaseData(ctx context.Context, mongoConfig config.MongoConfig, collectionID string, releaseDate *time.Time) error {
	conn, err := mongo.NewDBConnection(ctx, mongoConfig)
	if err != nil {
		return err
	}

	err = conn.UpdateReleaseDate(ctx, collectionID, *releaseDate)
	if err != nil {
		return err
	}

	return nil
}

func validateInputValue(inputVal string) {
	if strings.TrimSpace(inputVal) == "" {
		fmt.Fprintf(os.Stderr, "invalid input value")
		os.Exit(1)
	}
}

func convertToTime(releaseDate string) (*time.Time, error) {
	rd, err := time.Parse(time.RFC3339, releaseDate)
	if err != nil {
		return nil, err
	}

	return &rd, nil
}
