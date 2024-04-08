package discord

// main.go はアプリケーションのエントリーポイントです.
// 依存関係の設定とアプリケーションの起動を担当します.

import (
	"fmt"
	"net/http"
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

	var (
		DISCORD_BOT_TOKEN = os.Getenv("DISCORD_BOT_TOKEN")
		DISCORD_APP_ID    = os.Getenv("DISCORD_APP_ID")
	)

	// debug:
	fmt.Println("using DISCORD_BOT_TOKEN: ", DISCORD_BOT_TOKEN[:5])

	dg, _ := discordgo.New("Bot " + DISCORD_BOT_TOKEN)
	// tmp: 互換性上の理由から省略
	// cf. https://medium.com/@lapfed255/writing-modern-discord-bots-on-go-9e107bb7fcaa
	// if err != nil {
	// 	fmt.Println("エラーが発生しました: ", err)
	// 	return
	// }

	bot := NewBot(dg, DISCORD_APP_ID)
	bot.Setup()

	// ゲートウェイセッションを開放
	// これで, discord からのイベントを受信
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening websocket: ", err)
	}

	// サーバ起動
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080" // cloud-run が提供する PORT 環境変数がない場合のデフォルト値
			fmt.Println("Defaulting to port: ", port)
		}

		fmt.Println("Listening on port: ", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			panic(err)
		}
	}()

	fmt.Println("Bot が正常に起動しました. ctrl+c で終了します.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}
