package tools

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var C http.Client

type HttpTest struct {
	Url      string
	Response *http.Response
}

func init() {
	C = http.Client{
		Transport: &http.Transport{
			// 忽略证书错误
			// TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(25 * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*20)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}
}

func NewHttp(url string) *HttpTest {
	return &HttpTest{Url: url, Response: nil}
}

func (h *HttpTest) HttpGet() (err error) {

	h.Response, err = C.Get(h.Url)

	if err != nil {
		return
	}
	defer h.Response.Body.Close()
	return
}

func (h *HttpTest) HttpPost() (err error) {

	req, err := http.NewRequest("POST", h.Url, strings.NewReader(""))

	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Header.Set("User-Agent", "Paw/3.1.7 (Macintosh; OS X/10.13.6) GCDHTTPRequest")
	req.Header.Set("Connection", "close")

	h.Response, err = C.Do(req)

	if err != nil {
		return
	}
	defer h.Response.Body.Close()
	return

}

func httpPostForm() {

	resp, err := http.PostForm("http://www.01happy.com/demo/accept.php",

		url.Values{"key": {"Value"}, "id": {"123"}})

	if err != nil {
		// handle error
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}

func httpDo() {

	client := &http.Client{}

	req, err := http.NewRequest("POST", "http://www.01happy.com/demo/accept.php", strings.NewReader("name=cjb"))

	if err != nil {
		// handle error
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	req.Header.Set("Cookie", "name=anny")

	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		// handle error
	}
	fmt.Println(string(body))
}
