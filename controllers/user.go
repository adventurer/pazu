package controllers

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"publish/cache"
	"publish/models"
	"publish/session"
	"time"

	"github.com/kataras/iris"
	uuid "github.com/satori/go.uuid"
)

func (c *Controllers) Users(ctx iris.Context) {
	blob, _ := json.Marshal(cache.MemUsers)
	ctx.Write(blob)
}

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
		ctx.Redirect("/task/index")
	} else {
		ctx.WriteString(fmt.Sprintf("%s", "账号或密码错误"))
	}
}

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

func (c *Controllers) UserLogout(ctx iris.Context) {
	ctx.SetCookieKV("auth_key", "")
	ctx.Redirect("/user/login", 302)
}
