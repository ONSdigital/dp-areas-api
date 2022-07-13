package filereader

import (
	"arearelationshipimport/config"
	"context"
	"fmt"
	"os"
)

type LocalReader struct {
	config *config.ImportConfig
}

func (l LocalReader) GetFile(ctx context.Context, tbl string) (*os.File, error) {
	f, err := os.OpenFile(fmt.Sprintf("%s/table_%s.csv", l.config.GetAreaRelationshipDir(), tbl), os.O_RDONLY, 0777)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func NewLocalReader(config *config.ImportConfig) FileReader {
	return &LocalReader{config: config}
}
