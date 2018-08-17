package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"path"
	"publish/cache"
	"publish/models"
	"publish/tools"
	"strconv"
	"strings"
	"time"

	"github.com/kataras/iris"
)

func (c *Controllers) ProjectList(ctx iris.Context) {
	pro := new(models.Project)
	list := pro.List()
	blob, _ := json.Marshal(list)
	ctx.Write(blob)
}

func (c *Controllers) Projects(ctx iris.Context) {
	blob, _ := json.Marshal(cache.MemProject)
	ctx.Write(blob)
}

func (c *Controllers) ProjectDel(ctx iris.Context) {
	id, err := ctx.URLParamInt("id")
	if err != nil {
		ctx.WriteString(fmt.Sprintf("%s", err))
		return
	}
	p := models.Project{Id: id}
	p.Del()
	ctx.Redirect("/project/index", 302)
}

func (c *Controllers) ProjectCopy(ctx iris.Context) {
	id, err := ctx.URLParamInt("id")
	if err != nil {
		ctx.WriteString(fmt.Sprintf("%s", err))
		return
	}
	p := models.Project{}
	project, err := p.Find(id)
	if err != nil {
		ctx.WriteString(fmt.Sprintf("%s", err))
		return
	}
	project.Id = models.Project{}.Id
	_, err = project.New()
	if err != nil {
		ctx.WriteString(fmt.Sprintf("%s", err))
		return
	}
	ctx.Redirect("/project/index", 302)
}

// 新增页面
func (c *Controllers) ProjectNew(ctx iris.Context) {
	ctx.View("project/new.html")
}

// 编辑页面
func (c *Controllers) ProjectEdit(ctx iris.Context) {
	id := ctx.URLParam("id")
	p := models.Project{}
	project, err := p.Find(id)
	if err != nil {
		ctx.WriteString(fmt.Sprintf("%s", err))
		return
	}
	ctx.ViewData("project", project)
	ctx.View("project/edit.html")
}

// 编辑提交
func (c *Controllers) ProjectEditCommit(ctx iris.Context) {
	p := models.Project{}
	err := ctx.ReadForm(&p)
	p.UpdatedAt = time.Now()
	if err != nil {
		log.Println(err)
		return
	}
	_, err = p.Edit()
	if err != nil {
		log.Println(err)
		return
	}
	ctx.Redirect("/project/index", 302)
}

// 新增提交
func (c *Controllers) ProjectCommit(ctx iris.Context) {
	p := models.Project{}
	err := ctx.ReadForm(&p)
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	if err != nil {
		log.Println(err)
		return
	}
	_, err = p.New()
	if err != nil {
		log.Println(err)
		return
	}
	ctx.Redirect("/project/index", 302)
}

func (c *Controllers) ProjectInitialize(ctx iris.Context) {
	// 得到当前项目配置
	id := ctx.URLParam("id")
	proObj := new(models.Project)
	project, err := proObj.Find(id)
	if err != nil {
		ctx.WriteString(fmt.Sprintf("%s", err))
		return
	}
	if project.Id <= 0 {
		ctx.WriteString("未找到项目" + id)
		return
	}
	// 检查本地目录
	log.Println("检查本地目录")
	command := new(tools.Command)
	_, err = command.PathGen(project.DeployFrom)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// 检查仓库配置
	log.Println("检查仓库配置")
	_, err = command.PathExists(project.DeployFrom + "/.git")
	if err != nil {
		log.Println("未找到仓库：" + project.DeployFrom)
		log.Println("初始化仓库：" + project.DeployFrom)
		cmd := "git clone -q " + project.RepoUrl + " " + project.DeployFrom
		err = command.LocalCommand(cmd)
		if err != nil {
			log.Println(err.Error())
			return
		}
	} else {
		log.Println("存在仓库git目录:", project.DeployFrom+"/.git")
	}

	// 检查远程备份目录
	log.Println("检查远程仓库目录")
	proPos := strings.LastIndex(project.DeployFrom, "/")
	proName := project.DeployFrom[proPos+1:]
	// command.Host = "192.168.3.208"
	// command.Port = 22

	port := strings.Split(project.Hosts, ":")
	command.Host = port[0]
	command.Port, _ = strconv.Atoi(port[1])

	x := models.NewDefaultReturn()

	err = command.RemoteCommand("ls " + project.ReleaseLibrary + proName)
	if err != nil {
		log.Println("不存在远程备份目录：" + project.ReleaseLibrary + proName)
		log.Println("开始创建远程备份目录：" + project.ReleaseLibrary + proName)
		command.RemoteCommand("mkdir -p " + project.ReleaseLibrary + proName)
		x.Msg = "创建了远程仓库目录:" + project.ReleaseLibrary + proName

		log.Println(err.Error())
	}
	x.Msg = "存在远程仓库目录:" + project.ReleaseLibrary + proName
	log.Println("存在远程仓库目录：" + project.ReleaseLibrary + proName)

	// 检查远程项目目录
	log.Println("检查远程项目目录")
	err = command.RemoteCommand("ls " + project.ReleaseTo)
	if err != nil {
		log.Println("不存在远程项目目录：" + project.ReleaseTo)

		x.Code = 0
		x.Msg = x.Msg + "失败,不存在远程项目目录：" + project.ReleaseTo
		blob, _ := json.Marshal(x)
		ctx.Write(blob)

	} else {
		log.Println("存在远程项目目录：" + project.ReleaseTo)

		x.Code = 1
		x.Msg = x.Msg + "，存在远程项目目录：" + project.ReleaseTo + ","
		err = command.RemoteCommand("ls " + project.ReleaseTo + "/" + path.Base(project.DeployFrom))
		if err != nil {
			x.Msg = x.Msg + "且不存在项目：" + path.Base(project.DeployFrom)
		} else {
			x.Msg = x.Msg + "但存在项目：" + path.Base(project.DeployFrom)
		}
		blob, _ := json.Marshal(x)
		ctx.Write(blob)
	}

	log.Println("初始化结束")
}

func (c *Controllers) ProjectIndex(ctx iris.Context) {
	blob, _ := json.Marshal(cache.MemProject)
	ctx.View("project/index.html")
	a := "<script> var a = " + string(blob) + "</script>"
	ctx.Write([]byte(a))
}
