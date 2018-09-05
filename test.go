package main

import (
	"time"
	"fmt"
)

const (
	testSamplesPath = "/home/wangjie/Gopath/src/rc/sample"
)

func singleTest() {
	dataset := loadHandledSample()
	trainSet, testSet := splitDataset(dataset, 0.7)
	testSetSize := len(testSet.points)
	correct := 0
	start := time.Now()
	fmt.Println("train num:", len(trainSet.points), "test num:", testSetSize)
	for i := 0; i < testSetSize; i++ {
		fmt.Println("(", i, "/", testSetSize, ")", correct)
		if predict(trainSet, testSet.points[i].data) == testSet.points[i].label {
			correct++
		}
	}
	fmt.Println("correct:", correct, "(", float64(correct) / float64(testSetSize), ")")
	fmt.Println("time spent:", time.Now().Sub(start).String())
}

func test() {
	samples := loadSamples(testSamplesPath)
	dataset := loadHandledSample()
	correct := 0
	start := time.Now()
	fmt.Println("train num:", len(dataset.points), "test num:", len(samples))
	for i, sample := range samples {
		bi := createBinaryImage(sample.image)
		bi = filter(bi)
		bi4 := split(bi)
		re := ""
		for _, bi := range bi4 {
			re += string(predict(dataset, bi.toString()))
		}
		if re == sample.label {
			correct++
		}
		fmt.Println(i, correct)
	}
	fmt.Println("correct:", correct, "(", float64(correct) / float64(len(samples)), ")")
	fmt.Println("time spent:", time.Now().Sub(start).String())
}
