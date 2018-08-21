package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"publish/models"
	"publish/session"
	"publish/tools"
	"publish/websocket"
	"strconv"
	"strings"
	"time"

	"github.com/kataras/iris"
)

func (c *Controllers) TaskIndex(ctx iris.Context) {
	t := new(models.Task)
	list := t.List()
	blob, _ := json.Marshal(list)
	ctx.View("task/index.html")
	a := "<script> var a = " + string(blob) + "</script>"
	ctx.Write([]byte(a))
}

// 删除提交
func (c *Controllers) TaskDel(ctx iris.Context) {
	taskID, err := ctx.URLParamInt("taskid")
	if err != nil {
		ctx.WriteString("need taskid")
		return
	}
	t := models.Task{Id: taskID}
	t.Del()
	ctx.Redirect("/task/index", 302)
}

// 查询版本记录
func (c *Controllers) TaskCommitList(ctx iris.Context) {
	params := ctx.FormValues()

	commit := params["commit"][0]
	id := params["task"][0]
	commit = commit[:7]

	t1 := new(models.Project)
	project, err := t1.Find(id)
	if err != nil {
		ctx.WriteString(fmt.Sprintf("%s", err))
		return
	}

	command := new(tools.Command)
	output, err := command.LocalCommandOutput("cd " + project.DeployFrom + " && " + "git log -1 " + commit + " --name-only")
	if err != nil {
		log.Println(err)
		return
	}
	ctx.WriteString(string(output))
}

// task分页记录
func (c *Controllers) TaskPage(ctx iris.Context) {
	pageNo, err := ctx.Params().GetInt("pageNo")
	if err != nil {
		ctx.WriteString("need pageNo")
		return
	}
	records := make([]models.Task, 0)
	t := new(models.Task)
	err = t.Page(pageNo, &records)
	if err != nil {
		log.Println(err)
		return
	}

	blob, _ := json.Marshal(records)
	ctx.Write(blob)
}

// 新任务添加
func (c *Controllers) TaskNewCommit(ctx iris.Context) {
	task := models.Task{}
	err := ctx.ReadForm(&task)
	if err != nil {
		log.Println(err)
		return
	}
	task.CommitId = task.CommitId[:7]

	s := session.Sess.Start(ctx)
	userID, err := s.GetInt("user_id")
	if err != nil {
		ctx.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}
	task.UserId = userID
	task.EnableRollback = 1
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	t := new(models.Task)
	lastRecord := t.FindLast(task.ProjectId)
	task.LinkId = time.Now().Format("20060102-150405")
	task.ExLinkId = lastRecord.LinkId
	_, err = t.New(task)
	if err != nil {
		log.Println(err)
		return
	}
	ctx.Redirect("/task/index", 302)
}

// 新任务页面
func (c *Controllers) TaskNew(ctx iris.Context) {
	id := ctx.URLParam("id")

	t := new(models.Task)
	task, err := t.FindUndo(id)
	if err != nil {
		ctx.WriteString(fmt.Sprintf("%s", err))
	}

	if task.Id > 0 {
		ctx.WriteString(fmt.Sprintf("存在未提交的项目，先提交或删除%d", task.Id))
		return
	}

	t1 := new(models.Project)
	project, err := t1.Find(id)
	if err != nil {
		ctx.WriteString(fmt.Sprintf("%s", err))
		return
	}

	s := session.Sess.Start(ctx)
	userRole, err := s.GetInt("user_role")
	if err != nil {
		ctx.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}
	if userRole != 2 && project.Level == 3 {
		ctx.WriteString("只有管理员能创建线上环境")
		return
	}

	command := new(tools.Command)
	port := strings.Split(project.Hosts, ":")
	command.Host = port[0]
	command.Port, _ = strconv.Atoi(port[1])

	_, err = command.LocalCommandOutput("ls " + project.DeployFrom)
	if err != nil {
		ctx.WriteString(fmt.Sprintf("不存在本地项目目录，需要初始化"))
		log.Println(err)
		return
	}

	_, err = command.LocalCommandOutput("cd " + project.DeployFrom + " && " + "git pull")
	if err != nil {
		log.Println(err)
		return
	}

	branchList, err := command.LocalCommandOutput("cd " + project.DeployFrom + " && " + "git branch -a")
	if err != nil {
		log.Println(err)
		return
	}
	// fmt.Printf("%q", branchList)
	branchListArr := strings.Split((strings.TrimSpace(string(branchList))), "\n")

	logList, err := command.LocalCommandOutput("cd " + project.DeployFrom + " && " + "git log -20 --pretty=\"%h - %an - %s - %cD\"")
	if err != nil {
		log.Println(err)
		return
	}

	_, err = command.LocalCommandOutput("cd " + project.DeployFrom + " && " + "git pull")
	if err != nil {
		log.Println(err)
		return
	}
	// fmt.Printf("%q", logList)
	logListtArr := strings.Split((strings.TrimSpace(string(logList))), "\n")

	ctx.ViewData("branchList", branchListArr)
	ctx.ViewData("logList", logListtArr)
	ctx.ViewData("project", project)
	ctx.ViewData("taskId", id)
	ctx.View("task/new.html")
}

func (c *Controllers) Task(ctx iris.Context) {
	t := new(models.Task)
	list := t.List()
	blob, _ := json.Marshal(list)
	ctx.Write(blob)
}

func (c *Controllers) TaskDetail(ctx iris.Context) {
	t := new(models.Task)
	detail := t.Task(2069)
	blob, _ := json.Marshal(detail)
	ctx.Write(blob)
}

// 发布页面
func (c *Controllers) TaskDeploy(ctx iris.Context) {
	ctx.View("task/deploy.html")
}

// 部署页面
func (c *Controllers) TaskSubmmit(ctx iris.Context) {
	id, err := ctx.URLParamInt("taskid")
	if err != nil {
		log.Println(err)
		return
	}
	t := new(models.Task)
	task := t.Task(id)

	t1 := new(models.Project)
	project, err := t1.Find(task.ProjectId)
	if err != nil {
		ctx.WriteString(fmt.Sprintf("%s", err))
		return
	}

	ctx.ViewData("taskID", id)
	ctx.ViewData("task", task)
	ctx.ViewData("project", project)
	ctx.View("task/submmit.html")
}

// 部署动作
func (c *Controllers) TaskShift(ctx iris.Context) {
	id, err := ctx.URLParamInt("id")

	if err != nil {
		log.Println(err)
		return
	}
	t := new(models.Task)
	task := t.Find(id)

	p := new(models.Project)
	project, err := p.Find(task.ProjectId)
	if err != nil {
		ctx.WriteString(fmt.Sprintf("%s", err))
		return
	}

	switch task.FileTransmissionMode {
	case 1:
		err = fullDeploy(project, task)
	case 2:
		err = listDeploy(project, task)
	}
	if err != nil {
		websocket.Broadcast(websocket.Conn, fmt.Sprintf("部署发生了错误%s", err))
		log.Println(err)
		return
	}

	t.SetStatus(task.Id, 3)

	result := models.NewDefaultReturn()
	result.Sta = 1
	blob, _ := json.Marshal(result)
	ctx.Write(blob)
}

// 完整部署
func fullDeploy(project *models.Project, task *models.Task) (err error) {
	websocket.Broadcast(websocket.Conn, fmt.Sprintf("run remote command:%s\r\n", "开始完整部署"))
	// 文件打包
	destFile, err := tools.Compress(project.DeployFrom, "repos/"+task.LinkId)
	if err != nil {
		websocket.Broadcast(websocket.Conn, fmt.Sprintf("upload: %s\r\n", err.Error()))
		log.Println(err.Error())
		return
	}
	// 上传到服务器并链接
	command := new(tools.Command)

	port := strings.Split(project.Hosts, ":")
	command.Host = port[0]
	command.Port, _ = strconv.Atoi(port[1])
	// 上传文件
	err = command.FileUpload(destFile, project.ReleaseLibrary+path.Base(project.DeployFrom)+"/"+task.LinkId+".tar.gz")
	if err != nil {
		log.Println(err)
		return
	}

	// 远程解压文件
	err = command.RemoteCommand("tar -xvf " + project.ReleaseLibrary + path.Base(project.DeployFrom) + "/" + task.LinkId + ".tar.gz -C " + project.ReleaseLibrary + path.Base(project.DeployFrom))
	if err != nil {
		log.Println(err)
		return
	}
	// 链接
	err = command.RemoteCommand("ln -sfn " + project.ReleaseLibrary + path.Base(project.DeployFrom) + "/" + task.LinkId + " " + project.ReleaseTo + path.Base(project.DeployFrom))
	if err != nil {
		log.Println(err)
		return
	}
	// 删除gz包
	err = command.RemoteCommand("rm -rf " + project.ReleaseLibrary + path.Base(project.DeployFrom) + "/*.tar.gz")
	if err != nil {
		log.Println(err)
		return
	}
	err = os.Remove(destFile)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// 部署后命令
	cmds := strings.Split(strings.TrimSpace(project.PostRelease), "\r\n")
	for _, cmd := range cmds {
		err := command.RemoteCommand(cmd)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return
}

// 列表部署
func listDeploy(project *models.Project, task *models.Task) (err error) {
	websocket.Broadcast(websocket.Conn, fmt.Sprintf("run remote command:%s\r\n", "开始列表部署"))
	// 文件打包
	files := strings.Split(strings.TrimSpace(task.FileList), "\r\n")
	for k, v := range files {
		files[k] = project.DeployFrom + "/" + v
	}
	var destFile string
	destFile, err = tools.CompressFiles(files, project, "repos/"+task.LinkId)
	if err != nil {
		log.Println(err)
		return
	}

	command := new(tools.Command)

	// 更新本地代码到最新版本
	err = command.LocalCommand("git -C " + project.DeployFrom + " pull")
	if err != nil {
		log.Println(err)
		return
	}

	port := strings.Split(project.Hosts, ":")
	command.Host = port[0]
	command.Port, _ = strconv.Atoi(port[1])

	// 上传服务器并链接
	command.FileUpload(destFile, project.ReleaseLibrary+path.Base(project.DeployFrom)+"/"+task.LinkId+".tar.gz")

	// 备份当前版本
	err = command.RemoteCommand("cp -arf " + project.ReleaseTo + path.Base(project.DeployFrom) + "/. " + project.ReleaseLibrary + path.Base(project.DeployFrom) + "/" + task.LinkId)
	if err != nil {
		log.Println(err)
		return
	}

	// 合并文件
	err = command.RemoteCommand("tar -xvf " + project.ReleaseLibrary + path.Base(project.DeployFrom) + "/" + task.LinkId + ".tar.gz -C " + project.ReleaseLibrary + path.Base(project.DeployFrom) + "/" + task.LinkId)
	if err != nil {
		log.Println(err)
		return
	}

	// 链接
	err = command.RemoteCommand("ln -sfn " + project.ReleaseLibrary + path.Base(project.DeployFrom) + "/" + task.LinkId + " " + project.ReleaseTo + path.Base(project.DeployFrom))
	if err != nil {
		log.Println(err)
		return
	}

	// 删除gz包
	err = command.RemoteCommand("rm -rf " + project.ReleaseLibrary + path.Base(project.DeployFrom) + "/*.tar.gz")
	if err != nil {
		log.Println(err)
		return
	}

	// 删除文件包
	err = os.Remove(destFile)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// 部署后命令
	cmds := strings.Split(strings.TrimSpace(project.PostRelease), "\r\n")
	for _, cmd := range cmds {
		err := command.RemoteCommand(cmd)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return
}
