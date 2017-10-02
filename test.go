package main

func main() {

	// 基本函数测试
	//img := getCaptcha()
	//saveImage(img, "./raw.png")
	//img = medianFilterImage(img)
	//saveImage(img, "./dst.png")
	//bi := baseHandler(img)
	//bi = removeArouldBlank(bi)
	//displayBininaryImage(bi)

	//sb, err := loadBI("./" + handledSamplePathName + "/0/1")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//displayBininaryImage(sb.BI())

	// 处理样本
	//handleSample()

	// 训练
	//probabilityTrain()

	// 验证
	probalityVerify()

	// 打码获取样本
	//for i := 0; i < 1E5 ; i += 1 {
	//	if err := p.DefaultProcess.Open(); err != nil {
	//		fmt.Println(err)
	//		os.Exit(1)
	//	}
	//	err := catchAnCaptchaAndTest()
	//	if err != nil {
	//		fmt.Println(err.Error())
	//	}
	//	p.DefaultProcess.Close()
	//}

}
