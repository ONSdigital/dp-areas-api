package main

import (
	"arearelationshipimport/config"
	"arearelationshipimport/filereader"
	"arearelationshipimport/filewriter"
	"arearelationshipimport/importer"
	"context"
	"os"
)

func main() {
	ctx := context.Background()
	cnf := config.GetImportConfig()

	var areaImporter *importer.AreaImporter
	if cnf.GetShouldUseS3Source() {
		areaImporter = getS3ToPGImporter(cnf)
	} else {
		areaImporter = getLocalToPGImporter(cnf)
	}
	confirmation := areaImporter.ConfirmConfigsFromUser()
	if !confirmation {
		os.Exit(0)
	}

	areaImporter.StartImport(ctx)
}

func getLocalToPGImporter(cnf *config.ImportConfig) *importer.AreaImporter {
	localReader := filereader.NewLocalReader(cnf)
	dbWriter := filewriter.NewPostgresWriter(cnf)
	areaImporter := importer.NewAreaImporter(cnf, localReader, dbWriter)
	return areaImporter
}

func getS3ToPGImporter(cnf *config.ImportConfig) *importer.AreaImporter {
	s3Reader := filereader.NewS3Reader(cnf)
	dbWriter := filewriter.NewPostgresWriter(cnf)
	areaImporter := importer.NewAreaImporter(cnf, s3Reader, dbWriter)
	return areaImporter
}
