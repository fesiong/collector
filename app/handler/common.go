package handler

import (
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