package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main()  {
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	Logfile, _ := os.Create("opendoor.log")
	gin.DefaultWriter = io.MultiWriter(Logfile)
	router := gin.Default()
	router.LoadHTMLGlob("template/*")
	router.GET("/getpic", getting)
	router.StaticFile("/qr.png", "qr.png")
	router.Run(":7788")
}

func getting(context *gin.Context) {
	run()
	context.HTML(http.StatusOK, "index.html", gin.H{
		"pic" : "qr.png",
	})
}

func run()  {
	login_cl := &http.Client{}
	login_val := make(url.Values)
	login_val.Set("iccode", "app_login")
	login_val.Set("phone", "15065362098")
	login_val.Set("pwd", "fa53d132a53d33931c51597fcfec0a74")
	login_values := login_val.Encode()
	req_login, err := http.NewRequest("POST", "http://epay.sdzy.cn:7110/school.ashx", strings.NewReader(login_values))

	if err != nil {
		log.Fatal("POST构建失败")
	}
	req_login.Header.Set("Cookie", "ASP.NET_SessionId=5ositwhlome30kbyxjmfrekq")
	req_login.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req_login.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 7.1.2; vivo X9 Build/N2G47H; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/55.0.2883.91 Mobile Safari/537.36")

	login_cl.Do(req_login)

	defer req_login.Body.Close()

	client := &http.Client{}
	v := make(url.Values)
	v.Set("iccode", "app_barcode")
	v.Set("studentno", "201710143052")
	v.Set("bartype", "0")
	v.Set("patroncode", "10000000")
	valuse := v.Encode()
	req, err := http.NewRequest("POST", "http://epay.sdzy.cn:7110/school.ashx", strings.NewReader(valuse))

	if err != nil {
		log.Fatal("POST构建失败")
	}

	req.Header.Set("Cookie", "ASP.NET_SessionId=5ositwhlome30kbyxjmfrekq")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 7.1.2; vivo X9 Build/N2G47H; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/55.0.2883.91 Mobile Safari/537.36")

	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Line:55")
	}
	type Mes struct {
		Barcode string
		Answercode string
		Dscrp string
	}
	var s Mes
	err = json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("Line:65")
	}
	fmt.Println(s.Barcode)
	err = qrcode.WriteFile(s.Barcode, qrcode.Medium, 256, "qr.png")
	if err != nil {
		fmt.Println("qr write error")
	}
}
