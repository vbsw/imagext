/*
 *          Copyright 2021, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *      (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package imagext provides functions for images. It's written for fast prototyping. Results may depend on format.
package imagext

import (
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
)

// Gray converts values r (red), g (green) and b (blue)
// to a value of gray and returns it.
func Gray(r, g, b uint8) uint8 {
	return uint8((uint(r)*1742 + uint(g)*5859 + uint(b)*591) >> 13)
}

// NewGray converts argument img to a new image of grayscale and returns it.
// Result differes depending on image format.
func NewGray(img image.Image) *image.Gray {
	var gray *image.Gray
	if img != nil {
		bounds := img.Bounds()
		xMin := bounds.Min.X
		xMax := bounds.Max.X
		yMin := bounds.Min.Y
		yMax := bounds.Max.Y
		width := xMax - xMin
		height := yMax - yMin
		if width*height > 0 {
			gray = image.NewGray(image.Rect(0, 0, width, height))
			switch imgStruct := img.(type) {
			case *image.RGBA:
				toGrayImageRGBA(imgStruct.Pix, imgStruct.Stride, xMin, xMax, yMin, yMax, gray)
			case *image.RGBA64:
				toGrayImageRGBA64(imgStruct.Pix, imgStruct.Stride, xMin, xMax, yMin, yMax, gray)
			case *image.Alpha:
				toGrayImageAlpha(imgStruct.Pix, imgStruct.Stride, xMin, xMax, yMin, yMax, gray)
			case *image.Alpha16:
				toGrayImageAlpha16(imgStruct.Pix, imgStruct.Stride, xMin, xMax, yMin, yMax, gray)
			case *image.CMYK:
				toGrayImageCMYK(imgStruct.Pix, imgStruct.Stride, xMin, xMax, yMin, yMax, gray)
			case *image.Gray:
				toGrayImageGray(imgStruct.Pix, imgStruct.Stride, xMin, xMax, yMin, yMax, gray)
			case *image.Gray16:
				toGrayImageGray16(imgStruct.Pix, imgStruct.Stride, xMin, xMax, yMin, yMax, gray)
			case *image.NRGBA:
				toGrayImageNRGBA(imgStruct.Pix, imgStruct.Stride, xMin, xMax, yMin, yMax, gray)
			case *image.NRGBA64:
				toGrayImageNRGBA64(imgStruct.Pix, imgStruct.Stride, xMin, xMax, yMin, yMax, gray)
			case *image.Paletted:
				toGrayImagePaletted(imgStruct.Pix, imgStruct.Palette, imgStruct.Stride, xMin, xMax, yMin, yMax, gray)
			default:
				toGrayImageGeneric(img, xMin, xMax, yMin, yMax, gray)
			}
		} else {
			gray = image.NewGray(image.Rect(0, 0, 0, 0))
		}
	} else {
		gray = image.NewGray(image.Rect(0, 0, 0, 0))
	}
	return gray
}

// LoadImage reads image from file and returns it.
func LoadImage(path string) image.Image {
	var img image.Image
	if len(path) > 0 {
		file, err := os.Open(path)
		if err == nil {
			defer file.Close()
			ext := filepath.Ext(path)
			if ext == ".jpg" || ext == ".jpeg" {
				img, _ = jpeg.Decode(file)
			} else if ext == ".png" || ext == ".apng" {
				img, _ = png.Decode(file)
			} else if ext == ".gif" {
				img, _ = gif.Decode(file)
			} else {
				img, _, _ = image.Decode(file)
			}
		}
	}
	return img
}

// SaveImage saves image to file. Format is recognized from
// extension in path. Default is PNG.
func SaveImage(path string, img image.Image) error {
	if len(path) > 0 {
		file, err := os.Create(path)
		if err == nil {
			defer file.Close()
			ext := filepath.Ext(path)
			if ext == ".jpg" || ext == ".jpeg" {
				opt := jpeg.Options{100}
				err = jpeg.Encode(file, img, &opt)
			} else if ext == ".gif" {
				err = gif.Encode(file, img, nil)
			} else {
				err = png.Encode(file, img)
			}
		}
		return err
	}
	return nil
}

// ToMonochrome convertes image to black and white.
// Higher threshold means darker image.
func ToMonochrome(img *image.Gray, threshold uint8) {
	if (img.Rect.Max.X-img.Rect.Min.X)*(img.Rect.Max.Y-img.Rect.Min.Y) > 0 {
		for i, gray := range img.Pix {
			if gray < threshold {
				img.Pix[i] = 0
			} else {
				img.Pix[i] = 255
			}
		}
	}
}

// ToMedian sets median values for each pixel in its
// area size*size. Median value in {9, 5, 17} is 9.
func ToMedian(img *image.Gray, size uint) {
	if (img.Rect.Max.X-img.Rect.Min.X)*(img.Rect.Max.Y-img.Rect.Min.Y) > 0 && size > 1 {
		lenImg := img.Rect.Max.X - img.Rect.Min.X
		hist := make([]uint8, 256, 256)
		histZero := make([]uint8, 256, 256)
		offLines := int(size) / 2
		lines := newLines(img, int(size))
		limit := img.Rect.Max.Y - offLines
		idxLastLine := len(lines) - 1
		for y := img.Rect.Min.Y; y < limit; y++ {
			offImg := (y - img.Rect.Min.Y) * img.Stride
			offImgNew := offImg + offLines*img.Stride
			copy(lines[idxLastLine][offLines:], img.Pix[offImgNew:offImgNew+lenImg:offImgNew+lenImg])
			shiftLines(lines, lines[idxLastLine])
			for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
				fillHistogram(hist, histZero, lines, int(size), x)
				img.Pix[offImg+x-img.Rect.Min.X] = median(hist)
			}
		}
		whiteLine := lines[idxLastLine]
		setArrayValues(whiteLine, 255)
		for y, x := limit, img.Rect.Min.X; x < img.Rect.Max.X && y < img.Rect.Max.Y; x, y = x+1, y+1 {
			offImg := (y - img.Rect.Min.Y) * img.Stride
			shiftLines(lines, whiteLine)
			fillHistogram(hist, histZero, lines, int(size), x)
			img.Pix[offImg+x-img.Rect.Min.X] = median(hist)
		}
	}
}

// ToAvarage sets avarage values for each pixel in its
// area size*size. Avarage value of {9, 5, 16} is 10.
func ToAvarage(img *image.Gray, size uint) {
	if (img.Rect.Max.X-img.Rect.Min.X)*(img.Rect.Max.Y-img.Rect.Min.Y) > 0 && size > 1 {
		lenImg := img.Rect.Max.X - img.Rect.Min.X
		hist := make([]uint8, 256, 256)
		histZero := make([]uint8, 256, 256)
		offLines := int(size) / 2
		lines := newLines(img, int(size))
		limit := img.Rect.Max.Y - offLines
		idxLastLine := len(lines) - 1
		for y := img.Rect.Min.Y; y < limit; y++ {
			offImg := (y - img.Rect.Min.Y) * img.Stride
			offImgNew := offImg + offLines*img.Stride
			copy(lines[idxLastLine][offLines:], img.Pix[offImgNew:offImgNew+lenImg:offImgNew+lenImg])
			shiftLines(lines, lines[idxLastLine])
			for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
				fillHistogram(hist, histZero, lines, int(size), x)
				img.Pix[offImg+x-img.Rect.Min.X] = avarage(hist, int(size))
			}
		}
		whiteLine := lines[idxLastLine]
		setArrayValues(whiteLine, 255)
		for y := limit; y < img.Rect.Max.Y; y++ {
			offImg := (y - img.Rect.Min.Y) * img.Stride
			shiftLines(lines, whiteLine)
			for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
				fillHistogram(hist, histZero, lines, int(size), x)
				img.Pix[offImg+x-img.Rect.Min.X] = avarage(hist, int(size))
			}
		}
	}
}

func setArrayValues(array []uint8, value uint8) {
	for i := range array {
		array[i] = value
	}
}

func newLines(img *image.Gray, size int) [][]uint8 {
	lines := make([][]uint8, size, size)
	lenImg := img.Rect.Max.X - img.Rect.Min.X
	lenLine := lenImg + size - 1
	offLines := size / 2
	for i := range lines {
		lines[i] = make([]uint8, lenLine, lenLine)
		if i == 0 {
			setArrayValues(lines[0], 255)
		} else {
			copy(lines[i], lines[0])
		}
	}
	for i, y := offLines-1, img.Rect.Min.Y; i >= 0 && y < img.Rect.Max.Y; i, y = i-1, y+1 {
		offImg := (y - img.Rect.Min.Y) * img.Stride
		copy(lines[i][offLines:], img.Pix[offImg:offImg+lenImg:offImg+lenImg])
	}
	return lines
}

func shiftLines(lines [][]uint8, line0 []uint8) {
	copy(lines[1:], lines[0:len(lines)-1])
	lines[0] = line0
}

func fillHistogram(hist, histZero []uint8, lines [][]uint8, size, x int) {
	copy(hist, histZero)
	for _, line := range lines {
		for j := 0; j < size; j++ {
			idx := line[x+j]
			hist[idx]++
		}
	}
}

func sum(values []uint8) uint {
	var sum uint
	for _, v := range values {
		sum += uint(v)
	}
	return sum
}

func median(hist []uint8) uint8 {
	var leftSumPrev uint
	var rightSumPrev uint
	left := 0
	right := len(hist)
	for left < right {
		middle := (left + right) / 2
		if left != middle {
			leftSum := sum(hist[left:middle:middle])
			rightSum := sum(hist[middle:right:right])
			if leftSumPrev+leftSum > rightSum+rightSumPrev {
				right = middle
				rightSumPrev += rightSum
			} else if leftSumPrev+leftSum < rightSum+rightSumPrev {
				left = middle
				leftSumPrev += leftSum
			} else {
				left = middle
			}
		} else {
			return uint8(left)
		}
	}
	return 0
}

func avarage(hist []uint8, size int) uint8 {
	var sum uint
	for i, v := range hist {
		sum += (uint(i) * uint(v))
	}
	sum /= uint(size * size)
	return uint8(sum)
}

func cmykToGray(c, m, y, k uint) uint8 {
	kDiff := 255 - k
	r := ((k * c >> 8) + kDiff) - c
	g := ((k * m >> 8) + kDiff) - m
	b := ((k * y >> 8) + kDiff) - y
	// 0.299 * R + 0.587 * G + 0.114 * B
	return uint8((uint(r)*2449 + uint(g)*4809 + uint(b)*934) >> 13)
}

func toGrayImageRGBA(pix []uint8, stride, xMin, xMax, yMin, yMax int, gray *image.Gray) {
	i := 0
	for y := yMin; y < yMax; y++ {
		for x := xMin; x < xMax; x++ {
			offset := (y-yMin)*stride + (x-xMin)*4
			r := pix[offset]
			g := pix[offset+1]
			b := pix[offset+2]
			gray.Pix[i] = Gray(r, g, b)
			i++
		}
	}
}

func toGrayImageRGBA64(pix []uint8, stride, xMin, xMax, yMin, yMax int, gray *image.Gray) {
	i := 0
	for y := yMin; y < yMax; y++ {
		for x := xMin; x < xMax; x++ {
			offset := (y-yMin)*stride + (x-xMin)*8
			r := pix[offset+1]
			g := pix[offset+3]
			b := pix[offset+5]
			gray.Pix[i] = Gray(r, g, b)
			i++
		}
	}
}

func toGrayImageAlpha(pix []uint8, stride, xMin, xMax, yMin, yMax int, gray *image.Gray) {
	i := 0
	for y := yMin; y < yMax; y++ {
		for x := xMin; x < xMax; x++ {
			offset := (y-yMin)*stride + (x-xMin)*1
			a := pix[offset]
			gray.Pix[i] = 255 - a
			i++
		}
	}
}

func toGrayImageAlpha16(pix []uint8, stride, xMin, xMax, yMin, yMax int, gray *image.Gray) {
	i := 0
	for y := yMin; y < yMax; y++ {
		for x := xMin; x < xMax; x++ {
			offset := (y-yMin)*stride + (x-xMin)*2
			a := pix[offset+1]
			gray.Pix[i] = 255 - a
			i++
		}
	}
}

func toGrayImageCMYK(pix []uint8, stride, xMin, xMax, yMin, yMax int, gray *image.Gray) {
	i := 0
	for y := yMin; y < yMax; y++ {
		for x := xMin; x < xMax; x++ {
			offset := (y-yMin)*stride + (x-xMin)*4
			c := pix[offset]
			m := pix[offset+1]
			z := pix[offset+2]
			k := pix[offset+3]
			gray.Pix[i] = cmykToGray(uint(c), uint(m), uint(z), uint(k))
			i++
		}
	}
}

func toGrayImageGray(pix []uint8, stride, xMin, xMax, yMin, yMax int, gray *image.Gray) {
	length := xMax - xMin
	i := 0
	for y := yMin; y < yMax; y++ {
		offset := (y-yMin)*stride + xMin
		copy(gray.Pix[i:i+length:i+length], pix[offset:offset+length:offset+length])
		i += length
	}
}

func toGrayImageGray16(pix []uint8, stride, xMin, xMax, yMin, yMax int, gray *image.Gray) {
	i := 0
	for y := yMin; y < yMax; y++ {
		for x := xMin; x < xMax; x++ {
			offset := (y-yMin)*stride + (x-xMin)*2
			gray.Pix[i] = pix[offset+1]
			i++
		}
	}
}

func toGrayImageNRGBA(pix []uint8, stride, xMin, xMax, yMin, yMax int, gray *image.Gray) {
	i := 0
	for y := yMin; y < yMax; y++ {
		for x := xMin; x < xMax; x++ {
			offset := (y-yMin)*stride + (x-xMin)*4
			r := pix[offset]
			g := pix[offset+1]
			b := pix[offset+2]
			gray.Pix[i] = Gray(r, g, b)
			i++
		}
	}
}

func toGrayImageNRGBA64(pix []uint8, stride, xMin, xMax, yMin, yMax int, gray *image.Gray) {
	i := 0
	for y := yMin; y < yMax; y++ {
		for x := xMin; x < xMax; x++ {
			offset := (y-yMin)*stride + (x-xMin)*8
			r := pix[offset+1]
			g := pix[offset+3]
			b := pix[offset+5]
			gray.Pix[i] = Gray(r, g, b)
			i++
		}
	}
}

func toGrayImagePaletted(pix []uint8, palette []color.Color, stride, xMin, xMax, yMin, yMax int, gray *image.Gray) {
	i := 0
	for y := yMin; y < yMax; y++ {
		for x := xMin; x < xMax; x++ {
			offset := (y-yMin)*stride + (x-xMin)*1
			grayColor := color.GrayModel.Convert(palette[offset])
			r, _, _, _ := grayColor.RGBA()
			gray.Pix[i] = uint8(r >> 8)
			i++
		}
	}
}

func toGrayImageGeneric(img image.Image, xMin, xMax, yMin, yMax int, gray *image.Gray) {
	i := 0
	for y := yMin; y < yMax; y++ {
		for x := xMin; x < xMax; x++ {
			grayColor := color.GrayModel.Convert(img.At(x, y))
			r, _, _, _ := grayColor.RGBA()
			gray.Pix[i] = uint8(r >> 8)
			i++
		}
	}
}
