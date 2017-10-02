package main

import (
	"image"
	"image/gif"
	"image/png"
	"net/http"
	"os"

	"encoding/binary"
	"errors"
	"fmt"
	"github.com/disintegration/gift"
)

type BinaryImage [][]bool
type ColumnOfBI []bool
type StandardBI [standardSize]bool

const (
	captchaUrl            = "http://xk1.ahu.cn/CheckCode.aspx"
	loginUrl              = "http://xk1.ahu.cn/default2.aspx"
	threshold             = 45000
	standardWidth         = 12
	standardHeight        = 25
	standardSize          = standardHeight * standardWidth
	samplePathName        = "sample"
	handledSamplePathName = "handledSample"
)

// 下载验证码
func getCaptcha() image.Image {
	res, _ := http.Get(captchaUrl)
	img, _ := gif.Decode(res.Body)
	return img
}

// 保存图片
func saveImage(img image.Image, path string) error {
	o, err := os.Create(path)
	if err != nil {
		return err
	}
	defer o.Close()
	png.Encode(o, img)
	return nil
}

// 基本处理
func baseHandler(src image.Image) BinaryImage {
	srcRect := src.Bounds()
	width := srcRect.Dx()
	height := srcRect.Dy()

	// 象征该点的二维bool数组
	result := make(BinaryImage, width)
	for x := 0; x < width; x += 1 {
		column := make(ColumnOfBI, height)
		for y := 0; y < height; y += 1 {
			_, _, b, _ := src.At(x, y).RGBA()
			if b > uint32(threshold) {
				column[y] = false
			} else {
				column[y] = true
			}
		}
		result[x] = column
	}

	return cutBI(result, 1, 1, width-1, height-1)
}

// 在终端输出图形
func displayBininaryImage(data BinaryImage) {
	width := len(data)
	height := len(data[0])
	for y := 0; y < height; y += 1 {
		row := ""
		for x := 0; x < width; x += 1 {
			if data[x][y] {
				row += "■"
			} else {
				row += "□"
			}
		}
		fmt.Println(row)
	}
}

// 滤波
func medianFilterImage(src image.Image) image.Image {
	g := gift.New(
		gift.UnsharpMask(1, 1, 0),
		gift.Median(3, false),
	)
	dst := image.NewRGBA(g.Bounds(src.Bounds()))
	g.Draw(dst, src)
	return dst
}

// 去除四周空白
func removeArouldBlank(bi BinaryImage) BinaryImage {
	width := len(bi)
	height := len(bi[0])
	minX, minY, maxX, maxY := width, height, 0, 0
	for x := 0; x < width; x += 1 {
		for y := 0; y < height; y += 1 {
			if bi[x][y] {
				if x > maxX {
					maxX = x
				}

				if y > maxY {
					maxY = y
				}

				if x < minX {
					minX = x
				}

				if y < minY {
					minY = y
				}
			}
		}
	}
	return cutBI(bi, minX, minY, maxX, maxY)
}

// 裁剪
func cutBI(bi BinaryImage, minX, minY, maxX, maxY int) BinaryImage {
	newWidth := maxX - minX + 1
	newBI := make(BinaryImage, newWidth)
	for x := 0; x < newWidth; x += 1 {
		newBI[x] = bi[minX+x][minY : maxY+1]
	}
	return newBI
}

// 转成标准大小的图像
func standardizeBI(bi BinaryImage) StandardBI {
	var result = StandardBI{}
	bi = removeArouldBlank(bi)
	width := len(bi)
	for i := 0; i < standardWidth; i += 1 {
		if i < width {
			for j := 0; j < standardHeight; j += 1 {
				if j < len(bi[i]) {
					result[i*standardWidth+j] = bi[i][j]
				}
			}
		}
	}
	return result
}

// 切分字母
func cutLetter(bi BinaryImage) ([4]BinaryImage, error) {
	width := len(bi)
	height := len(bi[0])
	columnPixs := make([]int, width)
	for x := 0; x < width; x += 1 {
		for y := 0; y < height; y += 1 {
			if bi[x][y] {
				columnPixs[x] += 1
			}
		}
	}
	span := [4][2]int{}
	cur := 0
	for l := 0; l < width; l += 1 {
		if columnPixs[l] == 0 {
			continue
		}
		for r := l; r < width+1; r += 1 {
			if cur > 3 {
				break
			}
			if r == width || columnPixs[r] == 0 {
				if r-l < 3 {

				} else if r-l < 15 {
					span[cur] = [2]int{l, r - 1}
					cur += 1
				} else {
					minP := height
					m := 0
					for p := (r+l)/2 - 3; p < (r+l)/2+3; p += 1 {
						if p > width || p < 0 {
							continue
						}

						if columnPixs[p] < minP {
							minP = columnPixs[p]
							m = p
						}
					}
					if cur+1 < 4 {
						span[cur] = [2]int{l, m}
						span[cur+1] = [2]int{m, r - 1}
						cur += 2
					}
				}
				l = r + 1
				continue
			}
		}
	}
	fmt.Println(cur)
	var result = [4]BinaryImage{}
	if cur != 4 {
		return result, errors.New("cut failed")
	}

	for i := 0; i < 4; i += 1 {
		result[i] = cutBI(bi, span[i][0], 0, span[i][1], height)
	}
	return result, nil
}

// 打开图片
func openPngImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// 图片相似度
func similarity(a, b BinaryImage) int {
	width := intMin(len(a), len(b))
	height := intMin(len(a[0]), len(b[0]))
	result := 0
	for x := 0; x < width; x += 1 {
		for y := 0; y < height; y += 1 {
			if a[x][y] == b[x][y] {
				result += 1
			}
		}
	}
	return result
}

// 二进制保存图片
func saveBI(bi StandardBI, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := binary.Write(file, binary.LittleEndian, bi); err != nil {
		return err
	}
	return nil
}

// 读取二进制图片
func loadBI(path string) (StandardBI, error) {
	file, err := os.Open(path)
	var data = make([]bool, standardSize)
	var result = StandardBI{}
	if err != nil {
		return result, err
	}
	defer file.Close()
	if err := binary.Read(file, binary.LittleEndian, data); err != nil {
		return result, err
	}
	for i := 0; i < standardSize; i += 1 {
		result[i] = data[i]
	}
	return result, nil
}

func (sb *StandardBI) BI() BinaryImage {
	bi := make(BinaryImage, standardWidth)
	for i := 0; i < standardWidth; i += 1 {
		bi[i] = sb[(i * standardWidth):((i+1)*standardWidth - 1)]
	}
	return bi
}

func intMin(a, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}
