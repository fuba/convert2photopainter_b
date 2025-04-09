package dither

import (
	"image"
	"image/color"
	"image/draw"
)

// FindClosestColorFunc はカラーパレットから最も近い色を見つける関数の型定義
type FindClosestColorFunc func(color.Color) color.Color

// FloydSteinberg はFloyd-Steinbergディザリングアルゴリズムを実装する
func FloydSteinberg(img *image.RGBA, palette []color.Color, findClosest FindClosestColorFunc) *image.RGBA {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// 結果画像
	result := image.NewRGBA(bounds)
	draw.Draw(result, bounds, img, bounds.Min, draw.Src)

	// 各ピクセルを処理
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 現在のピクセルを取得
			oldPixel := result.RGBAAt(x, y)
			
			// パレットから最も近い色を見つける
			newColor := findClosest(oldPixel)
			newPixel := color.RGBAModel.Convert(newColor).(color.RGBA)
			
			// 新しい色に設定
			result.SetRGBA(x, y, newPixel)
			
			// 量子化誤差を計算
			errR := int(oldPixel.R) - int(newPixel.R)
			errG := int(oldPixel.G) - int(newPixel.G)
			errB := int(oldPixel.B) - int(newPixel.B)
			
			// 誤差拡散（Floyd-Steinbergアルゴリズム）
			// 右のピクセル (7/16)
			if x+1 < width {
				adjustPixel(result, x+1, y, errR, errG, errB, 7.0/16.0)
			}
			
			// 左下のピクセル (3/16)
			if x-1 >= 0 && y+1 < height {
				adjustPixel(result, x-1, y+1, errR, errG, errB, 3.0/16.0)
			}
			
			// 下のピクセル (5/16)
			if y+1 < height {
				adjustPixel(result, x, y+1, errR, errG, errB, 5.0/16.0)
			}
			
			// 右下のピクセル (1/16)
			if x+1 < width && y+1 < height {
				adjustPixel(result, x+1, y+1, errR, errG, errB, 1.0/16.0)
			}
		}
	}
	
	return result
}

// 指定したピクセルに誤差を追加する
func adjustPixel(img *image.RGBA, x, y, errR, errG, errB int, factor float64) {
	pixel := img.RGBAAt(x, y)
	
	// 誤差を拡散
	r := clamp(int(pixel.R) + int(float64(errR)*factor))
	g := clamp(int(pixel.G) + int(float64(errG)*factor))
	b := clamp(int(pixel.B) + int(float64(errB)*factor))
	
	// 新しい色を設定
	img.SetRGBA(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), pixel.A})
}

// 値を0～255の範囲に収める
func clamp(v int) int {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return v
}
