package server

import "github.com/gin-gonic/gin"

var serv = gin.Default()

func Serv() *gin.Engine {
	return serv
}
