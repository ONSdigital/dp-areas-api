package filereader

import (
	"context"
	"os"
)

type FileReader interface {
	GetFile(ctx context.Context, tbl string) (*os.File, error)
}
