package controllers

import (
	"encoding/json"
	"publish/models"

	"github.com/kataras/iris"
)

func (c *Controllers) Conf(ctx iris.Context) {
	pro := new(models.Project)
	list := pro.List()
	blob, _ := json.Marshal(list)
	ctx.Write(blob)
}
