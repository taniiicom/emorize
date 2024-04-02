package main

import (
	"emorize/src/interface/discord"
)

const (
	FONT_PATH = "public/fonts/ZenMaruGothic-Medium.ttf"
)

func main() {
	// discord bot の起動
	discord.Discord()
}
