package color

import (
	"fmt"
	"math/rand"
)

// 定数
const (
	COLOR_RED    = "#e57272"
	COLOR_ORANGE = "#e5ac72"
	COLOR_YELLOW = "#e5e572"
	COLOR_GREEN  = "#72e572"
	COLOR_CYAN   = "#72e5e5"
	COLOR_BLUE   = "#7272e5"
	COLOR_PURPLE = "#e572e5"
)

type ColorService struct {
}

func NewColorService() *ColorService {
	return &ColorService{}
}

func (s *ColorService) ConvHexColor(colorText string) (string, error) {
	var hexColor string
	switch colorText {
	case "red":
		hexColor = COLOR_RED
	case "orange":
		hexColor = COLOR_ORANGE
	case "yellow":
		hexColor = COLOR_YELLOW
	case "green":
		hexColor = COLOR_GREEN
	case "cyan":
		hexColor = COLOR_CYAN
	case "blue":
		hexColor = COLOR_BLUE
	case "purple":
		hexColor = COLOR_PURPLE
	default:
		return "", fmt.Errorf("invalid color")
	}

	return hexColor, nil
}

func (s *ColorService) GetRandomColor() string {
	colors := []string{
		COLOR_RED,
		COLOR_ORANGE,
		COLOR_YELLOW,
		COLOR_GREEN,
		COLOR_CYAN,
		COLOR_BLUE,
		COLOR_PURPLE,
	}

	// random
	color := colors[rand.Intn(len(colors))]

	return color
}
