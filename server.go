package main

import (
	"fmt"
	"net/http"
)

var localDataSet dataset

func requestHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.Write([]byte("hi"))
		return
	}
	req.ParseMultipartForm(10 << 10)
	_, fileHeader, _ := req.FormFile("img")
	if fileHeader.Header.Get("content-type") != "image/png" {
		res.WriteHeader(400)
	}
	file, _ := fileHeader.Open()
	img, _ := openImageFromFile(file)
	code := recognize(img, localDataSet)
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	res.Write([]byte("result: {{" + code + "}}"))
}

func startServer() {
	localDataSet = loadHandledSample()
	http.HandleFunc("/", requestHandler)
	channel := make(chan error)
	go func() {
		err := http.ListenAndServe(":3010", nil)
		channel <- err
	}()
	fmt.Println("server started")
	fmt.Println(<-channel)
}
