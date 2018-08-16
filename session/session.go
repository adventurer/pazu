package session

import (
	"time"

	"github.com/kataras/iris/sessions"
)

var Sess = new(sessions.Sessions)

func init() {
	Sess = sessions.New(sessions.Config{
		// Cookie string, the session's client cookie name, for example: "mysessionid"
		//
		// Defaults to "irissessionid"
		Cookie: "mysessionid",
		// it's time.Duration, from the time cookie is created, how long it can be alive?
		// 0 means no expire.
		// -1 means expire when browser closes
		// or set a value, like 2 hours:
		Expires: time.Hour * 24,
		// if you want to invalid cookies on different subdomains
		// of the same host, then enable it.
		// Defaults to false.
		DisableSubdomainPersistence: true,
		// AllowReclaim will allow to
		// Destroy and Start a session in the same request handler.
		// All it does is that it removes the cookie for both `Request` and `ResponseWriter` while `Destroy`
		// or add a new cookie to `Request` while `Start`.
		//
		// Defaults to false.
		AllowReclaim: true,
	})
}
