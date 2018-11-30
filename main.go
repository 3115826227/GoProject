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
	"fmt"
)

type RespData struct {
	Status    int         `json:"status"`
	Data      interface{} `json:"data"`
	ErrorCode interface{} `json:"errorCode"`
}

const (
	url         = "https://tuiguang.baidu.com/request.ajax?path=GET/authInfo"
	cookie      = "BAIDUID=78B8B9FA589A387F0291E281FFB52771:FG=1; BIDUPSID=78B8B9FA589A387F0291E281FFB52771; PSTM=1508221353; MCITY=-%3A; BDORZ=B490B5EBF6F3CD402E515D22BCDA1598; SAMPLING_USER_ID=25847076; ISFROMV3=25847076; NSTI=25847076; H_PS_PSSID=1449_21104_18559_27889; delPer=0; PSINO=2; OPENAB=4; uc_login_unique=0f89443520d47a40b31a25d7742bfc2f; SIGNIN_UC=70a2711cf1d3d9b1a82d2f87d633bd8a02935843088; uc_recom_mark=cmVjb21tYXJrXzI1ODQ3MDc2; __cas__st__3=3515250e9d36e0a88ad96395542669385c4f7390b72b51592b4f88d9b7f6825b3ba00bc18a2fda41ab2fc16a; __cas__id__3=25847076; __cas__rn__=293584308; Hm_lvt_ab31944d33b258d42a263a7c78b303a3=1543546636,1543556117,1543556352,1543571863; Hm_lpvt_ab31944d33b258d42a263a7c78b303a3=1543571863; __bsi=11242949868071590397_00_39_N_R_15_0303_c02f_Y"
	contentType = "application/x-www-form-urlencoded"
	userAgent   = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 Safari/537.36"
	RespBody    = "userid=25847076&token=3515250e9d36e0a88ad96395542669385c4f7390b72b51592b4f88d9b7f6825b3ba00bc18a2fda41ab2fc16a"
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

func getTimeNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

//定时任务
func TimeSend(method string) {
	logger.Printf("Start Now. The time is %v", getTimeNow())
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
		if err != nil {
			logger.Printf("Now is %v,This is Error: Cookie is invalid\n", getTimeNow())
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
