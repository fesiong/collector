package request

type Install struct {
	Database    string `form:"database" validate:"required"`
	User        string `form:"user" validate:"required"`
	Password    string `form:"password" validate:"required"`
	Host        string `form:"host" validate:"required"`
	Port        int    `form:"port" validate:"required"`
	Charset     string `form:"charset"`
	TablePrefix string `form:"table_prefix" validate:"required"`
}
