package system

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/fifsky/goblog/core/ding"
)

//Config contains application configuration for active gin mode
type Config struct {
	Database      DatabaseConfig
	SessionSecret string `json:"session_secret"`
	DingToken     string `json:"ding_token"`
}

//DatabaseConfig contains database connection info
type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string //database name
	User     string
	Password string
}

//current loaded config
var config *Config

//LoadConfig unmarshals config for current GIN_MODE
func LoadConfig() {
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
