package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"unicode/utf8"
)

type configData struct {
	MySQL mySQLConfig
	Server serverConfig
	Collector collectorConfig
	Content contentConfig
}

var ExecPath string

func initJSON() {
	sep := string(os.PathSeparator)
	ExecPath, _ = os.Getwd()
	length := utf8.RuneCountInString(ExecPath)
	lastChar := ExecPath[length-1:]
	if lastChar != sep {
		ExecPath = ExecPath + sep
	}

	bytes, err := ioutil.ReadFile(fmt.Sprintf("%sconfig.json", ExecPath))
	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
		os.Exit(-1)
	}

	configStr := string(bytes[:])
	reg := regexp.MustCompile(`/\*.*\*/`)

	configStr = reg.ReplaceAllString(configStr, "")
	bytes = []byte(configStr)

	if err := json.Unmarshal(bytes, &jsonData); err != nil {
		fmt.Println("Invalid Config: ", err.Error())
		os.Exit(-1)
	}

	//load Mysql
	MySQLConfig = jsonData.MySQL
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		MySQLConfig.User, MySQLConfig.Password, MySQLConfig.Host, MySQLConfig.Port, MySQLConfig.Database, MySQLConfig.Charset)
	MySQLConfig.Url = url

	//load server
	ServerConfig = jsonData.Server
	ServerConfig.ExecPath = ExecPath

	//load collector
	CollectorConfig = jsonData.Collector

	//load content
	ContentConfig = jsonData.Content

	fmt.Println(MySQLConfig, ServerConfig, CollectorConfig, ContentConfig)
}

var jsonData configData
var MySQLConfig mySQLConfig
var ServerConfig serverConfig
var CollectorConfig collectorConfig
var ContentConfig contentConfig

func init() {
	initJSON()
}
