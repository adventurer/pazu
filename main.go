package main

import (
	"publish/routes"

	"publish/tools"

	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
)

var (
	cookieNameForSessionID = "mycookiesessionnameid"
	sess                   = sessions.New(sessions.Config{Cookie: cookieNameForSessionID, AllowReclaim: true})
)

func main() {

	app := iris.New()

	app.RegisterView(iris.HTML("./views", ".html").Layout("layouts/layout.html").Reload(true))

	route := new(routes.Routes)
	route.InitRoute(app)

	app.StaticWeb("./assets", "./assets")

	tools.SetupWebsocket(app)

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
