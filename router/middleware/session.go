package middleware

import (
	"github.com/fifsky/goblog/config"
	"github.com/gin-gonic/gin"
	"github.com/ilibs/sessions"
	"github.com/ilibs/sessions/cookie"
)

func Sessions() gin.HandlerFunc {
	store := cookie.NewStore([]byte(config.App.Common.SessionSecret))
	store.Options(sessions.Options{HttpOnly: true, MaxAge: 7 * 86400, Path: "/"}) //Also set Secure: true if using SSL, you should though
	return sessions.Sessions("gin-session", store)
}
