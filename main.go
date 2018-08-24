package main

import (
	"publish/health"
	"publish/routes"
	"publish/websocket"

	"github.com/kataras/iris"
)

func main() {
	// 发送邮件给自己
	// mail.SendEmail(mail.NewEmail("16620808100@163.com", "tile", "content", "html"))
	app := iris.New()

	app.RegisterView(iris.HTML("./views", ".html").Layout("layouts/layout.html").Reload(true))

	route := new(routes.Routes)
	route.InitRoute(app)

	app.StaticWeb("./assets", "./assets")
	app.StaticWeb("./html", "./html")

	websocket.SetupWebsocket(app)
	go health.Check()

	app.Run(iris.Addr(":8088"), iris.WithConfiguration(iris.Configuration{
		DisableStartupLog:                 false,
		DisableInterruptHandler:           false,
		DisablePathCorrection:             false,
		EnablePathEscape:                  true,
		FireMethodNotAllowed:              false,
		DisableBodyConsumptionOnUnmarshal: false,
		DisableAutoFireStatusCode:         false,
		TimeFormat:                        "Mon, 02 Jan 2006 15:04:05 GMT",
		Charset:                           "UTF-8",
	}))

}
