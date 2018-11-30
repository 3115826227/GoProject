package main

import (
	"net/http"
	"io/ioutil"
	"strings"
	"github.com/robfig/cron"
	"log"
	"os"
	"time"
	"encoding/json"
)

type RespData struct {
	Status    int         `json:"status"`
	Data      interface{} `json:"data"`
	ErrorCode interface{} `json:"errorCode"`
}

const (
	contentType = "application/x-www-form-urlencoded"
	userAgent   = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 Safari/537.36"
)

var logger *log.Logger

//创建日志
func InitLog() {
	file, err := os.OpenFile("Logger.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	logger = log.New(file, "", log.LstdFlags|log.Llongfile)
}

//发送请求
func HttpRequest(method string) (res []byte, err error) {
	c := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(RespBody))
	if err != nil {
		logger.Println(err)
		return
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Cookie", cookie)
	resp, err := c.Do(req)
	if err != nil {
		logger.Println(err)
		return
	}
	defer resp.Body.Close()
	if res, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	return
}

//定时任务
func TimeSend(method string) {
	c := cron.New()
	spec := "@every 5m"
	c.AddFunc(spec, func() {
		//fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		resp, err := HttpRequest(method)
		if err != nil {
			return
		}
		var respData = RespData{}
		err = json.Unmarshal(resp, &respData)
		timeNow := time.Now().Format("2006-01-02 15:04:05")
		if err != nil {
			logger.Printf("Now is %v,This is Error: Cookie is invalid\n", timeNow)
			return
		}
	})
	c.Start() //开启定时任务
	select {}
	//defer cron.Stop()  //停止定时任务
}

func main() {
	InitLog()
	TimeSend("POST")
}
