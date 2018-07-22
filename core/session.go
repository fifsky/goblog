package core

import (
	"github.com/gin-gonic/gin"
	"github.com/ilibs/sessions"
	"github.com/ilibs/sessions/cookie"
	"github.com/fifsky/goblog/core/config"
)

func SetSessions(router *gin.Engine) {
	store := cookie.NewStore([]byte(config.App.Common.SessionSecret))
	store.Options(sessions.Options{HttpOnly: true, MaxAge: 7 * 86400, Path: "/"}) //Also set Secure: true if using SSL, you should though
	router.Use(sessions.Sessions("gin-session", store))
}