package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Session *discordgo.Session
	AppId   string
}

func NewBot(session *discordgo.Session, AppId string) *Bot {
	return &Bot{Session: session, AppId: AppId}
}

func (b *Bot) Setup() {
	// slash-command

	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "ping pong",
		},
		{
			Name:        "emorize",
			Description: "Generate Custom-Emoji from Text",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "text",
					Description: "e.g. Thank_You, 気に_なる",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name",
					Description: "e.g. thx, kininaru",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "color",
					Description: "text color / e.g. red, orange, yellow, green, cyan, blue, purple",
					Required:    false,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "red",
							Value: "red",
						},
						{
							Name:  "orange",
							Value: "orange",
						},
						{
							Name:  "yellow",
							Value: "yellow",
						},
						{
							Name:  "green",
							Value: "green",
						},
						{
							Name:  "cyan",
							Value: "cyan",
						},
						{
							Name:  "blue",
							Value: "blue",
						},
						{
							Name:  "purple",
							Value: "purple",
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "font",
					Description: "text font / e.g. round-gothic, mincho, headline, handwriting, dot",
					Required:    false,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "round-gothic",
							Value: "round-gothic",
						},
						{
							Name:  "mincho",
							Value: "mincho",
						},
						{
							Name:  "headline",
							Value: "headline",
						},
						{
							Name:  "handwriting",
							Value: "handwriting",
						},
						{
							Name:  "dot",
							Value: "dot",
						},
					},
				},
			},
		},
	}

	// slash-command の一括登録
	// guildID が空文字列の場合はグローバルコマンドとしてすべての"サーバ"に適用
	_, err := b.Session.ApplicationCommandBulkOverwrite(b.AppId, "", commands)
	if err != nil {
		fmt.Println("error creating slash-commands: ", err)
		panic(err)
	}
	fmt.Println("slash-commands registered")

	// -

	// interaction-handler

	// interaction に対する応答を処理するハンドラを登録
	// handleSlashCommandInteraction をイベントリスナーとして登録
	b.Session.AddHandler(handleSlashCommandInteraction)

}

func handleSlashCommandInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// slash-commands (= ApplicationCommand) によるリクエストか判定
	if i.Type == discordgo.InteractionApplicationCommand {
		// 具体的なコマンドごとの処理を実行
		commandHandler(s, i)
	}
}
