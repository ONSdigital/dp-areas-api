package filewriter

import (
	"context"
	"os"
)

type FileWriter interface {
	Write(ctx context.Context, file *os.File, tbl string) (int64, error)
}
