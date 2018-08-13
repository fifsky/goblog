package main

import (
	"os"
	"fmt"
	"flag"
	"log"
	"io/ioutil"
	"strconv"

	_ "github.com/go-sql-driver/mysql"

	"github.com/fifsky/goblog/core"
	"github.com/fifsky/goblog/config"
	"github.com/ilibs/gosql"
	"github.com/fifsky/goblog/server"
	"github.com/fifsky/goblog/router"
)

func main() {
	gosql.Connect(config.App.DB)
	flag.Parse()
	cmd := flag.Arg(0)
	if cmd == "install" {
		_, err := config.ImportDB()
		if err != nil {
			fmt.Println("Import DB Error:" + err.Error())
			log.Fatalf("import error %s", err)
		}
		return
	}

	setPid(os.Getpid())
	//定时提醒
	go core.StartCron()

	//router
	router.Route(server.Serv())
	server.Run(":8080")
}

func setPid(pid int) {
	d := []byte(strconv.Itoa(pid))
	err := ioutil.WriteFile("./blog.pid", d, 0644)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
}
