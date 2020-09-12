package route

import (
	"collector/app/handler"
	"collector/config"
	"fmt"
	"github.com/kataras/iris/v12"
)

func Register(app *iris.Application) {
	app.Use(Cors)

	app.OnErrorCode(iris.StatusNotFound, handler.NotFound)
	app.OnErrorCode(iris.StatusInternalServerError, handler.NotFound)

	app.HandleDir("/", fmt.Sprintf("%spublic", config.ExecPath))
	app.Get("/", handler.Inspect, handler.Index)

	app.Get("/install", handler.Install)
	app.Post("/install", handler.InstallForm)

	app.Get("/source", handler.Inspect, handler.ArticleSource)
	app.Get("/article", handler.Inspect, handler.ArticleList)
	app.Get("/keywords", handler.Inspect, handler.Keywords)
	app.Get("/setting", handler.Inspect, handler.DefaultSetting)
	app.Get("/publish", handler.Inspect, handler.PublishSetting)

	app.Post("/setting", handler.InspectJson, handler.DefaultSettingForm)
	app.Post("/publish", handler.InspectJson, handler.PublishSettingForm)

	app.Get("/api/index/echarts", handler.InspectJson, handler.IndexEchartsApi)

	app.Get("/api/article/list", handler.InspectJson, handler.ArticleListApi)
	app.Post("/api/article/delete", handler.InspectJson, handler.ArticleDeleteApi)

	app.Get("/api/article/source/list", handler.InspectJson, handler.ArticleSourceListApi)
	app.Post("/api/article/source/delete", handler.InspectJson, handler.ArticleSourceDeleteApi)
	app.Post("/api/article/source/save", handler.InspectJson, handler.ArticleSourceSaveApi)
	app.Get("/api/setting", handler.InspectJson, handler.DefaultSettingApi)
	app.Get("/api/publish", handler.InspectJson, handler.PublishSettingApi)
}

func Cors(ctx iris.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	if ctx.Request().Method == "OPTIONS" {
		ctx.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,PATCH,OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Api, Accept, Authorization, Version, Token")
		ctx.StatusCode(204)
		return
	}
	ctx.Next()
}
