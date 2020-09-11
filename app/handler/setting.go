package handler

import (
	"collector/app/request"
	"collector/config"
	"github.com/kataras/iris/v12"
	"strings"
)

func DefaultSetting(ctx iris.Context) {
	ctx.View("setting/index.html")
}

func PublishSetting(ctx iris.Context) {

	ctx.View("setting/publish.html")
}

func DefaultSettingApi(ctx iris.Context) {
	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  "",
		"data": config.JsonData.Collector,
	})
}

func PublishSettingApi(ctx iris.Context) {
	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  "",
		"data": config.JsonData.Content,
	})
}

func DefaultSettingForm(ctx iris.Context) {
	var req request.DefaultSetting
	if err := ctx.ReadForm(&req); err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	config.JsonData.Collector.ErrorTimes = req.ErrorTimes
	config.JsonData.Collector.Channels = req.Channels
	config.JsonData.Collector.TitleMinLength = req.TitleMinLength
	config.JsonData.Collector.ContentMinLength = req.ContentMinLength
	config.JsonData.Collector.TitleExclude = req.TitleExclude
	config.JsonData.Collector.TitleExcludePrefix = req.TitleExcludePrefix
	config.JsonData.Collector.TitleExcludeSuffix = req.TitleExcludeSuffix
	config.JsonData.Collector.ContentExclude = req.ContentExclude
	config.JsonData.Collector.ContentExcludeLine = req.ContentExcludeLine

	err := config.WriteConfig()
	if err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	config.InitJSON()

	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  "配置成功",
	})
}

func PublishSettingForm(ctx iris.Context) {
	var req request.ContentSetting
	if err := ctx.ReadForm(&req); err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	config.JsonData.Content.AutoPublish = req.AutoPublish
	config.JsonData.Content.TableName = req.TableName
	config.JsonData.Content.IdField = req.IdField
	config.JsonData.Content.TitleField = req.TitleField
	config.JsonData.Content.CreatedTimeField = req.CreatedTimeField
	config.JsonData.Content.KeywordsField = req.KeywordsField
	config.JsonData.Content.DescriptionField = req.DescriptionField
	config.JsonData.Content.ContentTableName = req.ContentTableName
	config.JsonData.Content.ContentIdField = req.ContentIdField
	config.JsonData.Content.ContentField = req.ContentField
	config.JsonData.Content.RemoteUrl = req.RemoteUrl
	config.JsonData.Content.ContentType = req.ContentType

	var headers []config.KeyValue
	for _, v := range req.Headers {
		vv := strings.Split(v, ":")
		if len(vv) >= 2 {
			headers = append(headers, config.KeyValue{
				Key:   vv[0],
				Value: vv[1],
			})
		}
	}
	var cookies []config.KeyValue
	for _, v := range req.Cookies {
		vv := strings.Split(v, ":")
		if len(vv) >= 2 {
			cookies = append(cookies, config.KeyValue{
				Key:   vv[0],
				Value: vv[1],
			})
		}
	}
	var extraFields []config.KeyValue
	for _, v := range req.ExtraFields {
		vv := strings.Split(v, ":")
		if len(vv) >= 2 {
			extraFields = append(extraFields, config.KeyValue{
				Key:   vv[0],
				Value: vv[1],
			})
		}
	}

	config.JsonData.Content.Headers = headers
	config.JsonData.Content.Cookies = cookies
	config.JsonData.Content.ExtraFields = extraFields

	err := config.WriteConfig()
	if err != nil {
		ctx.JSON(iris.Map{
			"code": config.StatusFailed,
			"msg":  err.Error(),
		})
		return
	}

	config.InitJSON()

	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  "配置成功",
	})
}
