package controllers

import "github.com/kataras/iris"

func (c *Controllers) RunlogIndex(ctx iris.Context) {
	ctx.View("runlog/index.html")
}
