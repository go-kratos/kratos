package service

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg" // for image decode
	_ "image/png"
	"strconv"
)

// ImageScorer .
type ImageScorer struct {
	Score float32
	step  int
	img   image.Image
	gray  [][]uint8
	avg   [][]uint8
	bin   [][]uint64
}

// NewImageScorer .
func NewImageScorer(img image.Image, step int) (is *ImageScorer, err error) {
	is = &ImageScorer{
		Score: 0,
		img:   img,
		step:  step,
	}
	return
}

// CompareWithPure 与纯色图片比较hash距离
func (is *ImageScorer) CompareWithPure() {
	is.grayMatrix()
	is.averageMatrix()
	is.binaryMatrix()
	hashCode := is.matrixHash()
	var count float32
	for i := 0; i < len(hashCode); i++ {
		if hashCode[i] != byte(48) {
			count++
		}
	}
	is.Score = count / float32(len(hashCode))
}

// grayImgMatrix get gray image matrix from image.Image
func (is *ImageScorer) grayMatrix() {
	bounds := is.img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y
	dst := make([][]uint8, height)
	for i := 0; i < height; i++ {
		dst[i] = make([]uint8, width)
	}
	// 图像的边界不一定从（0,0）开始
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := is.img.At(x, y).RGBA()
			gray16 := uint32(0.3*float32(r) + 0.59*float32(g) + 0.11*float32(b))
			dst[y][x] = uint8(gray16 >> 8)
		}
	}
	is.gray = dst
}

// averageMatrix .
func (is *ImageScorer) averageMatrix() {
	dst := make([][]uint8, len(is.gray)/is.step)
	for i := 0; i < len(is.gray)/is.step; i++ {
		dst[i] = make([]uint8, len(is.gray[0])/is.step)
	}
	for i := range dst {
		for j := range dst[i] {
			avg := 0
			xhead := i * is.step
			xtail := (i + 1) * is.step
			yhead := j * is.step
			ytail := (j + 1) * is.step
			for _, v := range is.gray[xhead:xtail] {
				for _, u := range v[yhead:ytail] {
					avg += int(u)
				}
			}
			dst[i][j] = uint8(avg / is.step / is.step)
		}
	}
	is.avg = dst
}

// binaryMatrix .
func (is *ImageScorer) binaryMatrix() {
	height := len(is.avg)
	width := len(is.avg[0])
	dst := make([][]uint64, height)
	for i := range dst {
		dst[i] = make([]uint64, width)
	}
	for i, v := range is.avg {
		x := i * is.step
		for j, u := range v {
			y := j * is.step
			for k := 0; k < is.step; k++ {
				for l := 0; l < is.step; l++ {
					dst[i][j] = dst[i][j] << 1
					if u < is.gray[x+k][y+l] {
						dst[i][j]++
					}
				}
			}
		}
	}
	is.bin = dst
}

// matrixHash .
func (is *ImageScorer) matrixHash() []byte {
	var buffer bytes.Buffer
	for i := range is.bin {
		for _, v := range is.bin[i] {
			// hex := DecToHex(v)
			hex := fmt.Sprintf("%016s", strconv.FormatUint(v, 16))
			buffer.WriteString(hex)
		}
	}
	return buffer.Bytes()
}
