package main

import (
	"image"
	"fmt"
	"regexp"
	"time"

	p "github.com/benbjohnson/phantomjs"
)

const (
	loginURL = "http://xk1.ahu.cn/default2.aspx"
)

func manually(img image.Image) string {
	fmt.Println("enter code:")
	var code string
	fmt.Scanf("%s", &code)
	return code
}

func saveIfSuccess(img image.Image, code string) {
	saveImage(img, "./sample/" + code + ".png")
	fmt.Println("correct:", code)
}

func getCaptchaAndTest(recgImg func(image.Image) string, success func(image.Image, string)) error {
	fmt.Println("loading...")
	page, err := p.CreateWebPage()
	if err != nil {
		return err
	}
	defer page.Close()

	if err := page.Open(loginURL); err != nil {
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
	if err := page.SetClipRect(p.Rect{Top: 0, Left: 0, Width: 72, Height: 27}); err != nil {
		return err
	}
	if err := page.Render("tmp.png", "png", 100); err != nil {
		return err
	}
	rawimg, err := openImage("tmp.png")
	if err != nil {
		return err
	}

	code := recgImg(rawimg)

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
	time.Sleep(2 * time.Second)
	html, err := page.Evaluate(`
	function() {
		return document.body.innerHTML;
	}
	`)
	if err != nil {
		return err
	}
	re, _ := regexp.Compile("用户名不存在")
	REResult := re.FindString(html.(string))
	if REResult == "" {
		return nil
	}
	if success != nil {
		success(rawimg, code)
	}
	return nil
}

func maunallyGetCode() {
	for {
		p.DefaultProcess.Open()
		err := getCaptchaAndTest(manually, saveIfSuccess)
		if err != nil {
			fmt.Println(err.Error())
		}
		p.DefaultProcess.Close()
	}
}
