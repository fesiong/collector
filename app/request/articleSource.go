package request

type ArticleSource struct {
	ID         int    `form:"id"`
	Url        string `form:"url" validate:"required"`
	ErrorTimes int    `form:"error_times"`
	UrlType    int    `form:"url_type"`
}

type Article struct {
	ID int `form:"id"`
}
