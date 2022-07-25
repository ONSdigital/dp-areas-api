package importer

import (
	"arearelationshipimport/config"
	"arearelationshipimport/filereader"
	"arearelationshipimport/filewriter"
	"context"
	"fmt"
	"os"

	"github.com/fatih/color"
	"strings"
)

type AreaImporter struct {
	config      *config.ImportConfig
	source      filereader.FileReader
	destination filewriter.FileWriter
}

func NewAreaImporter(config *config.ImportConfig, source filereader.FileReader, destination filewriter.FileWriter) *AreaImporter {
	return &AreaImporter{config: config, source: source, destination: destination}
}

func (ai *AreaImporter) ConfirmConfigsFromUser() bool {
	var s string
	print := color.New(color.FgHiWhite, color.Bold, color.BgRed)
	if ai.config.GetShouldUseS3Source() {
		print.Println("***************************** Importing from S3 storage *****************************")
	} else {
		print.Println("**************************** Importing from local storage ****************************")
	}

	fmt.Printf("StartImport Config %+v \n", ai.config)

	fmt.Printf("If everything is correct please proceed with (y/N): ")
	_, err := fmt.Scan(&s)
	if err != nil {
		panic(err)
	}

	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	if s == "y" || s == "yes" {
		return true
	}
	return false
}

func (ai *AreaImporter) StartImport(ctx context.Context) {
	for _, table := range ai.config.GetAreaRelationshipSchema() {
		ai.ImportForTable(ctx, table)
	}
}

func (ai *AreaImporter) ImportForTable(ctx context.Context, tbl string) {
	file, err := ai.source.GetFile(ctx, tbl)
	if err != nil {
		fmt.Printf("error occurred while reading file from source. error: %+v", err)
	}

	insertedRows, err := ai.destination.Write(ctx, file, tbl)
	if err != nil {
		fmt.Printf("error occurred while writing file to destination. error: %+v", err)
	}

	defer func(fileToClose *os.File) {
		_ = fileToClose.Close()
	}(file)

	fmt.Printf("Total imported %s: %d \n", tbl, insertedRows)
}
