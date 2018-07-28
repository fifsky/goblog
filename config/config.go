package config

import (
	"os"
	"log"
	"path/filepath"
	"database/sql"

	"github.com/ilibs/gosql"
	"github.com/fifsky/goblog/core/ding"
	"github.com/ilibs/logger"
	"github.com/fifsky/goconf"
	"github.com/gin-gonic/gin"
	"strings"
)

type common struct {
	Env         string `json:"env"`
	Debug       string `json:"debug"`
	ConfigPath  string `json:"config_path"`
	StoragePath string `json:"storage_path"`
	DingToken   string `json:"dingtoken"`
	SessionSecret string `json:"session_secret"`
}

type app struct {
	Common common                   `conf:"common"`
	Log    logger.Config            `conf:"log"`
	DB     map[string]*gosql.Config `conf:"database"`
}

var App = &app{}

func init() {
	argsInit()
	Load(ExtArgs)
}

func Load(args map[string]string) {
	//env
	env := os.Getenv("APP_ENV")
	if env == "local" {
		gin.SetMode(gin.DebugMode)
	} else {
		env = "prod"
		gin.SetMode(gin.ReleaseMode)
	}

	appPath := args["config"]
	if appPath == "" {
		//获得程序路径从里面获取到goblog的路径
		execpath, err := os.Getwd()
		if err == nil {
			appPath = execpath[0: strings.Index(execpath, "/goblog")+7]
		}
	}

	App.Common.ConfigPath = filepath.Join(appPath, "config") + "/"

	conf, err := goconf.NewConfig(App.Common.ConfigPath + env)
	if err != nil {
		logger.Fatal("json config path error %s", err.Error())
	}

	//load config
	if err := conf.Load(App); err != nil {
		log.Fatal("Config Error:", err)
	}

	if !filepath.IsAbs(App.Common.StoragePath) {
		App.Common.StoragePath = filepath.Join(appPath, App.Common.StoragePath) + "/"
	}

	ding.DING_TALK_TOKEN = App.Common.DingToken

	//debug model
	if args["debug"] != "" {
		App.Common.Debug = args["debug"]
	}

	//debug
	if App.Common.Debug == "on" {
		//log level
		App.Log.LogLevel = "debug"
		//log model
		App.Log.LogMode = "std"
	}

	if args["show-sql"] == "on" {
		for _, d := range App.DB {
			d.ShowSql = true
		}
	} else if args["show-sql"] == "off" {
		for _, d := range App.DB {
			d.ShowSql = false
		}
	}

	logger.Setting(func(c *logger.Config) {
		c.LogMode = App.Log.LogMode
		c.LogLevel = App.Log.LogLevel
		c.LogMaxFiles = App.Log.LogMaxFiles
		c.LogPath = App.Common.StoragePath + "logs/"
		c.LogSentryDSN = App.Log.LogSentryDSN
		c.LogSentryType = App.Log.LogSentryType
	})
}

func ImportDB() ([]sql.Result, error) {
	sqlpath := "./db/blog.sql"
	rst, err := gosql.Import(sqlpath)
	return rst, err
}
