package convert

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"  // GIF形式のサポート
	_ "image/jpeg" // JPEG形式のサポート
	_ "image/png"  // PNG形式のサポート
	"os"

	"github.com/example/convert2photopainter/internal/dither"
	"github.com/example/convert2photopainter/internal/resize"
	"golang.org/x/image/bmp" // BMPサポート
)

// PhotoPainter(B)の6色
var (
	colorBlack  = color.RGBA{0, 0, 0, 255}
	colorWhite  = color.RGBA{255, 255, 255, 255}
	colorRed    = color.RGBA{255, 0, 0, 255}
	colorGreen  = color.RGBA{0, 255, 0, 255}
	colorBlue   = color.RGBA{0, 0, 255, 255}
	colorYellow = color.RGBA{255, 255, 0, 255}
)

// 利用可能なカラーパレット
var Palette = []color.Color{
	colorBlack,
	colorWhite,
	colorRed,
	colorGreen,
	colorBlue,
	colorYellow,
}

// Options は画像変換のオプションを表す
type Options struct {
	Width      int  // 出力画像の幅
	Height     int  // 出力画像の高さ
	AutoRotate bool // 縦長画像を自動的に90度回転
	Verbose    bool // 詳細なログ出力
}

// 最も近い色をパレットから見つける
func findClosestColor(c color.Color) color.Color {
	r, g, b, _ := c.RGBA()
	r, g, b = r>>8, g>>8, b>>8 // 16bitから8bitに変換

	minDistance := float64(3 * 255 * 255)
	var closestColor color.Color

	for _, paletteColor := range Palette {
		pr, pg, pb, _ := paletteColor.RGBA()
		pr, pg, pb = pr>>8, pg>>8, pb>>8

		// ユークリッド距離の2乗を計算
		distance := float64((int(r)-int(pr))*(int(r)-int(pr)) +
			(int(g)-int(pg))*(int(g)-int(pg)) +
			(int(b)-int(pb))*(int(b)-int(pb)))

		if distance < minDistance {
			minDistance = distance
			closestColor = paletteColor
		}
	}

	return closestColor
}

// ConvertImage は画像を変換してBMPファイルとして保存する
func ConvertImage(inputPath, outputPath string, options Options) error {
	// 画像ファイルを開く
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("ファイルを開けませんでした: %v", err)
	}
	defer file.Close()

	// 画像をデコード
	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("画像のデコードに失敗しました: %v", err)
	}

	// 画像をリサイズ（必要に応じて回転）
	resizedImg := resize.ResizeImage(img, options.Width, options.Height, options.AutoRotate)

	// 新しいRGBA画像を作成
	bounds := resizedImg.Bounds()
	rgbaImg := image.NewRGBA(bounds)
	draw.Draw(rgbaImg, bounds, resizedImg, bounds.Min, draw.Src)

	// Floyd-Steinbergディザリングを適用
	result := dither.FloydSteinberg(rgbaImg, Palette, findClosestColor)

	// BMP形式で保存
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("出力ファイルの作成に失敗しました: %v", err)
	}
	defer outFile.Close()

	// 24bit BMPで保存
	err = bmp.Encode(outFile, result)
	if err != nil {
		return fmt.Errorf("BMPエンコードに失敗しました: %v", err)
	}

	return nil
}
