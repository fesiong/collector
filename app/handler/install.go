package handler

import (
	"collector/app/request"
	"collector/config"
	"collector/core"
	"collector/services"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris/v12"
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
	_, err := gorm.Open("mysql", mysqlUrl)
	if err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	config.JsonData.MySQL.Database = req.Database
	config.JsonData.MySQL.User = req.User
	config.JsonData.MySQL.Password = req.Password
	config.JsonData.MySQL.Host = req.Host
	config.JsonData.MySQL.Port = req.Port
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
	services.DB.AutoMigrate(&core.Article{}, &core.ArticleData{}, &core.ArticleSource{})

	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  "采集工具初始化成功",
	})
}
