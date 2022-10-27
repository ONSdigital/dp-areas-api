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
		fmt.Errorf("error getting config: %v", err)
	}

	for true {
		var collectionId string
		var releaseDate string
		fmt.Print("Enter the collection id: ")
		fmt.Scanf("%s", &collectionId)
		validateInputValue(collectionId)

		fmt.Print("Enter the release date: ")
		fmt.Scanf("%s", &releaseDate)
		validateInputValue(releaseDate)
		rd, err := convertToTime(releaseDate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error while converting release date to time.Time type")
			os.Exit(1)
		}

		err = updateTopicReleaseData(ctx, conf.MongoConfig, collectionId, rd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error while updating topic release date")
			os.Exit(1)
		}
		fmt.Printf("Collection Id %s has been updated with release date %v \n", collectionId, rd)
	}
}
func updateTopicReleaseData(ctx context.Context, mongoConfig config.MongoConfig, collectionId string, releaseDate *time.Time) error {
	conn, err := mongo.NewDBConnection(ctx, mongoConfig)
	if err = conn.UpdateReleaseDate(ctx, collectionId, *releaseDate); err != nil {
		return err
	}

	return nil
}

func validateInputValue(inputVal string) {
	if len(strings.TrimSpace(inputVal)) == 0 {
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
