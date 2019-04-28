package drawimg

import (
	"errors"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
)

// Format is an image file format.
type Format int

// Image file formats.
const (
	JPEG Format = iota
	PNG
	GIF
)

func (f Format) String() string {
	switch f {
	case JPEG:
		return "JPEG"
	case PNG:
		return "PNG"
	case GIF:
		return "GIF"
	default:
		return "Unsupported"
	}
}

var (
	// ErrUnsupportedFormat means the given image format (or file extension) is unsupported.
	ErrUnsupportedFormat = errors.New("imaging: unsupported image format")
)

// Decode reads an image from r.
func Decode(r io.Reader) (image.Image, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	return toNRGBA(img), nil
}

// Open loads an image from file
func Open(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, err := Decode(file)
	return img, err
}

// Encode writes the image img to w in the specified format (JPEG, PNG, GIF, TIFF or BMP).
func Encode(w io.Writer, img image.Image, format Format) error {
	var err error
	switch format {
	case JPEG:
		var rgba *image.RGBA
		if nrgba, ok := img.(*image.NRGBA); ok {
			if nrgba.Opaque() {
				rgba = &image.RGBA{
					Pix:    nrgba.Pix,
					Stride: nrgba.Stride,
					Rect:   nrgba.Rect,
				}
			}
		}
		if rgba != nil {
			err = jpeg.Encode(w, rgba, &jpeg.Options{Quality: 95})
		} else {
			err = jpeg.Encode(w, img, &jpeg.Options{Quality: 95})
		}
	case PNG:
		err = png.Encode(w, img)
	case GIF:
		err = gif.Encode(w, img, &gif.Options{NumColors: 256})
	default:
		err = ErrUnsupportedFormat
	}
	return err
}

// Save saves the image to file with the specified filename.
// The format is determined from the filename extension: "jpg" (or "jpeg"), "png", "gif" are supported.
func Save(img image.Image, filename string) (err error) {
	formats := map[string]Format{
		".jpg":  JPEG,
		".jpeg": JPEG,
		".png":  PNG,
		".gif":  GIF,
	}
	ext := strings.ToLower(filepath.Ext(filename))
	f, ok := formats[ext]
	if !ok {
		return ErrUnsupportedFormat
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return Encode(file, img, f)
}

// NewNRGBA creates a new image with the specified width and height, and fills it with the specified color.
func NewNRGBA(width, height int, fillColor color.Color) *image.NRGBA {
	if width <= 0 || height <= 0 {
		return &image.NRGBA{}
	}
	dst := image.NewNRGBA(image.Rect(0, 0, width, height))
	c := color.NRGBAModel.Convert(fillColor).(color.NRGBA)
	if c.R == 0 && c.G == 0 && c.B == 0 && c.A == 0 {
		return dst
	}
	cs := []uint8{c.R, c.G, c.B, c.A}
	// fill the first row
	for x := 0; x < width; x++ {
		copy(dst.Pix[x*4:(x+1)*4], cs)
	}
	// copy the first row to other rows
	for y := 1; y < height; y++ {
		copy(dst.Pix[y*dst.Stride:y*dst.Stride+width*4], dst.Pix[0:width*4])
	}
	return dst
}

// Clone returns a copy of the given image.
func Clone(img image.Image) *image.NRGBA {
	srcBounds := img.Bounds()
	srcMinX := srcBounds.Min.X
	srcMinY := srcBounds.Min.Y
	dstBounds := srcBounds.Sub(srcBounds.Min)
	dstW := dstBounds.Dx()
	dstH := dstBounds.Dy()
	dst := image.NewNRGBA(dstBounds)
	switch src := img.(type) {
	case *image.NRGBA:
		rowSize := srcBounds.Dx() * 4
		parallel(dstH, func(partStart, partEnd int) {
			for dstY := partStart; dstY < partEnd; dstY++ {
				di := dst.PixOffset(0, dstY)
				si := src.PixOffset(srcMinX, srcMinY+dstY)
				copy(dst.Pix[di:di+rowSize], src.Pix[si:si+rowSize])
			}
		})
	case *image.NRGBA64:
		parallel(dstH, func(partStart, partEnd int) {
			for dstY := partStart; dstY < partEnd; dstY++ {
				di := dst.PixOffset(0, dstY)
				si := src.PixOffset(srcMinX, srcMinY+dstY)
				for dstX := 0; dstX < dstW; dstX++ {
					dst.Pix[di+0] = src.Pix[si+0]
					dst.Pix[di+1] = src.Pix[si+2]
					dst.Pix[di+2] = src.Pix[si+4]
					dst.Pix[di+3] = src.Pix[si+6]
					di += 4
					si += 8
				}
			}
		})
	case *image.RGBA:
		parallel(dstH, func(partStart, partEnd int) {
			for dstY := partStart; dstY < partEnd; dstY++ {
				di := dst.PixOffset(0, dstY)
				si := src.PixOffset(srcMinX, srcMinY+dstY)
				for dstX := 0; dstX < dstW; dstX++ {
					a := src.Pix[si+3]
					dst.Pix[di+3] = a
					switch a {
					case 0:
						dst.Pix[di+0] = 0
						dst.Pix[di+1] = 0
						dst.Pix[di+2] = 0
					case 0xff:
						dst.Pix[di+0] = src.Pix[si+0]
						dst.Pix[di+1] = src.Pix[si+1]
						dst.Pix[di+2] = src.Pix[si+2]
					default:
						var tmp uint16
						tmp = uint16(src.Pix[si+0]) * 0xff / uint16(a)
						dst.Pix[di+0] = uint8(tmp)
						tmp = uint16(src.Pix[si+1]) * 0xff / uint16(a)
						dst.Pix[di+1] = uint8(tmp)
						tmp = uint16(src.Pix[si+2]) * 0xff / uint16(a)
						dst.Pix[di+2] = uint8(tmp)
					}
					di += 4
					si += 4
				}
			}
		})
	case *image.RGBA64:
		parallel(dstH, func(partStart, partEnd int) {
			for dstY := partStart; dstY < partEnd; dstY++ {
				di := dst.PixOffset(0, dstY)
				si := src.PixOffset(srcMinX, srcMinY+dstY)
				for dstX := 0; dstX < dstW; dstX++ {
					a := src.Pix[si+6]
					dst.Pix[di+3] = a
					switch a {
					case 0:
						dst.Pix[di+0] = 0
						dst.Pix[di+1] = 0
						dst.Pix[di+2] = 0
					case 0xff:
						dst.Pix[di+0] = src.Pix[si+0]
						dst.Pix[di+1] = src.Pix[si+2]
						dst.Pix[di+2] = src.Pix[si+4]
					default:
						var tmp uint16
						tmp = uint16(src.Pix[si+0]) * 0xff / uint16(a)
						dst.Pix[di+0] = uint8(tmp)
						tmp = uint16(src.Pix[si+2]) * 0xff / uint16(a)
						dst.Pix[di+1] = uint8(tmp)
						tmp = uint16(src.Pix[si+4]) * 0xff / uint16(a)
						dst.Pix[di+2] = uint8(tmp)
					}
					di += 4
					si += 8
				}
			}
		})
	case *image.Gray:
		parallel(dstH, func(partStart, partEnd int) {
			for dstY := partStart; dstY < partEnd; dstY++ {
				di := dst.PixOffset(0, dstY)
				si := src.PixOffset(srcMinX, srcMinY+dstY)
				for dstX := 0; dstX < dstW; dstX++ {
					c := src.Pix[si]
					dst.Pix[di+0] = c
					dst.Pix[di+1] = c
					dst.Pix[di+2] = c
					dst.Pix[di+3] = 0xff
					di += 4
					si++
				}
			}
		})
	case *image.Gray16:
		parallel(dstH, func(partStart, partEnd int) {
			for dstY := partStart; dstY < partEnd; dstY++ {
				di := dst.PixOffset(0, dstY)
				si := src.PixOffset(srcMinX, srcMinY+dstY)
				for dstX := 0; dstX < dstW; dstX++ {
					c := src.Pix[si]
					dst.Pix[di+0] = c
					dst.Pix[di+1] = c
					dst.Pix[di+2] = c
					dst.Pix[di+3] = 0xff
					di += 4
					si += 2
				}
			}
		})
	case *image.YCbCr:
		parallel(dstH, func(partStart, partEnd int) {
			for dstY := partStart; dstY < partEnd; dstY++ {
				di := dst.PixOffset(0, dstY)
				for dstX := 0; dstX < dstW; dstX++ {
					srcX := srcMinX + dstX
					srcY := srcMinY + dstY
					siy := src.YOffset(srcX, srcY)
					sic := src.COffset(srcX, srcY)
					r, g, b := color.YCbCrToRGB(src.Y[siy], src.Cb[sic], src.Cr[sic])
					dst.Pix[di+0] = r
					dst.Pix[di+1] = g
					dst.Pix[di+2] = b
					dst.Pix[di+3] = 0xff
					di += 4
				}
			}
		})

	case *image.Paletted:
		plen := len(src.Palette)
		pnew := make([]color.NRGBA, plen)
		for i := 0; i < plen; i++ {
			pnew[i] = color.NRGBAModel.Convert(src.Palette[i]).(color.NRGBA)
		}
		parallel(dstH, func(partStart, partEnd int) {
			for dstY := partStart; dstY < partEnd; dstY++ {
				di := dst.PixOffset(0, dstY)
				si := src.PixOffset(srcMinX, srcMinY+dstY)
				for dstX := 0; dstX < dstW; dstX++ {
					c := pnew[src.Pix[si]]
					dst.Pix[di+0] = c.R
					dst.Pix[di+1] = c.G
					dst.Pix[di+2] = c.B
					dst.Pix[di+3] = c.A
					di += 4
					si++
				}
			}
		})
	default:
		parallel(dstH, func(partStart, partEnd int) {
			for dstY := partStart; dstY < partEnd; dstY++ {
				di := dst.PixOffset(0, dstY)
				for dstX := 0; dstX < dstW; dstX++ {
					c := color.NRGBAModel.Convert(img.At(srcMinX+dstX, srcMinY+dstY)).(color.NRGBA)
					dst.Pix[di+0] = c.R
					dst.Pix[di+1] = c.G
					dst.Pix[di+2] = c.B
					dst.Pix[di+3] = c.A
					di += 4
				}
			}
		})

	}
	return dst
}

// toNRGBA converts any image type to *image.NRGBA with min-point at (0, 0).
func toNRGBA(img image.Image) *image.NRGBA {
	srcBounds := img.Bounds()
	if srcBounds.Min.X == 0 && srcBounds.Min.Y == 0 {
		if src0, ok := img.(*image.NRGBA); ok {
			return src0
		}
	}
	return Clone(img)
}

// AdjustFunc performs a gamma correction on the image and returns the adjusted image.
func AdjustFunc(img image.Image, fn func(c color.NRGBA) color.NRGBA) *image.NRGBA {
	src := toNRGBA(img)
	width := src.Bounds().Max.X
	height := src.Bounds().Max.Y
	dst := image.NewNRGBA(image.Rect(0, 0, width, height))
	parallel(height, func(partStart, partEnd int) {
		for y := partStart; y < partEnd; y++ {
			for x := 0; x < width; x++ {
				i := y*src.Stride + x*4
				j := y*dst.Stride + x*4
				r := src.Pix[i+0]
				g := src.Pix[i+1]
				b := src.Pix[i+2]
				a := src.Pix[i+3]
				c := fn(color.NRGBA{r, g, b, a})
				dst.Pix[j+0] = c.R
				dst.Pix[j+1] = c.G
				dst.Pix[j+2] = c.B
				dst.Pix[j+3] = c.A
			}
		}
	})
	return dst
}

// AdjustGamma performs a gamma correction on the image and returns the adjusted image.
// Gamma parameter must be positive. Gamma = 1.0 gives the original image.
// Gamma less than 1.0 darkens the image and gamma greater than 1.0 lightens it.
//
// Example:
//
//	dstImage = imaging.AdjustGamma(srcImage, 0.7)
//
func AdjustGamma(img image.Image, gamma float64) *image.NRGBA {
	e := 1.0 / math.Max(gamma, 0.0001)
	lut := make([]uint8, 256)
	for i := 0; i < 256; i++ {
		lut[i] = clamp(math.Pow(float64(i)/255.0, e) * 255.0)
	}
	fn := func(c color.NRGBA) color.NRGBA {
		return color.NRGBA{lut[c.R], lut[c.G], lut[c.B], c.A}
	}
	return AdjustFunc(img, fn)
}

func sigmoid(a, b, x float64) float64 {
	return 1 / (1 + math.Exp(b*(a-x)))
}

// AdjustSigmoid changes the contrast of the image using a sigmoidal function and returns the adjusted image.
// It's a non-linear contrast change useful for photo adjustments as it preserves highlight and shadow detail.
// The midpoint parameter is the midpoint of contrast that must be between 0 and 1, typically 0.5.
// The factor parameter indicates how much to increase or decrease the contrast, typically in range (-10, 10).
// If the factor parameter is positive the image contrast is increased otherwise the contrast is decreased.
//
// Examples:
//
//	dstImage = imaging.AdjustSigmoid(srcImage, 0.5, 3.0) // increase the contrast
//	dstImage = imaging.AdjustSigmoid(srcImage, 0.5, -3.0) // decrease the contrast
//
func AdjustSigmoid(img image.Image, midpoint, factor float64) *image.NRGBA {
	if factor == 0 {
		return Clone(img)
	}
	lut := make([]uint8, 256)
	a := math.Min(math.Max(midpoint, 0.0), 1.0)
	b := math.Abs(factor)
	sig0 := sigmoid(a, b, 0)
	sig1 := sigmoid(a, b, 1)
	e := 1.0e-6
	if factor > 0 {
		for i := 0; i < 256; i++ {
			x := float64(i) / 255.0
			sigX := sigmoid(a, b, x)
			f := (sigX - sig0) / (sig1 - sig0)
			lut[i] = clamp(f * 255.0)
		}
	} else {
		for i := 0; i < 256; i++ {
			x := float64(i) / 255.0
			arg := math.Min(math.Max((sig1-sig0)*x+sig0, e), 1.0-e)
			f := a - math.Log(1.0/arg-1.0)/b
			lut[i] = clamp(f * 255.0)
		}
	}
	fn := func(c color.NRGBA) color.NRGBA {
		return color.NRGBA{lut[c.R], lut[c.G], lut[c.B], c.A}
	}
	return AdjustFunc(img, fn)
}

// AdjustContrast changes the contrast of the image using the percentage parameter and returns the adjusted image.
// The percentage must be in range (-100, 100). The percentage = 0 gives the original image.
// The percentage = -100 gives solid grey image.
//
// Examples:
//
//	dstImage = imaging.AdjustContrast(srcImage, -10) // decrease image contrast by 10%
//	dstImage = imaging.AdjustContrast(srcImage, 20) // increase image contrast by 20%
//
func AdjustContrast(img image.Image, percentage float64) *image.NRGBA {
	percentage = math.Min(math.Max(percentage, -100.0), 100.0)
	lut := make([]uint8, 256)
	v := (100.0 + percentage) / 100.0
	for i := 0; i < 256; i++ {
		if 0 <= v && v <= 1 {
			lut[i] = clamp((0.5 + (float64(i)/255.0-0.5)*v) * 255.0)
		} else if 1 < v && v < 2 {
			lut[i] = clamp((0.5 + (float64(i)/255.0-0.5)*(1/(2.0-v))) * 255.0)
		} else {
			lut[i] = uint8(float64(i)/255.0+0.5) * 255
		}
	}
	fn := func(c color.NRGBA) color.NRGBA {
		return color.NRGBA{lut[c.R], lut[c.G], lut[c.B], c.A}
	}
	return AdjustFunc(img, fn)
}

// AdjustBrightness changes the brightness of the image using the percentage parameter and returns the adjusted image.
// The percentage must be in range (-100, 100). The percentage = 0 gives the original image.
// The percentage = -100 gives solid black image. The percentage = 100 gives solid white image.
//
// Examples:
//
//	dstImage = imaging.AdjustBrightness(srcImage, -15) // decrease image brightness by 15%
//	dstImage = imaging.AdjustBrightness(srcImage, 10) // increase image brightness by 10%
//
func AdjustBrightness(img image.Image, percentage float64) *image.NRGBA {
	percentage = math.Min(math.Max(percentage, -100.0), 100.0)
	lut := make([]uint8, 256)
	shift := 255.0 * percentage / 100.0
	for i := 0; i < 256; i++ {
		lut[i] = clamp(float64(i) + shift)
	}
	fn := func(c color.NRGBA) color.NRGBA {
		return color.NRGBA{lut[c.R], lut[c.G], lut[c.B], c.A}
	}
	return AdjustFunc(img, fn)
}

// Grayscale produces grayscale version of the image.
func Grayscale(img image.Image) *image.NRGBA {
	fn := func(c color.NRGBA) color.NRGBA {
		f := 0.299*float64(c.R) + 0.587*float64(c.G) + 0.114*float64(c.B)
		y := uint8(f + 0.5)
		return color.NRGBA{y, y, y, c.A}
	}
	return AdjustFunc(img, fn)
}

// Invert produces inverted (negated) version of the image.
func Invert(img image.Image) *image.NRGBA {
	fn := func(c color.NRGBA) color.NRGBA {
		return color.NRGBA{255 - c.R, 255 - c.G, 255 - c.B, c.A}
	}
	return AdjustFunc(img, fn)
}

// parallel starts parallel image processing based on the current GOMAXPROCS value.
// If GOMAXPROCS = 1 it uses no parallelization.
// If GOMAXPROCS > 1 it spawns N=GOMAXPROCS workers in separate goroutines.
func parallel(dataSize int, fn func(partStart, partEnd int)) {
	numGoroutines := 1
	partSize := dataSize
	numProcs := runtime.GOMAXPROCS(0)
	if numProcs > 1 {
		numGoroutines = numProcs
		partSize = dataSize / (numGoroutines * 10)
		if partSize < 1 {
			partSize = 1
		}
	}
	if numGoroutines == 1 {
		fn(0, dataSize)
	} else {
		var wg sync.WaitGroup
		wg.Add(numGoroutines)
		idx := uint64(0)
		for p := 0; p < numGoroutines; p++ {
			go func() {
				defer wg.Done()
				for {
					partStart := int(atomic.AddUint64(&idx, uint64(partSize))) - partSize
					if partStart >= dataSize {
						break
					}
					partEnd := partStart + partSize
					if partEnd > dataSize {
						partEnd = dataSize
					}
					fn(partStart, partEnd)
				}
			}()
		}
		wg.Wait()
	}
}

// absint returns the absolute value of i.
func absint(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

// clamp rounds and clamps float64 value to fit into uint8.
func clamp(x float64) uint8 {
	v := int64(x + 0.5)
	if v > 255 {
		return 255
	}
	if v > 0 {
		return uint8(v)
	}
	return 0
}
