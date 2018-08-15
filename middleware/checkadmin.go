package middleware

import (
	"fmt"
	"publish/session"

	"github.com/kataras/iris"
)

func CheckAdmin(ctx iris.Context) {
	s := session.Sess.Start(ctx)
	userRole, err := s.GetInt("user_role")
	if err != nil {
		ctx.Write([]byte(fmt.Sprintf("%s", err)))
		return
	}
	if userRole != 2 {
		ctx.WriteString("just admin access")
		return
	}
	ctx.Next()
	return
}
