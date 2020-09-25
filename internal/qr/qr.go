package qr

import (
	"bytes"
	"fmt"
	qrcode "github.com/skip2/go-qrcode"
	"image"
	"image/color"
	"image/png"
)

func Generate(id string)(image.Image, error){
	//generate qt
	raw, err := qrcode.Encode(id, qrcode.Highest, 256)
	if err != nil{
		return nil, fmt.Errorf("Failed to generate qr-code, err :%v", err)
	}
	// convert bytes to io.Reader
	reader := bytes.NewReader(raw)
	//	generate png
	img, err := png.Decode(reader)
	if err != nil{
		return nil, fmt.Errorf("failed to create png image, err: %v", err)
	}
	return img, nil
}

func ConsolePngWriter(img image.Image){
	levels := []string{" ", "░", "▒", "▓", "█"}

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			c := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			level := c.Y / 51 // 51 * 5 = 255
			if level == 5 {
				level--
			}
			fmt.Print(levels[level])
		}
		fmt.Print("\n")
	}
}