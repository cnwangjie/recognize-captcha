package main

import "fmt"

type binaryImage [][]bool

func (bi binaryImage) display() {
	width := len(bi)
	height := len(bi[0])
	for y := 0; y < height; y++ {
		row := ""
		for x := 0; x < width; x++ {
			if bi[x][y] {
				row += "■"
			} else {
				row += "□"
			}
		}
		fmt.Println(row)
	}
}

func (bi binaryImage) toVector() []bool {
	width := len(bi)
	height := len(bi[0])
	size := width * height
	re := make([]bool, size)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			re[x*height+y] = bi[x][y]
		}
	}
	return re
}

func (bi binaryImage) toString() string {
	v := bi.toVector()
	l := len(v)
	r := ""
	for i := 0; i < l; i++ {
		if v[i] {
			r += "1"
		} else {
			r += "0"
		}
	}
	return r
}
