package config

type collectorConfig struct {
	ErrorTimes int	`json:"error_times"`
	Channels int `json:"channels"`
	TitleMinLength int `json:"title_min_length"`
	ContentMinLength int `json:"content_min_length"`
	TitleExclude []string `json:"title_exclude"`
	TitleExcludePrefix []string `json:"title_exclude_prefix"`
	TitleExcludeSuffix []string `json:"title_exclude_suffix"`
	ContentExclude []string `json:"content_exclude"`
	ContentExcludeLine []string `json:"content_exclude_line"`
}

var defaultCollectorConfig = collectorConfig{
	ErrorTimes: 5,
	Channels: 5,
	TitleMinLength: 6,
	ContentMinLength: 200,
	TitleExclude: []string{
		"法律声明",
		"站点地图",
		"区长信箱",
		"政务服务",
		"政务公开",
		"领导介绍",
		"首页",
		"当前页",
		"当前位置",
		"来源：",
		"点击：",
		"关注我们",
		"浏览次数",
		"信息分类",
		"索引号",
	},
	TitleExcludePrefix: []string{
		"404",
		"403",
	},
	TitleExcludeSuffix: []string{
		"网",
		"政府",
		"门户",
	},
	ContentExclude: []string{
		"版权声明",
	},
	ContentExcludeLine: []string{
		"背景色：",
		"时间：",
		"作者：",
		"来源：",
		"编辑：",
		"时间:",
		"来源:",
		"作者:",
		"编辑:",
		"摄影：",
		"摄影:",
		"官方微信",
		"一篇：",
		"相关附件",
		"qrcode",
		"微信扫一扫",
		"用手机浏览",
		"打印正文",
		"浏览次数",
	},
}