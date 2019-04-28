package drawimg

import (
	"image"
	"image/color"
	"math"
)

func gaussianBlurKernel(x, sigma float64) float64 {
	return math.Exp(-(x*x)/(2*sigma*sigma)) / (sigma * math.Sqrt(2*math.Pi))
}

func blurHorizontal(src *image.NRGBA, kernel []float64) *image.NRGBA {
	radius := len(kernel) - 1
	width := src.Bounds().Max.X
	height := src.Bounds().Max.Y
	dst := image.NewNRGBA(image.Rect(0, 0, width, height))
	parallel(width, func(partStart, partEnd int) {
		for x := partStart; x < partEnd; x++ {
			start := x - radius
			if start < 0 {
				start = 0
			}
			end := x + radius
			if end > width-1 {
				end = width - 1
			}
			weightSum := 0.0
			for ix := start; ix <= end; ix++ {
				weightSum += kernel[absint(x-ix)]
			}
			for y := 0; y < height; y++ {
				var r, g, b, a float64
				for ix := start; ix <= end; ix++ {
					weight := kernel[absint(x-ix)]
					i := y*src.Stride + ix*4
					wa := float64(src.Pix[i+3]) * weight
					r += float64(src.Pix[i+0]) * wa
					g += float64(src.Pix[i+1]) * wa
					b += float64(src.Pix[i+2]) * wa
					a += wa
				}
				j := y*dst.Stride + x*4
				dst.Pix[j+0] = clamp(r / a)
				dst.Pix[j+1] = clamp(g / a)
				dst.Pix[j+2] = clamp(b / a)
				dst.Pix[j+3] = clamp(a / weightSum)
			}
		}
	})
	return dst
}

func blurVertical(src *image.NRGBA, kernel []float64) *image.NRGBA {
	radius := len(kernel) - 1
	width := src.Bounds().Max.X
	height := src.Bounds().Max.Y
	dst := image.NewNRGBA(image.Rect(0, 0, width, height))
	parallel(height, func(partStart, partEnd int) {
		for y := partStart; y < partEnd; y++ {
			start := y - radius
			if start < 0 {
				start = 0
			}
			end := y + radius
			if end > height-1 {
				end = height - 1
			}
			weightSum := 0.0
			for iy := start; iy <= end; iy++ {
				weightSum += kernel[absint(y-iy)]
			}
			for x := 0; x < width; x++ {
				var r, g, b, a float64
				for iy := start; iy <= end; iy++ {
					weight := kernel[absint(y-iy)]
					i := iy*src.Stride + x*4
					wa := float64(src.Pix[i+3]) * weight
					r += float64(src.Pix[i+0]) * wa
					g += float64(src.Pix[i+1]) * wa
					b += float64(src.Pix[i+2]) * wa
					a += wa
				}
				j := y*dst.Stride + x*4
				dst.Pix[j+0] = clamp(r / a)
				dst.Pix[j+1] = clamp(g / a)
				dst.Pix[j+2] = clamp(b / a)
				dst.Pix[j+3] = clamp(a / weightSum)
			}
		}
	})
	return dst
}

// Sharpen produces a sharpened version of the image.
func Sharpen(img image.Image, radius int, sigma float64) *image.NRGBA {
	if sigma <= 0 {
		// sigma parameter must be positive!
		return Clone(img)
	}
	src := toNRGBA(img)
	blurred := Blur(img, radius, sigma)
	width := src.Bounds().Max.X
	height := src.Bounds().Max.Y
	dst := image.NewNRGBA(image.Rect(0, 0, width, height))
	parallel(height, func(partStart, partEnd int) {
		for y := partStart; y < partEnd; y++ {
			for x := 0; x < width; x++ {
				i := y*src.Stride + x*4
				for j := 0; j < 4; j++ {
					k := i + j
					val := int(src.Pix[k])<<1 - int(blurred.Pix[k])
					if val < 0 {
						val = 0
					} else if val > 255 {
						val = 255
					}
					dst.Pix[k] = uint8(val)
				}
			}
		}
	})
	return dst
}

// Blur produces a blurred version of the image using a Gaussian function.
func Blur(img image.Image, radius int, sigma float64) *image.NRGBA {
	var dst *image.NRGBA
	if sigma <= 0 {
		// sigma parameter must be positive!
		return Clone(img)
	}
	src := toNRGBA(img)
	kernel := make([]float64, radius+1)
	for i := 0; i <= radius; i++ {
		kernel[i] = gaussianBlurKernel(float64(i), sigma)
	}
	dst = blurHorizontal(src, kernel)
	dst = blurVertical(dst, kernel)
	return dst
}

// IsTooBright assume that average brightness higher than 100 is too bright.
func IsTooBright(img image.Image) bool {
	var (
		pixCount, totalBrightness float64
	)
	pixCount = 0
	totalBrightness = 0
	AdjustFunc(img, func(c color.NRGBA) color.NRGBA {
		brightness := 0.2126*float64(c.R) + 0.7152*float64(c.G) + 0.0722*float64(c.B)
		totalBrightness += brightness
		pixCount++
		return c
	})
	averBrightness := totalBrightness / pixCount
	return averBrightness > 100
}
