package config

import (
	"encoding/json"
	"os"
	"fmt"
	"io/ioutil"

	"github.com/ilibs/gosql"
	"github.com/gin-gonic/gin"
	"github.com/fifsky/goblog/core/ding"
	"database/sql"
)

//Config contains application configuration for active gin mode
type Config struct {
	Database       map[string]*gosql.Config
	SessionSecret  string `json:"session_secret"`
	DingToken      string `json:"ding_token"`
}

//current loaded config
var config *Config

//LoadConfig unmarshals config for current GIN_MODE
func LoadConfig() {
	env := os.Getenv("APP_ENV")
	if env == "local" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	config = &Config{}

	file := ""
	switch gin.Mode() {
	case gin.DebugMode:
		file = "dev"
	case gin.ReleaseMode:
		file = "release"
	default:
		panic(fmt.Sprintf("Unknown gin mode %s", gin.Mode()))
	}

	data, err := ioutil.ReadFile("config/config_" + file + ".json")
	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
		panic(err)
	}

	if err := json.Unmarshal(data, config); err != nil {
		panic(err)
	}

	ding.DING_TALK_TOKEN = config.DingToken
}

//GetConfig returns actual config
func GetConfig() *Config {
	return config
}

func ImportDB() ([]sql.Result, error) {
	sqlpath := "./db/blog.sql"
	rst, err := gosql.Import(sqlpath)
	return rst, err
}