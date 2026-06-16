package main

import (
	"context"
	textemoji "emorize/src/domain/TextEmoji"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
)

const (
	FONT_PATH = "public/fonts/ZenMaruGothic-Medium.ttf"
)

type noopUploader struct{}

func (noopUploader) UploadToBucket(ctx context.Context, filePath string) (string, error) {
	return "", nil
}

func main() {
	// 使用するフォントファイルのパス。環境に合わせて適切なパスに修正してください。
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("failed to resolve current file path")
	}
	fontPath := filepath.Join(filepath.Dir(currentFile), "../../../../", FONT_PATH)
	// TextEmojiService インスタンスの作成
	service := textemoji.NewTextEmojiService(fontPath, noopUploader{})

	// 生成するテキストと色の設定
	text := "はやく_これに_なりたい"
	hexColor := "#FF5733" // 文字の色（赤）

	// テキスト絵文字の生成
	fileName, _, err := service.GenerateTextEmoji(text, hexColor)
	if err != nil {
		log.Fatalf("Failed to generate text emoji: %v", err)
	}

	fmt.Printf("Generated text emoji: %s\n", fileName)
}
