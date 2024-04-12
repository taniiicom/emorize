package textemoji

import (
	"context"
)

type TextEmojiRepository interface {
	UploadToBucket(ctx context.Context, filePath string) (string, error)
}
