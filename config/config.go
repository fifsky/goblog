package config

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fifsky/goblog/ding"
	"github.com/fifsky/goconf"
	"github.com/gin-gonic/gin"
	"github.com/ilibs/gosql"
	"github.com/ilibs/logger"
)

type common struct {
	Env           string `json:"env"`
	Debug         string `json:"debug"`
	Path          string `json:"path"`
	ConfigPath    string `json:"config_path"`
	StoragePath   string `json:"storage_path"`
	DingToken     string `json:"dingtoken"`
	SessionSecret string `json:"session_secret"`
}

type app struct {
	Common    common                   `conf:"common"`
	Log       logger.Config            `conf:"log"`
	DB        map[string]*gosql.Config `conf:"database"`
	StartTime time.Time
}

var App = &app{
	StartTime: time.Now(),
}

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
			src := "/goblog"
			appPath = execpath[0 : strings.Index(execpath, src)+len(src)]
		}
	}

	App.Common.Path = appPath
	App.Common.ConfigPath = filepath.Join(appPath, "config")

	conf, err := goconf.NewConfig(filepath.Join(App.Common.ConfigPath, env))
	if err != nil {
		logger.Fatal("json config path error %s", err.Error())
	}

	//load config
	if err := conf.Load(App); err != nil {
		log.Fatal("Config Error:", err)
	}

	if !filepath.IsAbs(App.Common.StoragePath) {
		App.Common.StoragePath = filepath.Join(appPath, App.Common.StoragePath)
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
		c.LogPath = filepath.Join(App.Common.StoragePath, "logs")
		c.LogSentryDSN = App.Log.LogSentryDSN
		c.LogSentryType = App.Log.LogSentryType
		c.LogDetail = App.Log.LogDetail
	})

	//test
	if os.Getenv("MYSQL_TEST_DSN") != "" {
		App.DB["default"].Dsn = os.Getenv("MYSQL_TEST_DSN")
	}
}

func ImportDB() ([]sql.Result, error) {
	sqlpath := "./config/blog.sql"
	rst, err := gosql.Import(sqlpath)
	return rst, err
}
