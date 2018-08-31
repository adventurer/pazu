package health

import (
	"encoding/json"
	"fmt"
	"log"
	"publish/cache"
	"publish/mail"
	"publish/models"
	"publish/tools"
	"publish/websocket"
	"strings"
	"time"
)

var waring = make([]models.Health, len(cache.MemHealthTable))

func init() {
	copy(waring, cache.MemHealthTable)
}

type CheckResult struct {
	Health models.Health
	Code   int
	Cost   int64
	Msg    string
}

func waringCheck(id int, err string) {
	for k, v := range waring {
		if v.Id == id {
			waring[k].Report = waring[k].Report - 1
			log.Println("found waring:", v.Name, v.Report)
		}
		if v.Report <= 0 {
			log.Println("waring send a mail:", v.Name)
			mail.SendEmail(mail.NewEmail("16620808100@163.com", "健康检查故障", v.Name+"发生故障:"+err, "text"))
			waring[k].Report = cache.MemHealthTable[k].Report
		}
	}
}

func waringClear(id int) {
	for k, v := range waring {
		if v.Id == id {
			waring[k].Report = cache.MemHealthTable[k].Report
		}
	}
}

func Check() {
	for {
		for _, v := range cache.MemHealthTable {
			if int(time.Now().Unix())%v.Interval == 0 {
				go send(v)
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func send(v models.Health) {
	var health tools.HttpTest
	var err error
	health.Url = v.Url
	begin := time.Now().UnixNano()
	if strings.TrimSpace(v.Method) == "post" {
		err = health.HttpPost()
	} else {
		err = health.HttpGet()
	}
	end := time.Now().UnixNano()

	var responseJson CheckResult

	if err != nil {
		waringCheck(v.Id, fmt.Sprintf("err:%s,status:%s", err, "cant access"))
		responseJson = CheckResult{Health: v, Code: -1, Cost: (end - begin) / 1000000, Msg: fmt.Sprintf("%s", err)}
	} else {
		if health.Response.StatusCode != 200 {
			waringCheck(v.Id, fmt.Sprintf("err:%s,status:%s", "", health.Response.Status))
		}
		waringClear(v.Id)
		responseJson = CheckResult{Health: v, Code: health.Response.StatusCode, Cost: (end - begin) / 1000000, Msg: fmt.Sprintf("%s", health.Response.Status)}
	}

	blob, _ := json.Marshal(responseJson)
	websocket.BroadcastHeath(websocket.Conn, fmt.Sprintf(string(blob)))
}
