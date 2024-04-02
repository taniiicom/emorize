package discord

import (
	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Session *discordgo.Session
}

func NewBot(session *discordgo.Session) *Bot {
	return &Bot{Session: session}
}

func (b *Bot) Setup() {
	b.Session.AddHandler(interactionCreate)
	// 他のセットアップや初期化のコードをここに追加
}

func interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		commandHandler(s, i)
	}
}
