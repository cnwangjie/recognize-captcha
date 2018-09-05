package main

import (
	"io"
	"strings"
	"path/filepath"
	"path"
	"io/ioutil"
	"os"
	"image/gif"
	"image/png"
	"image"
)

func openImageFromFile(file io.Reader) (image.Image, error) {
	img, err := png.Decode(file)
	if err == nil {
		return img, nil
	}
	img, err = gif.Decode(file)
	if err == nil {
		return img, nil
	}
	return nil, err
}

func openImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return openImageFromFile(file)
}

func mapDir(dir string, fn func(string)) {
	files, _ := ioutil.ReadDir(dir)
	for _, file := range files {
		fn(path.Join(dir, file.Name()))
	}
}

func loadSamples(dir string) []sample {
	files, _ := ioutil.ReadDir(dir)
	samples := make([]sample, len(files))
	for i, f := range files {
		file := path.Join(dir, f.Name())
		img, err := openImage(file)
		if err != nil {
			continue
		}
		label := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
		samples[i] = sample{img, label}
	}
	return samples
}
