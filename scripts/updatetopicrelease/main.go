package main

import (
	"context"
	"fmt"
	"github.com/ONSdigital/dp-topic-api/config"
	"github.com/ONSdigital/dp-topic-api/mongo"
	"os"
	"strings"
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

		err = updateTopicReleaseData(ctx, conf.MongoConfig, collectionId, releaseDate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error while updating topic release date")
			os.Exit(1)
		}
		fmt.Printf("Collection Id %s has been updated with release date %s \n", collectionId, releaseDate)
	}
}
func updateTopicReleaseData(ctx context.Context, mongoConfig config.MongoConfig, collectionId string, releaseDate string) error {
	conn, err := mongo.NewDBConnection(ctx, mongoConfig)
	_, err = conn.UpdateReleaseDate(ctx, collectionId, releaseDate)
	if err != nil {
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
