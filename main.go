package main

import (
	"publish/routes"
	"publish/websocket"

	"github.com/kataras/iris"
)

func main() {

	app := iris.New()

	app.RegisterView(iris.HTML("./views", ".html").Layout("layouts/layout.html").Reload(true))

	route := new(routes.Routes)
	route.InitRoute(app)

	app.StaticWeb("./assets", "./assets")

	websocket.SetupWebsocket(app)

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
