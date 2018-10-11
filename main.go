package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"

	"github.com/fifsky/goblog/config"
	"github.com/fifsky/goblog/core"
	"github.com/fifsky/goblog/server"
	"github.com/ilibs/gosql"
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
		fmt.Println("Database init success!")
		return
	}

	setPid(os.Getpid())
	//定时提醒
	go core.StartCron()
	server.Run(":8080")
}

func setPid(pid int) {
	d := []byte(strconv.Itoa(pid))
	err := ioutil.WriteFile("./blog.pid", d, 0644)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
}
