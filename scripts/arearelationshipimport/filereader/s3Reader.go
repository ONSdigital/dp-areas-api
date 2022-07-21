package filereader

import (
	"arearelationshipimport/config"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
)

type S3Reader struct {
	config *config.ImportConfig
}

func NewS3Reader(config *config.ImportConfig) FileReader {
	return &S3Reader{config: config}
}

func (s3Reader *S3Reader) GetFile(ctx context.Context, tbl string) (*os.File, error) {
	tempFilepath := fmt.Sprintf("%s/table_%s.csv", s3Reader.config.GetTempPath(), tbl)
	s3Filepath := fmt.Sprintf("%s/table_%s.csv", s3Reader.config.GetS3AreaRelationshipDir(), tbl)

	err := s3Reader.downloadS3File(tempFilepath, s3Filepath)
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(tempFilepath, os.O_RDONLY, 0777)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (s3Reader *S3Reader) downloadS3File(tempFilepath string, s3Filepath string) error {
	file, err := os.Create(tempFilepath)
	if err != nil {
		return err
	}
	defer file.Close()

	sess, _ := session.NewSession(&aws.Config{Region: aws.String(s3Reader.config.GetS3Region())})
	downloader := s3manager.NewDownloader(sess)
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(s3Reader.config.GetS3Bucket()),
			Key:    aws.String(s3Filepath),
		})
	if err != nil {
		return err
	}

	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")

	return nil
}
