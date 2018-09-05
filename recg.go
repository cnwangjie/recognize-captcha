package main

import (
	"fmt"
	"image"
)

func recognize(img image.Image, ds dataset) string {
	bi := createBinaryImage(img)
	bi = filter(bi)
	bi4 := split(bi)
	re := ""
	for i := 0; i < 4; i++ {
		re += string(predict(ds, bi4[i].toString()))
	}
	return re
}

func recognizeFileAndPrint() {
	img, _ := openImage(recgFilePath)
	dataset := loadHandledSample()
	code := recognize(img, dataset)
	fmt.Println("result: {{" + code + "}}")
}
