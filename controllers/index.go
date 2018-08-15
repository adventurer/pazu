package controllers

import "github.com/kataras/iris"

type Controllers struct{}

func (c *Controllers) Index(ctx iris.Context) {
	ctx.Writef("i am the king of king")
}
