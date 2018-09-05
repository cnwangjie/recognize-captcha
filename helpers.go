package main

import (
	"image"
	"image/png"
	"os"

	"fmt"
)

func saveImage(img image.Image, path string) error {
	o, err := os.Create(path)
	if err != nil {
		return err
	}
	defer o.Close()
	png.Encode(o, img)
	return nil
}

// 在终端输出图形
func displayBininaryImage(data binaryImage) {
	width := len(data)
	height := len(data[0])
	for y := 0; y < height; y++ {
		row := ""
		for x := 0; x < width; x++ {
			if data[x][y] {
				row += "■"
			} else {
				row += "□"
			}
		}
		fmt.Println(row)
	}
}
