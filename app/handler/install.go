package handler

import (
	"collector/app/request"
	"collector/config"
	"collector/services"
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris/v12"
	"strings"
)

func Install(ctx iris.Context) {
	if services.DB != nil {
		ctx.Redirect("/")
		return
	}

	ctx.View("install/index.html")
}

func InstallForm(ctx iris.Context) {
	if services.DB != nil {
		ctx.Redirect("/")
		return
	}
	var req request.Install
	if err := ctx.ReadForm(&req); err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	mysqlUrl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		req.User, req.Password, req.Host, req.Port, req.Database, config.MySQLConfig.Charset)
	db, err := gorm.Open("mysql", mysqlUrl)
	if err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}
	//执行数据库初始化操作
	box := packr.New("default", config.ExecPath+"default")
	sql, _ := box.FindString("mysql.sql")
	sql = strings.ReplaceAll(sql, "\r\n", "\n")
	sql = strings.ReplaceAll(sql, "\r", "\n")
	sqlSlice := strings.Split(sql, ";\n")
	for _, v := range sqlSlice {
		if v == "" {
			continue
		}
		db.Exec(v)
	}

	config.JsonData.MySQL.Database = req.Database
	config.JsonData.MySQL.User = req.User
	config.JsonData.MySQL.Password = req.Password
	config.JsonData.MySQL.Host = req.Host
	config.JsonData.MySQL.Port = req.Port
	config.JsonData.MySQL.TablePrefix = req.TablePrefix
	config.JsonData.MySQL.Url = mysqlUrl
	err = config.WriteConfig()
	if err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	config.InitJSON()
	services.InitDB()

	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  "采集工具初始化成功",
	})
}