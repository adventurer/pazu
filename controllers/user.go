package controllers

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"publish/cache"
	"publish/models"
	"publish/session"
	"time"

	"github.com/kataras/iris"
	uuid "github.com/satori/go.uuid"
)

// 列表
func (c *Controllers) Users(ctx iris.Context) {
	blob, _ := json.Marshal(cache.MemUsers)
	ctx.Write(blob)
}

// 登入
func (c *Controllers) UserLoginSubmit(ctx iris.Context) {
	user := models.User{}
	if err := ctx.ReadForm(&user); err != nil {
		ctx.WriteString(fmt.Sprintf("%s", err))
		return
	}
	userNew, err := user.FindByUsername()
	if err != nil {
		ctx.WriteString(fmt.Sprintf("%s", err))
	}
	oldpass := fmt.Sprintf("%x", md5.Sum([]byte(user.PasswordHash)))
	newpass := userNew.PasswordHash

	access := base64.URLEncoding.EncodeToString(uuid.NewV3(uuid.Must(uuid.NewV4()), oldpass).Bytes())

	if oldpass == newpass {
		userNew.AuthKey = access
		userNew.UpdatedAt = time.Now()
		userNew.SetAccessTocken()

		s := session.Sess.Start(ctx)
		s.Set("user_id", userNew.Id)
		s.Set("user_role", userNew.Role)

		// 登陆后必须刷新缓存
		cache.CacheUserHasTable()
		cache.CacheUsers()
		ctx.Redirect("/task/index")
	} else {
		ctx.WriteString(fmt.Sprintf("%s", "账号或密码错误"))
	}
}

// 注册提交
func (c *Controllers) UserRegisterSubmit(ctx iris.Context) {
	user := models.User{}
	if err := ctx.ReadForm(&user); err != nil {
		ctx.WriteString(fmt.Sprintf("%s", err))
		return
	}
	user.Status = 0
	user.CreatedAt = time.Now()
	has := md5.Sum([]byte(user.PasswordHash))
	user.PasswordHash = fmt.Sprintf("%x", has)
	if _, err := user.New(); err != nil {
		ctx.WriteString(fmt.Sprintf("%s", err))
		return
	}
	ctx.Redirect("/user/login")
}

// 登出
func (c *Controllers) UserLogout(ctx iris.Context) {
	ctx.SetCookieKV("auth_key", "")
	ctx.Redirect("/user/login", 302)
}

// 用户界面
func (c *Controllers) UserCtl(ctx iris.Context) {
	ctx.View("user/index.html")
}

// 用户列表
func (c *Controllers) UserList(ctx iris.Context) {
	u := models.User{}
	users, err := u.GetFullList()
	if err != nil {
		log.Println(err)
		return
	}
	blob, _ := json.Marshal(users)
	ctx.Write(blob)
}

// 激活
func (c *Controllers) UserActive(ctx iris.Context) {
	role, err := ctx.URLParamInt("role")
	if err != nil {
		ctx.WriteString(fmt.Sprintf("%s", err))
		return
	}
	id, err := ctx.URLParamInt("id")
	if err != nil {
		ctx.WriteString(fmt.Sprintf("%s", err))
		return
	}
	u := models.User{}
	u.SetRoleById(id, role)
	ctx.Redirect("/user/ctl")
}

// 删除
func (c *Controllers) UserDel(ctx iris.Context) {
	id, err := ctx.URLParamInt("id")
	if err != nil {
		ctx.WriteString(fmt.Sprintf("%s", err))
		return
	}
	u := models.User{}
	u.DelUserById(id)
	ctx.Redirect("/user/ctl")

}
