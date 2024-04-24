package textemoji

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// 定数
const CANVAS_HEIGHT = 128
const CANVAS_WIDTH = 128

type TextEmojiService struct {
	fontPath   string
	repository TextEmojiRepository
}

func NewTextEmojiService(fontPath string, repository TextEmojiRepository) *TextEmojiService {
	return &TextEmojiService{
		fontPath:   fontPath,
		repository: repository,
	}
}

func (s *TextEmojiService) GenerateTextEmoji(text string, hexColor string) (string, string, error) {
	fontBytes, err := os.ReadFile(s.fontPath)
	if err != nil {
		return "", "", err
	}

	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return "", "", err
	}

	img := image.NewRGBA(image.Rect(0, 0, 128, 128))
	draw.Draw(img, img.Bounds(), image.Transparent, image.Point{}, draw.Src)

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(f)
	c.SetClip(img.Bounds())
	c.SetDst(img)

	col, err := parseHexColor(hexColor)
	if err != nil {
		return "", "", err
	}
	c.SetSrc(image.NewUniform(col))

	if err := drawText(c, f, text, 128); err != nil {
		return "", "", err
	}

	outFile, err := os.CreateTemp("", "textemoji-*.png")
	if err != nil {
		return "", "", err
	}
	defer outFile.Close()

	if err := png.Encode(outFile, img); err != nil {
		return "", "", err
	}

	// bucket にアップロード
	// outFile.Name() のアップロード部分を追加
	bucketObjectUrl, err := s.repository.UploadToBucket(context.Background(), outFile.Name())
	if err != nil {
		// [ignorable error]
		fmt.Println("Error uploading to bucket:", err)
	}

	return outFile.Name(), bucketObjectUrl, nil
}

func parseHexColor(s string) (color.Color, error) {
	var r, g, b uint8
	_, err := fmt.Sscanf(s, "#%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return nil, err
	}
	return color.NRGBA{R: r, G: g, B: b, A: 0xff}, nil
}

// MeasureString: テキストを描画するのに必要な width を計算
func MeasureString(face font.Face, text string) (width fixed.Int26_6) {
	for _, x := range text {
		awidth, ok := face.GlyphAdvance(rune(x))
		if ok {
			width += awidth
		}
	}
	return width
}

func drawText(c *freetype.Context, font *truetype.Font, text string, width int) error {
	// "_" で改行を分割
	lines := strings.Split(text, "_")

	// fontSize, yPos: y 座標 を定義
	var fontSize float64

	// fontSize, yPos を計算
	if len(lines) == 0 {
		// err
		return fmt.Errorf("too few lines")
	} else if len(lines) < 4 {
		fontSize = (CANVAS_HEIGHT) / float64(len(lines))
	} else {
		// err
		return fmt.Errorf("too many lines")
	}

	// フォントサイズを設定
	c.SetFontSize(fontSize)

	// 分割された各行について描画
	for i, line := range lines {
		opts := truetype.Options{}
		opts.Size = fontSize
		// 設定したフォントサイズでフォントフェイスを生成
		face := truetype.NewFace(font, &opts)

		// 文字列の表示幅を計算
		txtWidth := MeasureString(face, line).Round()
		var scale float64 = 1.0
		// 描画幅が指定された幅を超える場合、スケールを計算してフォントサイズを調整
		if txtWidth > width {
			scale = float64(width) / float64(txtWidth)
		}

		// スケールに基づいてフォントサイズを調整
		c.SetFontSize(fontSize * scale)
		// テキストを中央揃えで描画するための開始ポイントを計算
		pt := freetype.Pt(int(math.Round((float64(CANVAS_WIDTH)-(float64(txtWidth)*scale))/2)), (i+1)*int(fontSize*scale))
		_, err := c.DrawString(line, pt)
		if err != nil {
			return err
		}
	}

	return nil
}
