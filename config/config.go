package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"
)

type configData struct {
	MySQL     mySQLConfig     `json:"mysql"`
	Server    serverConfig    `json:"server"`
	Collector collectorConfig `json:"collector"`
	Content   contentConfig   `json:"content"`
}

var ExecPath string

func InitJSON() {
	sep := string(os.PathSeparator)
	root := filepath.Dir(os.Args[0])
	ExecPath, _ = filepath.Abs(root)
	if strings.Contains(ExecPath, "/T") || strings.Contains(ExecPath, "temp") {
		ExecPath, _ = os.Getwd()
	}
	length := utf8.RuneCountInString(ExecPath)
	lastChar := ExecPath[length-1:]
	if lastChar != sep {
		ExecPath = ExecPath + sep
	}

	//生成public目录
	_, err := os.Stat(ExecPath + "public")
	if err != nil && os.IsNotExist(err) {
		err = os.Mkdir(ExecPath+"public", os.ModePerm)
		if err != nil {
			fmt.Println("无法创建public目录: ", err.Error())
			os.Exit(-1)
		}
	}

	buf, err := ioutil.ReadFile(fmt.Sprintf("%sconfig.json", ExecPath))
	configStr := ""
	if err != nil {
		//文件不存在
		fmt.Println("根目录下不存在配置文件config.json")
		os.Exit(-1)
	}
	configStr = string(buf[:])
	reg := regexp.MustCompile(`/\*.*\*/`)

	configStr = reg.ReplaceAllString(configStr, "")
	buf = []byte(configStr)

	if err := json.Unmarshal(buf, &JsonData); err != nil {
		fmt.Println("配置文件格式有误: ", err.Error())
		os.Exit(-1)
	}

	//load Mysql
	MySQLConfig = JsonData.MySQL
	if MySQLConfig.Database != "" {
		url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			MySQLConfig.User, MySQLConfig.Password, MySQLConfig.Host, MySQLConfig.Port, MySQLConfig.Database, MySQLConfig.Charset)
		MySQLConfig.Url = url
	}

	//load server
	ServerConfig = JsonData.Server
	ServerConfig.ExecPath = ExecPath

	//load collector
	CollectorConfig = loadCollectorConfig(JsonData.Collector)

	//load content
	ContentConfig = JsonData.Content
}

var JsonData configData
var MySQLConfig mySQLConfig
var ServerConfig serverConfig
var CollectorConfig collectorConfig
var ContentConfig contentConfig

func init() {
	InitJSON()
}

func loadCollectorConfig(collector collectorConfig) collectorConfig {
	if collector.ErrorTimes == 0 {
		collector.ErrorTimes = defaultCollectorConfig.ErrorTimes
	}
	if collector.Channels == 0 {
		collector.Channels = defaultCollectorConfig.Channels
	}
	if collector.TitleMinLength == 0 {
		collector.TitleMinLength = defaultCollectorConfig.TitleMinLength
	}
	if collector.ContentMinLength == 0 {
		collector.ContentMinLength = defaultCollectorConfig.ContentMinLength
	}
	for _, v := range defaultCollectorConfig.TitleExclude {
		exists := false
		for _, vv := range collector.TitleExclude {
			if vv == v {
				exists = true
			}
		}
		if !exists {
			collector.TitleExclude = append(collector.TitleExclude, v)
		}
	}
	for _, v := range defaultCollectorConfig.TitleExcludePrefix {
		exists := false
		for _, vv := range collector.TitleExcludePrefix {
			if vv == v {
				exists = true
			}
		}
		if !exists {
			collector.TitleExcludePrefix = append(collector.TitleExcludePrefix, v)
		}
	}
	for _, v := range defaultCollectorConfig.TitleExcludeSuffix {
		exists := false
		for _, vv := range collector.TitleExcludeSuffix {
			if vv == v {
				exists = true
			}
		}
		if !exists {
			collector.TitleExcludeSuffix = append(collector.TitleExcludeSuffix, v)
		}
	}
	for _, v := range defaultCollectorConfig.ContentExclude {
		exists := false
		for _, vv := range collector.ContentExclude {
			if vv == v {
				exists = true
			}
		}
		if !exists {
			collector.ContentExclude = append(collector.ContentExclude, v)
		}
	}
	for _, v := range defaultCollectorConfig.ContentExcludeLine {
		exists := false
		for _, vv := range collector.ContentExcludeLine {
			if vv == v {
				exists = true
			}
		}
		if !exists {
			collector.ContentExcludeLine = append(collector.ContentExcludeLine, v)
		}
	}

	return collector
}

func WriteConfig() error {
	//将现有配置写回文件
	configFile, err := os.OpenFile(fmt.Sprintf("%sconfig.json", ExecPath), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	defer configFile.Close()

	buff := &bytes.Buffer{}

	buf, err := json.MarshalIndent(JsonData, "", "\t")
	if err != nil {
		return err
	}
	buff.Write(buf)

	_, err = io.Copy(configFile, buff)
	if err != nil {
		return err
	}

	return nil
}
