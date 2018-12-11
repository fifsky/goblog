package config

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/fifsky/goblog/ding"
	"github.com/fifsky/goconf"
	"github.com/gin-gonic/gin"
	"github.com/ilibs/gosql"
	"github.com/verystar/logger"
)

type common struct {
	Env            string `json:"env"`
	Debug          string `json:"debug"`
	Path           string `json:"path"`
	ConfigPath     string `json:"config_path"`
	StoragePath    string `json:"storage_path"`
	StaticDomain   string `json:"static_domain"`
	DingToken      string `json:"dingtoken"`
	SessionSecret  string `json:"session_secret"`
	TCaptchaId     string `json:"tcaptcha_id"`
	TCaptchaSecret string `json:"tcaptcha_secret"`
}

type ossConf struct {
	AccessKey    string `json:"access_key"`
	AccessSecret string `json:"access_secret"`
	Endpoint     string `json:"endpoint"`
	Bucket       string `json:"bucket"`
}

type app struct {
	Common    common                   `conf:"common"`
	Log       logger.Config            `conf:"log"`
	DB        map[string]*gosql.Config `conf:"database"`
	OSS       ossConf                  `conf:"oss"`
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
		_, file, _, _ := runtime.Caller(0)
		appPath = filepath.Dir(filepath.Dir(file))
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
