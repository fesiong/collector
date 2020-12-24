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

	api := app.Party("/api", handler.InspectJson)
	{
		api.Get("/index/echarts", handler.IndexEchartsApi)

		api.Get("/article/list", handler.ArticleListApi)
		api.Post("/article/delete", handler.ArticleDeleteApi)
		api.Post("/article/publish", handler.ArticlePublishApi)
		api.Post("/article/catch", handler.ArticleCatchApi)

		api.Get("/article/source/list", handler.ArticleSourceListApi)
		api.Post("/article/source/delete", handler.ArticleSourceDeleteApi)
		api.Post("/article/source/save", handler.ArticleSourceSaveApi)
		api.Post("/article/source/catch", handler.ArticleSourceCatchApi)
		api.Get("/setting", handler.DefaultSettingApi)
		api.Get("/publish", handler.PublishSettingApi)
	}
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
