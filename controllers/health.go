package controllers

import (
	"fmt"
	"publish/cache"
	"publish/models"
	"time"

	"github.com/kataras/iris"
)

func (c *Controllers) HealthAppIndex(ctx iris.Context) {
	ctx.View("health/index.html")
}

// 健康检查添加页
func (c *Controllers) HealthAppAdd(ctx iris.Context) {
	ctx.View("health/add.html")
}

// 健康检查提交
func (c *Controllers) HealthAppAddCommit(ctx iris.Context) {
	form := models.Health{}
	err := ctx.ReadForm(&form)
	if err != nil {
		ctx.WriteString(fmt.Sprintf("%s", err))
		return
	}
	form.CreatedAt = time.Now()
	form.UpdatedAt = time.Now()
	form.New()
	cache.CacheHealthTable()
	ctx.Redirect("/health/app")
}

func (c *Controllers) HealthProcess(ctx iris.Context) {
}

// 运行日志检查
func (c *Controllers) HealthLog(ctx iris.Context) {
}
