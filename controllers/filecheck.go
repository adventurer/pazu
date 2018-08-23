package controllers

import (
	"fmt"
	"log"
	"path"
	"publish/models"
	"publish/tools"
	"strconv"
	"strings"

	"github.com/kataras/iris"
)

type filecheck struct {
	Search int
	Id     int
	File   string
}

func (c *Controllers) FileCheck(ctx iris.Context) {
	if ctx.Method() == "GET" {
		ctx.View("filecheck/index.html")
		return
	}

	form := filecheck{}
	err := ctx.ReadForm(&form)
	if err != nil {
		ctx.WriteString(fmt.Sprintf("%s", err))
		return
	}

	if form.Id == 0 {
		ctx.View("filecheck/index.html")
		return
	}

	p := new(models.Project)
	project, err := p.Find(form.Id)
	if err != nil {
		ctx.WriteString(fmt.Sprintf("%s", err))
		return
	}

	command := new(tools.Command)
	port := strings.Split(project.Hosts, ":")
	command.Host = port[0]
	command.Port, _ = strconv.Atoi(port[1])

	cmd := "cat " + project.ReleaseTo + path.Base(project.DeployFrom) + "/" + form.File
	output, err := command.RemoteCommandOutput(cmd)
	if err != nil {
		ctx.WriteString(fmt.Sprintf("%s", err))
		return
	}
	log.Println(output)

	ctx.ViewData("output", output)
	ctx.View("filecheck/index.html")
}
