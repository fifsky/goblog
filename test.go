package main

import (
	"fmt"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/fifsky/goblog/system"
	"github.com/fifsky/goblog/helpers/tuling"
)

func main() {

	APP_ENV := os.Getenv("APP_ENV")
	if APP_ENV == "local" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	fmt.Println("Run Mode:", gin.Mode())

	system.LoadConfig()

	content, err := tuling.Say("测试测试", "")
	fmt.Println(content,err)
}
