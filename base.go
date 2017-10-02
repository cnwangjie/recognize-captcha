package main

import (
	"fmt"
	p "github.com/benbjohnson/phantomjs"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"time"
)

// 处理样本
func handleSample() {
	samples, _ := ioutil.ReadDir("./" + samplePathName)
	for _, v := range samples {
		img, err := openPngImage("./" + samplePathName + "/" + v.Name())
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		img = medianFilterImage(img)
		bi := baseHandler(img)
		bi = removeArouldBlank(bi)
		bi4, err := cutLetter(bi)
		if err != nil {
			fmt.Println(v.Name() + " cut failed")
			continue
		}
		for i := 0; i < 4; i += 1 {
			letter := string(v.Name()[i])
			letterData := standardizeBI(bi4[i])
			letterDir := "./" + handledSamplePathName + "/" + letter + "/"
			files, err := ioutil.ReadDir(letterDir)
			if err != nil {
				if os.IsNotExist(err) {
					os.MkdirAll(letterDir, os.FileMode(0755))
				}
			}
			filesum := len(files)
			err = saveBI(letterData, letterDir+strconv.Itoa(filesum))
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}

// 获取样本检验结果
func catchAnCaptchaAndTest() error {
	fmt.Println("loading...")
	page, err := p.CreateWebPage()
	if err != nil {
		return err
	}
	defer page.Close()

	if err := page.Open(loginUrl); err != nil {
		return err
	}
	//if err := page.IncludeJS("https://cdn.bootcss.com/jquery/3.2.1/jquery.min.js"); err != nil {
	//	return err
	//}

	_, err = page.Evaluate(`
	function() {
		var code = document.querySelector('#icode');
		code.style.position = 'fixed';
		code.style.left = 0;
		code.style.top = 0;
	}
	`)
	time.Sleep(1 * time.Second)
	if err != nil {
		return err
	}
	if err := page.SetClipRect(p.Rect{0, 0, 72, 27}); err != nil {
		return err
	}
	if err := page.Render("tmp.png", "png", 100); err != nil {
		return err
	}
	rawimg, err := openPngImage("tmp.png")
	if err != nil {
		return err
	}
	img := medianFilterImage(rawimg)
	saveImage(rawimg, "./tmp.png")
	bi := baseHandler(img)
	displayBininaryImage(bi)
	fmt.Println("enter code:")
	code := ""

	// manually
	fmt.Scanf("%s", &code)

	_, err = page.Evaluate(`
	function() {
		document.querySelector('#txtUserName').value = '123';
		document.querySelector('#TextBox2').value = '123';
		document.querySelector('#txtSecretCode').value = '` + code + `';
		document.querySelector('#Button1').click();
	}
	`)
	if err != nil {
		return err
	}
	fmt.Println("verifying...")
	time.Sleep(2 * time.Second)
	html, err := page.Evaluate(`
	function() {
		return document.body.innerHTML;
	}
	`)
	if err != nil {
		return err
	}
	if err := page.SetClipRect(p.Rect{0, 0, 1080, 768}); err != nil {
		return err
	}
	if err := page.Render("s.png", "png", 100); err != nil {
		return err
	}
	re, _ := regexp.Compile("用户名不存在")
	REResult := re.FindString(html.(string))
	//fmt.Println(REResult)
	if REResult == "" {
		fmt.Println("code wrong")
		return nil
	}

	saveImage(rawimg, "./"+samplePathName+"/"+code+".png")
	fmt.Println("code: " + code)
	return nil
}
