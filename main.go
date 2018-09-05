package main

import (
	"strconv"
	"strings"
	"os"
	"fmt"
)

var (
	samplePath = "./sample"
	handledSamplePath = "./handledSample"
	sampleLimit = 2000
	recgFilePath string
)

func handleArgs() {
	args := os.Args
	argsLen := len(args)
	params := make(map[string]string)
	for i := 0; i < argsLen; i++ {
		if strings.HasPrefix(args[i], "--") {
			kv := strings.Split(strings.TrimPrefix(args[i], "--"), "=")
			params[kv[0]] = kv[1]
		} else if strings.HasPrefix(args[i], "-") {
			k := strings.TrimPrefix(args[i], "-")
			params[k] = args[i + 1]
		}
	}
	for k, v := range params {
		switch k {
		case "f", "file": recgFilePath = v
		case "l", "limit": sampleLimit, _ = strconv.Atoi(v)
		case "s", "sample": samplePath = v
		case "h", "handled": handledSamplePath = v
		default:
			fmt.Println("invalid arg:", k)
			os.Exit(2)
		}
	}
}

func main() {
	var cmd string
	if len(os.Args) < 2 {
		cmd = "test"
	} else {
		cmd = os.Args[1]
	}
	handleArgs()
	switch cmd {
	case "test": test()
	case "singleTest": singleTest()
	case "handleSample": handleSample()
	case "recgonize": recognizeFileAndPrint()
	case "serve": startServer()
	case "getCode": maunallyGetCode()
	default: fmt.Println("invalid command:", cmd)
	}
}
