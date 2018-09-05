package main

import (
	"image"
	"bufio"
	"os"
	"path"
	"fmt"
	"math/rand"
	"strings"
	"io/ioutil"
)

type sample struct {
	image image.Image
	label string
}

func handleSample() {
	samples, _ := ioutil.ReadDir(samplePath)
	files := make(map[byte](*os.File))
	fmt.Println("sample count:", len(files))
	for _, sample := range samples {
		sampleName := sample.Name()
		if !strings.HasSuffix(sampleName, ".png") {
			continue
		}
		filePath := path.Join(samplePath, sampleName)
		fmt.Println(filePath)
		img, err := openImage(filePath)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		bi := createBinaryImage(img)
		bi = filter(bi)
		bi4 := split(bi)
		for i := 0; i < 4; i++ {
			letter := sampleName[i]
			letterFile, exists := files[letter]
			if !exists {
				letterFile, _ = os.OpenFile(path.Join(handledSamplePath, string(letter)), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				files[letter] = letterFile
			}
			letterFile.WriteString(bi4[i].toString() + "\n")
		}
	}
	fmt.Println("done!")
}

func loadHandledSample() dataset {
	samples, _ := ioutil.ReadDir(handledSamplePath)
	re := make([]point, 4)
	labelsum := 0
	for _, sample := range samples {
		labelsum++
		letter := sample.Name()[0]
		file, _ := os.Open(path.Join(handledSamplePath, string(letter)))
		buf := bufio.NewReader(file)
		for {
			line, _, err := buf.ReadLine()
			if err != nil {
				break
			}
			re = append(re, point{string(line), letter})
		}
	}
	return dataset{re, labelsum}
}

func splitDataset(ds dataset, d float64) (dataset, dataset) {
	trainSize := int(float64(len(ds.points)) * d)
	rand.Shuffle(len(ds.points), func (i, j int) {
		ds.points[i], ds.points[j] = ds.points[j], ds.points[i]
	})
	return dataset{ds.points[0:trainSize - 1], ds.labelsum}, dataset{ds.points[trainSize:len(ds.points)-1], ds.labelsum}
}
