package routes

import (
	"publish/controllers"
	"publish/middleware"

	"github.com/kataras/iris"
)

type Routes struct{}

var (
	controller = new(controllers.Controllers)
)

func (r *Routes) InitRoute(app *iris.Application) {

	app.Any("/user/login", func(ctx iris.Context) {
		ctx.ViewLayout(iris.NoLayout)
		ctx.View("user/login.html")
	})

	app.Any("/user/register", func(ctx iris.Context) {
		ctx.ViewLayout(iris.NoLayout)
		ctx.View("user/register.html")
	})

	app.Any("/user/register-submit", controller.UserRegisterSubmit)
	app.Any("/user/login-submit", controller.UserLoginSubmit)

	adminRoutes := app.Party("/", middleware.Authentication, middleware.CheckAdmin)
	{

		adminRoutes.Get("/project/index", controller.ProjectIndex)
		adminRoutes.Get("/project/init", controller.ProjectInitialize)
		adminRoutes.Get("/project/new", controller.ProjectNew)
		adminRoutes.Post("/project/commit", controller.ProjectCommit)
		adminRoutes.Get("/project/edit", controller.ProjectEdit)
		adminRoutes.Post("/project/edit-commit", controller.ProjectEditCommit)
		adminRoutes.Get("/project/del", controller.ProjectDel)
		adminRoutes.Get("/project/copy", controller.ProjectCopy)
	}

	usersRoutes := app.Party("/", middleware.Authentication)
	{
		// 用户管理

		usersRoutes.Get("/users", controller.Users)
		usersRoutes.Get("/users/list", controller.Conf)
		usersRoutes.Get("/users/add", controller.Conf)
		usersRoutes.Get("/user/logout", controller.UserLogout)

		// 首页
		usersRoutes.Get("/", controller.ProjectIndex)

		// 项目配置
		usersRoutes.Post("/conf", controller.Conf)
		usersRoutes.Get("/projects", controller.Projects)

		// 我的上线单
		usersRoutes.Any("/task/index", controller.TaskIndex)
		usersRoutes.Any("/task", controller.Task)
		usersRoutes.Any("/task/detail", controller.TaskDetail)
		usersRoutes.Any("/task/deploy", controller.TaskDeploy)
		usersRoutes.Get("/task/submmit", controller.TaskSubmmit)
		usersRoutes.Get("/task/shift", controller.TaskShift)
		usersRoutes.Get("/task/new", controller.TaskNew)
		usersRoutes.Any("/task/commit", controller.TaskCommitList)
		usersRoutes.Any("/task/tasknewcommit", controller.TaskNewCommit)
		usersRoutes.Any("/task/page/{pageNo:int}", controller.TaskPage)

		// 文件校验
		usersRoutes.Get("/check", controller.Conf)

		// 运行日志
		usersRoutes.Any("/runlog/index", controller.RunlogIndex)

	}

}
