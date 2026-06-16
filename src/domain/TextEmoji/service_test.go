package textemoji_test

import (
	"context"
	textemoji "emorize/src/domain/TextEmoji"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

const (
	FONT_PATH = "public/fonts/ZenMaruGothic-Medium.ttf"
)

type noopUploader struct{}

func (noopUploader) UploadToBucket(ctx context.Context, filePath string) (string, error) {
	return "", nil
}

func TestGenerateTextEmoji(t *testing.T) {
	// 使用するフォントファイルのパス。環境に合わせて適切なパスに修正してください。
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to resolve test file path")
	}
	fontPath := filepath.Join(filepath.Dir(currentFile), "../../../", FONT_PATH)
	// TextEmojiService インスタンスの作成
	service := textemoji.NewTextEmojiService(fontPath, noopUploader{})

	// 生成するテキストと色の設定
	text := "はやく_これに_なりたい"
	hexColor := "#FF5733" // 文字の色（赤）

	// テキスト絵文字の生成
	fileName, _, err := service.GenerateTextEmoji(text, hexColor)
	if err != nil {
		t.Fatalf("Failed to generate text emoji: %v", err)
	}
	defer os.Remove(fileName)

	fmt.Printf("Generated text emoji: %s\n", fileName)
}
