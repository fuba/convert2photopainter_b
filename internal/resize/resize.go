package resize

import (
	"image"

	"golang.org/x/image/draw"
)

// RotateImage は画像を90度回転する
func RotateImage(img image.Image, clockwise bool) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	
	// 回転後の新しい画像を作成（幅と高さが入れ替わる）
	var rotated *image.RGBA
	if clockwise {
		rotated = image.NewRGBA(image.Rect(0, 0, height, width))
		
		// 時計回りに90度回転
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				// (x,y) -> (height-y-1, x)
				rotated.Set(height-y-1, x, img.At(x, y))
			}
		}
	} else {
		rotated = image.NewRGBA(image.Rect(0, 0, height, width))
		
		// 反時計回りに90度回転
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				// (x,y) -> (y, width-x-1)
				rotated.Set(y, width-x-1, img.At(x, y))
			}
		}
	}
	
	return rotated
}

// ResizeImage は画像を指定したサイズにリサイズする
// アスペクト比を保持し、必要に応じてクロップする
// autoRotate が true で画像が縦長の場合、90度回転して適切な方向に調整する
func ResizeImage(img image.Image, targetWidth, targetHeight int, autoRotate bool) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	
	// 画像が縦長で自動回転が有効な場合、90度回転
	// PhotoPainter(B)は常に横長表示のため、縦長画像は回転させる
	if autoRotate && float64(width)/float64(height) < 1.0 {
		// 時計回りに90度回転
		img = RotateImage(img, true)
		
		// 回転後のサイズを更新
		bounds = img.Bounds()
		width, height = bounds.Dx(), bounds.Dy()
	}
	
	// ソース画像のアスペクト比と目標アスペクト比を計算
	srcAspect := float64(width) / float64(height)
	dstAspect := float64(targetWidth) / float64(targetHeight)
	
	var srcRect image.Rectangle
	var dstRect image.Rectangle
	
	// アスペクト比に基づいてリサイズ方法を決定
	if srcAspect > dstAspect {
		// ソース画像が目標よりも横長の場合
		// 幅に合わせてリサイズし、高さ方向の中央部分を使用
		newHeight := int(float64(width) / dstAspect)
		y := (height - newHeight) / 2
		srcRect = image.Rect(0, y, width, y+newHeight)
	} else {
		// ソース画像が目標よりも縦長の場合
		// 高さに合わせてリサイズし、幅方向の中央部分を使用
		newWidth := int(float64(height) * dstAspect)
		x := (width - newWidth) / 2
		srcRect = image.Rect(x, 0, x+newWidth, height)
	}
	
	// 出力用の画像を作成
	dstRect = image.Rect(0, 0, targetWidth, targetHeight)
	dst := image.NewRGBA(dstRect)
	
	// リサイズアルゴリズムとしてBilinear（バイリニア）を使用
	// 他のオプション：NearestNeighbor, CatmullRom, ApproxBiLinear など
	draw.ApproxBiLinear.Scale(dst, dstRect, img, srcRect, draw.Over, nil)
	
	return dst
}
