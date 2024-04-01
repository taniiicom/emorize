package main

import (
	textemoji "emorize/src/domain/TextEmoji"
	"fmt"
	"log"
)

const (
	FONT_PATH = "public/fonts/ZenMaruGothic-Medium.ttf"
)

func main() {
	// 使用するフォントファイルのパス。環境に合わせて適切なパスに修正してください。
	fontPath := FONT_PATH
	// TextEmojiService インスタンスの作成
	service := textemoji.NewTextEmojiService(fontPath)

	// 生成するテキストと色の設定
	text := "気に_なる"
	hexColor := "#FF5733" // 文字の色（赤）

	// テキスト絵文字の生成
	fileName, err := service.GenerateTextEmoji(text, hexColor)
	if err != nil {
		log.Fatalf("Failed to generate text emoji: %v", err)
	}

	fmt.Printf("Generated text emoji: %s\n", fileName)
}
