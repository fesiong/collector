package route

import (
	"collector/app/handler"
	"collector/config"
	"github.com/kataras/iris/v12"
)

func Register(app *iris.Application) {
	app.Use(Cors)

	app.OnErrorCode(iris.StatusNotFound, handler.NotFound)
	app.OnErrorCode(iris.StatusInternalServerError, handler.NotFound)

	app.HandleDir("/", config.ServerConfig.ExecPath + "public")
	app.Get("/", handler.Inspect, handler.Index)

	app.Get("/install", handler.Install)
	app.Post("/install", handler.InstallForm)

	app.Get("/source", handler.ArticleSource)
	app.Get("/article", handler.ArticleList)
	app.Get("/keywords", handler.Keywords)
	app.Get("/setting", handler.DefaultSetting)
	app.Get("/publish", handler.PublishSetting)

	app.Post("/setting", handler.DefaultSettingForm)
	app.Post("/publish", handler.PublishSettingForm)

	app.Get("/api/index/echarts", handler.IndexEchartsApi)

	app.Get("/api/article/list", handler.ArticleListApi)
	app.Post("/api/article/delete", handler.ArticleDeleteApi)

	app.Get("/api/article/source/list", handler.ArticleSourceListApi)
	app.Post("/api/article/source/delete", handler.ArticleSourceDeleteApi)
	app.Post("/api/article/source/save", handler.ArticleSourceSaveApi)
	app.Get("/api/setting", handler.DefaultSettingApi)
	app.Get("/api/publish", handler.PublishSettingApi)
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
