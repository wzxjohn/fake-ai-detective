package main

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"math/rand"

	"github.com/dchest/captcha"
)

func generateCaptchaImage(traceID string) ([]byte, error) {
	img := captcha.NewImage(traceID, captcha.RandomDigits(6), 240, 80)
	buf := new(bytes.Buffer)
	_, err := img.WriteTo(buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func generateImage() ([]byte, error) {
	// 生成随机颜色的图片
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	color := color.RGBA{
		R: uint8(rand.Intn(256)),
		G: uint8(rand.Intn(256)),
		B: uint8(rand.Intn(256)),
		A: 255,
	}
	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			img.Set(x, y, color)
		}
	}

	// 将图片编码为 WebP
	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, img, &jpeg.Options{Quality: 80}); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
