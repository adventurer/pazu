package middleware

import (
	"log"
	"publish/session"

	"github.com/kataras/iris"
)

func Authentication(ctx iris.Context) {
	s := session.Sess.Start(ctx)
	userId, err := s.GetInt("user_id")
	if err != nil {
		log.Println(err)
		ctx.Redirect("/user/login", 302)
		return
	}
	if userId <= 0 {
		ctx.WriteString("invalid userid < 0")
		return
	}
	ctx.Next()
	return
}
