package core

import (
	"collector/services"
	"github.com/jinzhu/gorm"
	"time"
)

type Article struct {
	Id           int    `json:"id" gorm:"column:id;type:int(10) unsigned not null AUTO_INCREMENT;primary_key"`
	SourceId     int    `json:"source_id" gorm:"column:source_id;type:int(11) not null;default:0"`
	Title        string `json:"title" gorm:"column:title;type:varchar(190) not null;default:'';index:idx_title"`
	Keywords     string `json:"keywords" gorm:"column:keywords;type:varchar(250) not null;default:''"`
	Description  string `json:"description" gorm:"column:description;type:varchar(250) not null;default:''"`
	Content      string `json:"content" gorm:"-"`
	ArticleType  int    `json:"article_type" gorm:"column:article_type;type:tinyint(1) unsigned not null;default:0;index:idx_article_type"`
	OriginUrl    string `json:"origin_url" gorm:"column:origin_url;type:varchar(250) not null;default:'';index:idx_origin_url"`
	Author       string `json:"author" gorm:"column:author;type:varchar(100) not null;default:''"`
	Views        int    `json:"views" gorm:"column:views;type:int(10) not null;default:0;index:idx_views"`
	Status       int    `json:"status" gorm:"column:status;type:tinyint(1) unsigned not null;default:0;index:idx_status"`
	CreatedTime  int    `json:"created_time" gorm:"column:created_time;type:int(11) unsigned not null;default:0;index:idx_created_time"`
	UpdatedTime  int    `json:"updated_time" gorm:"column:updated_time;type:int(11) unsigned not null;default:0;index:idx_updated_time"`
	DeletedTime  int    `json:"-" gorm:"column:deleted_time;type:int(11) unsigned not null;default:0"`
	OriginDomain string `json:"-" gorm:"-"`
	OriginPath   string `json:"-" gorm:"-"`
	ContentText  string `json:"-" gorm:"-"`
	PubDate      string `json:"-" gorm:"-"`
}

type ArticleData struct {
	Id      int    `json:"id" gorm:"column:id;type:int(10) ;unsigned not null AUTO_INCREMENT;primary_key"`
	Content string `json:"content" gorm:"column:content;type:longtext;not null;default:''"`
}

type ArticleSource struct {
	Id         int    `json:"id" gorm:"column:id;type:int(10) unsigned not null AUTO_INCREMENT;primary_key"`
	Url        string `json:"url" gorm:"column:url;type:varchar(190) not null;default:'';index:idx_url"`
	UrlType    int    `json:"url_type" gorm:"column:url_type;type:tinyint(1) not null;default:0"`
	ErrorTimes int    `json:"error_times" gorm:"column:error_times;type:int(10) not null;default:0;index:idx_error_times"`
}

func (article *Article) Save(db *gorm.DB) error {
	if article.Id == 0 {
		article.CreatedTime = int(time.Now().Unix())
	}

	if err := db.Save(&article).Error; err != nil {
		return err
	}
	articleData := ArticleData{
		Id:      article.Id,
		Content: article.Content,
	}
	db.Save(&articleData)

	return nil
}

func (article *Article) Delete() error {
	db := services.DB
	if err := db.Delete(&article).Error; err != nil {
		return err
	}

	return nil
}

func (source *ArticleSource) Save() error {
	db := services.DB
	if err := db.Save(&source).Error; err != nil {
		return err
	}

	return nil
}

func (source *ArticleSource) Delete() error {
	db := services.DB
	if err := db.Delete(&source).Error; err != nil {
		return err
	}

	return nil
}
