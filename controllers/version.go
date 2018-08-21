package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"path"
	"publish/models"
	"publish/tools"
	"strconv"
	"strings"

	"github.com/kataras/iris"
)

// 版本管理页面
func (c *Controllers) VersionCtl(ctx iris.Context) {
	id, err := ctx.URLParamInt("id")
	if err != nil {
		ctx.WriteString(fmt.Sprintf("$s", err))
		return
	}
	ctx.ViewData("id", id)
	ctx.View("versions/index.html")
}

// 版本列表
func (c *Controllers) VersionList(ctx iris.Context) {
	projectId, err := ctx.URLParamInt("id")
	if err != nil {
		ctx.WriteString(fmt.Sprintf("%s", err))
		return
	}
	p := models.Project{Id: projectId}
	project, err := p.Find(p.Id)
	if err != nil {
		ctx.WriteString(fmt.Sprintf("%s", err))
		return
	}

	port := strings.Split(project.Hosts, ":")

	command := tools.Command{}
	command.Host = port[0]
	command.Port, _ = strconv.Atoi(port[1])
	var cmd string
	var output string
	cmd = "ls " + project.ReleaseLibrary + path.Base(project.DeployFrom)
	output, err = command.RemoteCommandOutput(cmd)
	if err != nil {
		ctx.WriteString(fmt.Sprintf("%s", err))
		return
	}

	versions := strings.Split(strings.TrimSpace(output), "\n")
	log.Println(len(versions))

	versionsModel := models.Version{}
	versionSlice := versionsModel.VersionList(versions)

	cmd = "ls " + project.ReleaseTo + path.Base(project.DeployFrom) + " -al"
	output, err = command.RemoteCommandOutput(cmd)
	if err != nil {
		ctx.WriteString(fmt.Sprintf("%s", err))
		return
	}
	for k, v := range versionSlice {
		if strings.Contains(output, v.Id) {
			versionSlice[k].Active = 1
		}
	}

	blob, _ := json.Marshal(versionSlice)
	ctx.Write(blob)

}

// 版本切换
func (c *Controllers) VersionSwitch(ctx iris.Context) {
	id := ctx.URLParam("id")
	projectID, err := ctx.URLParamInt("project")
	if err != nil {
		ctx.WriteString(fmt.Sprintf("%s", err))
		return
	}
	p := models.Project{}
	project, err := p.Find(projectID)
	if err != nil {
		ctx.WriteString(fmt.Sprintf("%s", err))
		return
	}

	command := tools.Command{}
	port := strings.Split(project.Hosts, ":")
	command.Host = port[0]
	command.Port, _ = strconv.Atoi(port[1])
	cmd := "ln -sfn " + project.ReleaseLibrary + path.Base(project.DeployFrom) + "/" + id + " " + project.ReleaseTo + path.Base(project.DeployFrom)
	_, err = command.RemoteCommandOutput(cmd)
	if err != nil {
		ctx.WriteString(fmt.Sprintf("%s", err))
		return
	}

	// 部署后命令
	cmds := strings.Split(strings.TrimSpace(project.PostRelease), "\r\n")
	for _, cmd := range cmds {
		err := command.RemoteCommand(cmd)
		if err != nil {
			ctx.WriteString(fmt.Sprintf("%s", err))
			return
		}
	}

	ctx.Redirect("/version/ctl?id=" + strconv.Itoa(projectID))
}

// 版本删除
func (c *Controllers) VersionDestory(ctx iris.Context) {

}
