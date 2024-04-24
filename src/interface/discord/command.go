package discord

import (
	color "emorize/src/domain/Color"
	textemoji "emorize/src/domain/TextEmoji"
	"emorize/src/infra/bucket"
	"encoding/base64"
	"fmt"
	"net/url"
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
		// [ack]
		sendAck(s, i)
		// [async]
		responseEmorize(s, i)
	}
}

func sendAck(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		fmt.Println("ack の送信に失敗しました: ", err)
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

func respondAsyncError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	fmt.Println("respondAsyncError: ", message)
	_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: "",
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "oops...: ",
				Description: message,
				Color:       0x1fd1da,
			},
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
	uploader := &bucket.R2Uploader{} // [di]
	te := textemoji.NewTextEmojiService(FONT_PATH, uploader)
	filePath, bucketObjectUrl, err := te.GenerateTextEmoji(text, hexColor)
	if err != nil {
		fmt.Println("Failed to generate text emoji: ", err)
		respondAsyncError(s, i, "Failed to generate text emoji")
		return
	}

	// png を base64 に変換 (Base64 encoding)
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Failed to read file: ", err)
		respondAsyncError(s, i, "Failed to read file")
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
		fmt.Println("Failed to add emoji: ", err)
		respondAsyncError(s, i, "Failed to add emoji\nhint1 : `MANAGE_EMOJI_AND_STICKER` permission is insufficient. Please change from `Server Settings` > `Roles` > `emorize` > `Permissions` > `Manage emojis`.\nhint2 : Your server's emoji slots may be full. Boost your server to increase capacity or organize existing emojis to free up space.")
		return
	}

	// Twitter 共有リンクを生成
	twitterText := url.QueryEscape("I created a new Emoji using #emorize! \n" + bucketObjectUrl + " \n\napp: emorize.megrio.com")
	twitterURL := "https://twitter.com/intent/tweet?text=" + twitterText

	// ボタンを作成
	// shareButton := discordgo.Button{
	// 	Label: "Share on X/Twitter",
	// 	Style: discordgo.LinkButton,
	// 	URL:   twitterURL,
	// }

	// respond
	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: "",
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       " <:" + emoji.Name + ":" + emoji.ID + "> : " + emoji.Name,
				Description: "New Custom-Emoji Created and Now Available!\nYou can use this emoji by typing `:" + emoji.Name + ":`. \n[share](" + twitterURL + ")",
				Color:       0x1fd1da,
				Image: &discordgo.MessageEmbedImage{
					URL: bucketObjectUrl,
				},
			},
		},
		// Components: []discordgo.MessageComponent{
		// 	discordgo.ActionsRow{
		// 		Components: []discordgo.MessageComponent{shareButton},
		// 	},
		// },
	})
	if err != nil {
		fmt.Println("応答に失敗しました: ", err)
	}
}
