package health

import (
	"encoding/json"
	"fmt"
	"publish/cache"
	"publish/models"
	"publish/tools"
	"publish/websocket"
	"time"
)

type CheckResult struct {
	Health models.Health
	Code   int
	Cost   int64
	Msg    string
}

func Check() {
	for {
		for _, v := range cache.MemHealthTable {
			go send(v)
		}
		time.Sleep(1 * time.Second)
	}
}

func send(v models.Health) {
	var health tools.HttpTest
	begin := time.Now().UnixNano()
	health.Url = v.Url
	err := health.HttpGet()
	end := time.Now().UnixNano()

	var responseJson CheckResult
	if err != nil {
		responseJson = CheckResult{Health: v, Code: -1, Cost: (end - begin) / 1000000, Msg: fmt.Sprintf("%s", err)}
	} else {
		responseJson = CheckResult{Health: v, Code: health.Response.StatusCode, Cost: (end - begin) / 1000000, Msg: "正常"}
	}

	blob, _ := json.Marshal(responseJson)
	websocket.BroadcastHeath(websocket.Conn, fmt.Sprintf(string(blob)))
}
