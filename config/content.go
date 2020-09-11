package config

type KeyValue struct {
	Key   string `json:"key" form:"key"`
	Value string `json:"value" form:"value"`
}

type contentConfig struct {
	AutoPublish      int        `json:"auto_publish"`
	TableName        string     `json:"table_name"`
	IdField          string     `json:"id_field"`
	TitleField       string     `json:"title_field"`
	CreatedTimeField string     `json:"created_time_field"`
	KeywordsField    string     `json:"keywords_field"`
	DescriptionField string     `json:"description_field"`
	AuthorField      string     `json:"author_field"`
	ViewsField       string     `json:"views_field"`
	ContentTableName string     `json:"content_table_name"`
	ContentIdField   string     `json:"content_id_field"`
	ContentField     string     `json:"content_field"`
	RemoteUrl        string     `json:"remote_url"`
	ContentType      string     `json:"content_type"`
	Headers          []KeyValue `json:"headers"`
	Cookies          []KeyValue `json:"cookies"`
	ExtraFields      []KeyValue `json:"extra_fields"`
}
