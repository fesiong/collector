package handler

import (
	"collector/app/provider"
	"collector/app/request"
	"collector/config"
	"collector/core"
	"github.com/kataras/iris/v12"
)

func Keywords(ctx iris.Context) {
	ctx.View("article/keywords.html")
}

func ArticleSource(ctx iris.Context) {
	ctx.View("article/source.html")
}

func ArticleList(ctx iris.Context) {
	ctx.View("article/list.html")
}

func ArticleListApi(ctx iris.Context) {
	currentPage := ctx.URLParamIntDefault("page", 1)
	pageSize := ctx.URLParamIntDefault("limit", 20)

	articleList, total, err := provider.GetArticleList(currentPage, pageSize)
	if err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	ctx.JSON(iris.Map{
		"code":  config.StatusOK,
		"msg":   "",
		"data":  articleList,
		"count": total,
	})
}

func ArticleDeleteApi(ctx iris.Context) {
	var req request.Article
	if err := ctx.ReadForm(&req); err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	article, err := provider.GetArticleById(req.ID)
	if err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	err = article.Delete()
	if err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  "删除成功",
	})
}

func ArticleSourceListApi(ctx iris.Context) {
	currentPage := ctx.URLParamIntDefault("page", 1)
	pageSize := ctx.URLParamIntDefault("limit", 20)

	sourceList, total, err := provider.GetArticleSourceList(currentPage, pageSize)
	if err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	ctx.JSON(iris.Map{
		"code":  config.StatusOK,
		"msg":   "",
		"data":  sourceList,
		"count": total,
	})
}

func ArticleSourceDeleteApi(ctx iris.Context) {
	var req request.ArticleSource
	if err := ctx.ReadForm(&req); err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	source, err := provider.GetArticleSourceById(req.ID)
	if err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	err = source.Delete()
	if err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  "删除成功",
	})
}

func ArticleSourceSaveApi(ctx iris.Context) {
	var req request.ArticleSource
	err := ctx.ReadForm(&req)
	if err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}
	var source *core.ArticleSource
	if req.ID > 0 {
		source, err = provider.GetArticleSourceById(req.ID)
		if err != nil {
			ctx.JSON(iris.Map{
				"code": config.StatusFailed,
				"msg":  err.Error(),
			})
			return
		}
	} else {
		source, err = provider.GetArticleSourceByUrl(req.Url)
		if err == nil {
			ctx.JSON(iris.Map{
				"code": config.StatusFailed,
				"msg":  "该数据源已存在，不用重复添加",
			})
			return
		}
		source = &core.ArticleSource{}
		source.Url = req.Url
	}

	if req.Url != "" {
		source.Url = req.Url
	}
	source.ErrorTimes = req.ErrorTimes

	err = source.Save()
	if err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}
	//添加完，马上抓取
	core.GetArticleLinks(source)

	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  "添加/修改成功",
		"data": source,
	})
}

func ArticlePublishApi(ctx iris.Context) {
	var req request.Article
	if err := ctx.ReadForm(&req); err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	article, err := provider.GetArticleById(req.ID)
	if err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	core.AutoPublish(article)

	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  "删除成功",
	})
}

func ArticleCatchApi(ctx iris.Context) {
	var req request.Article
	if err := ctx.ReadForm(&req); err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	article, err := provider.GetArticleById(req.ID)
	if err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	go core.GetArticleDetail(article)

	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  "抓取任务已执行",
	})
}

func ArticleSourceCatchApi(ctx iris.Context) {
	var req request.ArticleSource
	if err := ctx.ReadForm(&req); err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	source, err := provider.GetArticleSourceById(req.ID)
	if err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	go core.GetArticleLinks(source)

	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  "抓取任务执行",
	})
}