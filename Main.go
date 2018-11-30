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
	url         = "https://tuiguang.baidu.com/request.ajax?path=GET/authInfo"
	cookie      = "BAIDUID=78B8B9FA589A387F0291E281FFB52771:FG=1; BIDUPSID=78B8B9FA589A387F0291E281FFB52771; PSTM=1508221353; MCITY=-%3A; BDORZ=B490B5EBF6F3CD402E515D22BCDA1598; SAMPLING_USER_ID=25847076; delPer=0; PSINO=2; pgv_pvi=2268574720; pgv_si=s6550958080; uc_login_unique=2937b4ebdc89dad28776b9653723fc19; uc_recom_mark=cmVjb21tYXJrXzI1ODQ3MDc2; H_PS_PSSID=1449_21104_18559_27889; SIGNIN_UC=70a2711cf1d3d9b1a82d2f87d633bd8a02935543588; __cas__st__3=59307ec276c0ec398697602f3b41762a42faba2c72d205d60a967e4a90c872509522e1f483d78e3b8e0b66ee; __cas__id__3=25847076; __cas__rn__=293554358; Hm_lvt_ab31944d33b258d42a263a7c78b303a3=1542358715,1542358738,1543399577,1543541912; Hm_lpvt_ab31944d33b258d42a263a7c78b303a3=1543541920; __bsi=6815283152668088576_00_112_N_R_3_0303_c02f_Y"
	contentType = "application/x-www-form-urlencoded"
	userAgent   = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 Safari/537.36"
	RespBody    = "userid=25847076&token=59307ec276c0ec398697602f3b41762a42faba2c72d205d60a967e4a90c872509522e1f483d78e3b8e0b66ee"
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
