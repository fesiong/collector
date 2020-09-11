package request

type DefaultSetting struct {
	ErrorTimes         int      `form:"error_times"`
	Channels           int      `form:"channels"`
	TitleMinLength     int      `form:"title_min_length"`
	ContentMinLength   int      `form:"content_min_length"`
	TitleExclude       []string `form:"title_exclude[]"`
	TitleExcludePrefix []string `form:"title_exclude_prefix[]"`
	TitleExcludeSuffix []string `form:"title_exclude_suffix[]"`
	ContentExclude     []string `form:"content_exclude[]"`
	ContentExcludeLine []string `form:"content_exclude_line[]"`
}

type ContentSetting struct {
	AutoPublish      int      `form:"auto_publish"`
	TableName        string   `form:"table_name"`
	IdField          string   `form:"id_field"`
	TitleField       string   `form:"title_field"`
	CreatedTimeField string   `form:"created_time_field"`
	KeywordsField    string   `form:"keywords_field"`
	DescriptionField string   `form:"description_field"`
	AuthorField      string   `form:"author_field"`
	ViewsField       string   `form:"views_field"`
	ContentTableName string   `form:"content_table_name"`
	ContentIdField   string   `form:"content_id_field"`
	ContentField     string   `form:"content_field"`
	RemoteUrl        string   `form:"remote_url"`
	ContentType      string   `form:"content_type"`
	Headers          []string `form:"headers[]"`
	Cookies          []string `form:"cookies[]"`
	ExtraFields      []string `form:"extra_fields[]"`
}
