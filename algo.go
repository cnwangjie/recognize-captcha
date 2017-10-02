package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

var (
	probDat = map[string][standardSize]float64{}
)

// 像素匹配度
func probabilityTrain() {
	letters, _ := ioutil.ReadDir("./" + handledSamplePathName)
	result := map[string][standardSize]float64{}
	for _, v := range letters {
		letter := v.Name()
		samples, _ := ioutil.ReadDir("./" + handledSamplePathName + "/" + letter)
		sum := len(samples)
		letterRe := [standardSize]float64{}
		data := make([]StandardBI, sum)
		for i := 0; i < sum; i += 1 {
			data[i], _ = loadBI("./" + handledSamplePathName + "/" + letter + "/" + strconv.Itoa(i))
			displayBininaryImage(data[i].BI())
			for j := 0; j < standardHeight*standardWidth; j += 1 {
				if data[i][j] {
					letterRe[j] += 1
				}
			}
		}

		for j := 0; j < standardSize; j += 1 {
			letterRe[j] /= float64(sum)
		}
		result[letter] = letterRe
	}
	file, _ := os.Create("./data/prob.json")
	defer file.Close()
	enc := json.NewEncoder(file)
	enc.Encode(result)
}

func probabiltyCal(sb StandardBI) (string, float64) {
	if len(probDat) == 0 {
		dat, _ := os.Open("./data/prob.json")
		jsonStr, _ := ioutil.ReadAll(dat)
		json.Unmarshal([]byte(jsonStr), &probDat)
	}
	probability := map[string]float64{}
	for k, v := range probDat {
		for i := 0; i < standardSize; i += 1 {
			if sb[i] {
				probability[k] += v[i]
			}
		}
	}

	var maxProb = float64(0)
	maxProbLetter := ""
	for k, v := range probability {
		if v > maxProb {
			maxProbLetter = k
			maxProb = v
		}
	}

	var absoluteProb = float64(0)
	for _, v := range probDat[maxProbLetter] {
		absoluteProb += v
	}
	return maxProbLetter, maxProb / absoluteProb
}

func probalityVerify() {
	letters, _ := ioutil.ReadDir("./" + handledSamplePathName)
	sumInAll := 0
	correctInAll := 0
	for _, v := range letters {
		letter := v.Name()
		samples, _ := ioutil.ReadDir("./" + handledSamplePathName + "/" + letter)
		sum := len(samples)
		for i := 0; i < sum; i += 1 {
			sb, _ := loadBI("./" + handledSamplePathName + "/" + letter + "/" + strconv.Itoa(i))
			recgRe, prob := probabiltyCal(sb)
			sumInAll += 1
			if recgRe == letter {
				correctInAll += 1
			}
			fmt.Println(letter+":"+recgRe, prob)
		}
	}

	fmt.Println("done!")
	fmt.Println("sum:", sumInAll)
	fmt.Println("accuracy:", float64(correctInAll)/float64(sumInAll))
}
