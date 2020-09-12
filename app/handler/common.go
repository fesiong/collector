package handler

import (
	"collector/config"
	"collector/services"
	"github.com/kataras/iris/v12"
)

func NotFound(ctx iris.Context) {
	ctx.View("errors/404.html")
}

func Inspect(ctx iris.Context) {
	if services.DB == nil {
		ctx.Redirect("/install")
		return
	}

	ctx.Next()
}

func InspectJson(ctx iris.Context) {
	if services.DB == nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  "请先完成初始化操作",
		})
		return
	}

	ctx.Next()
}
