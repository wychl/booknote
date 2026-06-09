package export

import (
	"context"
)

type Screenshot interface {
	Capture(ctx context.Context, htmlPath, outputPath string) error
}
