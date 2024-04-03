package discord

// main.go はアプリケーションのエントリーポイントです.
// 依存関係の設定とアプリケーションの起動を担当します.

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// discord のエントリーポイント
func Discord() {
	// .env -> 環境変数
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file: ", err)
	}

	DISCORD_BOT_TOKEN := os.Getenv("DISCORD_BOT_TOKEN")

	dg, _ := discordgo.New("Bot " + DISCORD_BOT_TOKEN)
	// tmp: 互換性上の理由から省略
	// cf. https://medium.com/@lapfed255/writing-modern-discord-bots-on-go-9e107bb7fcaa
	// if err != nil {
	// 	fmt.Println("エラーが発生しました: ", err)
	// 	return
	// }

	bot := NewBot(dg)
	bot.Setup()

	// ゲートウェイセッションを開放
	// これで, discord からのイベントを受信
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening websocket: ", err)
	}

	fmt.Println("Bot が正常に起動しました. ctrl+c で終了します.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}
