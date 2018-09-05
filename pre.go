package main

import (
	"image"
)

const (
	thresholdR = 28270
	thresholdG = 28270
	thresholdB = 25700
)

func createBinaryImage(img image.Image) binaryImage {
	width := img.Bounds().Size().X
	height := img.Bounds().Size().Y
	re := make([][]bool, width)
	for x := 0; x < width; x++ {
		re[x] = make([]bool, height)
		for y := 0; y < height; y++ {
			r, g, b, _ := img.At(x, y).RGBA()
			re[x][y] = r < 28270 && g < 28270 && b > 25700
		}
	}
	return re
}

func filter(bi binaryImage) binaryImage {
	width := len(bi)
	height := len(bi[0])
	re := make([][]bool, width)
	for x := 0; x < width; x++ {
		re[x] = make([]bool, height)
		for y := 0; y < height; y++ {
			if !bi[x][y] {
				continue
			}
			num := 0
			for dx := -1; dx < 2; dx++ {
				for dy := -1; dy < 2; dy++ {
					if !(dx == 0 && dy == 0) && x + dx > -1 && x + dx < width && y + dy > -1 && y + dy < height && bi[x + dx][y + dy] {
						num++
					}
				}
			}
			re[x][y] = num > 2
			// i := []int{1, -1, 0, 0}
			// for d := 0; d < 4; d++ {
			// 	if x + i[d] > -1 && x + i[d] < width && y + i[3-d] > -1 && y + i[3-d] < height && bi[x + i[d]][y + i[3-d]] {
			// 		num++
			// 	}
			// }
			// re[x][y] = num > 0
		}
	}
	return re
}

func split(bi binaryImage) [4]binaryImage {
	width := len(bi)
	height := len(bi[0])
	col := make([]int, width)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if bi[x][y] {
				col[x]++
			}
		}
	}
	span := [4][2]int{}
	cur := 0
	for l := 0; l < width - 1; l++ {
		if col[l] == 0 {
			continue
		}
		for r := l + 1; r < width; r++ {
			if cur == 4 {
				break
			}
			if r - l > 10 {
				end := true
				for d := 1; d < 5 && d < width; d++ {
					if col[r + d] <= 1 {
						end = false
						break
					}
				}
				if end {
					span[cur] = [2]int{l, r}
					cur++
					l = r
					break
				}
			}
			if r - l > 9 && col[r] < 3 || r - l > 2 && col[r] == 0 {
				span[cur] = [2]int{l, r}
				cur++
				l = r
				break
			}
		}
	}
	re := [4]binaryImage{}
	for i := 0; i < 4; i++ {
		min, max := height, 0
		for y := 0; y < height; y++ {
			n := 0
			for x := span[i][0]; x < span[i][1] + 1; x++ {
				if bi[x][y] {
					n++
				}
			}
			if n != 0 {
				if y > max {
					max = y
				}
				if y < min {
					min = y
				}
			}
		}
		subWidth := span[i][1] - span[i][0] + 1
		re[i] = make(binaryImage, subWidth)
		for x := 0; x < subWidth; x++ {
			re[i][x] = bi[x + span[i][0]][min:max+1]
		}
	}
	return re
}
