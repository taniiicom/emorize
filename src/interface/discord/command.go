package discord

import (
	color "emorize/src/domain/Color"
	textemoji "emorize/src/domain/TextEmoji"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

const (
	FONT_PATH = "public/fonts/ZenMaruGothic-Medium.ttf"
)

func commandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {
	case "ping":
		responsePing(s, i)

	case "emorize":
		responseEmorize(s, i)
	}
}

func respondError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	fmt.Println("respondError: ", message)
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "oops...: " + message,
		},
	})
	if err != nil {
		fmt.Println("応答に失敗しました: ", err)
	}
}

func responsePing(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong!",
		},
	})
	if err != nil {
		fmt.Println("応答に失敗しました: ", err)
	}
}

func responseEmorize(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Options の値を取得
	var (
		text      string
		name      string
		colorText string = ""
	)
	for _, option := range i.ApplicationCommandData().Options {
		switch option.Name {
		case "text":
			text = option.StringValue()
		case "name":
			name = option.StringValue()
		case "color":
			colorText = option.StringValue()
		}
	}

	// Color
	var hexColor string
	col := color.NewColorService()
	if colorText == "" {
		hexColor = col.GetRandomColor()
	} else {
		hexColor, _ = col.ConvHexColor(colorText)
	}

	// TextEmoji
	te := textemoji.NewTextEmojiService(FONT_PATH)
	filePath, err := te.GenerateTextEmoji(text, hexColor)
	if err != nil {
		fmt.Println("Failed to generate text emoji: ", err)
		respondError(s, i, "Failed to generate text emoji")
		return
	}

	// png を base64 に変換 (Base64 encoding)
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Failed to read file: ", err)
		respondError(s, i, "Failed to read file")
		return
	}

	encodedEmoji := base64.StdEncoding.EncodeToString(data)

	newEmoji := &discordgo.EmojiParams{
		Name:  name,
		Image: "data:image/png;base64," + encodedEmoji,
	}

	// Emoji を guild に追加
	emoji, err := s.GuildEmojiCreate(i.GuildID, newEmoji)
	if err != nil {
		fmt.Println("Failed to create emoji: ", err)
		respondError(s, i, "Failed to create emoji")
		return
	}

	// respond
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "",
			// Embeds: []*discordgo.MessageEmbed{
			// 	{
			// 		Title: emoji.Name,
			// 		Description: "",
			// 		Image: &discordgo.MessageEmbedImage{
			// 			URL: emoji.,
			// 		},
			// 	},
			// },
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "New Custom-Emoji Created and Now Available!",
					Description: " <:" + emoji.Name + ":" + emoji.ID + "> : " + emoji.Name,
					Color:       0x1fd1da,
				},
			},
		},
	})
	if err != nil {
		fmt.Println("応答に失敗しました: ", err)
	}
}
