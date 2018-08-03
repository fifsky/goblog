package main

import (
	"os"
	"fmt"
	"flag"
	"io/ioutil"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gin-gonic/gin"
	"github.com/fifsky/goblog/core"
	"github.com/fifsky/goblog/config"
	"github.com/fifsky/goblog/route"
	"github.com/ilibs/gosql"
	"github.com/ilibs/logger"
)

func main() {
	gosql.Connect(config.App.DB)
	fmt.Println("Run Mode:", gin.Mode())
	flag.Parse()
	cmd := flag.Arg(0)
	if cmd == "install" {
		_, err := config.ImportDB()
		if err != nil {
			fmt.Println("Import DB Error:" + err.Error())
			logger.Error(err)
		}
		return
	}

	router := gin.Default()
	//路由
	route.Route(router)
	setPid(os.Getpid())
	//定时提醒
	go core.StartCron()

	router.Run(":8080")
}

func setPid(pid int) {
	d := []byte(strconv.Itoa(pid))
	err := ioutil.WriteFile("./blog.pid", d, 0644)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
}
