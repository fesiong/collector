package provider

import (
	"collector/core"
	"collector/services"
)

func GetArticleSourceList(currentPage int, pageSize int) ([]core.ArticleSource, int, error) {
	var sources []core.ArticleSource
	offset := (currentPage - 1) * pageSize
	var total int

	builder := services.DB.Model(core.ArticleSource{}).Order("id desc")
	if err := builder.Count(&total).Limit(pageSize).Offset(offset).Find(&sources).Error; err != nil {
		return nil, 0, err
	}

	return sources, total, nil
}

func GetArticleList(currentPage int, pageSize int) ([]core.Article, int, error) {
	var articles []core.Article
	offset := (currentPage - 1) * pageSize
	var total int

	builder := services.DB.Model(core.Article{}).Order("id desc")
	if err := builder.Count(&total).Limit(pageSize).Offset(offset).Find(&articles).Error; err != nil {
		return nil, 0, err
	}
	if len(articles) > 0 {
		for i, v := range articles {
			var articleData core.ArticleData
			if err := services.DB.Model(core.ArticleData{}).Where("`id` = ?", v.Id).First(&articleData).Error; err == nil {
				articles[i].Content = articleData.Content
			}
		}
	}
	return articles, total, nil
}

func GetArticleById(id int) (*core.Article, error) {
	var article core.Article
	if err := services.DB.Model(core.Article{}).Where("`id` = ?", id).First(&article).Error; err != nil {
		return nil, err
	}
	var articleData core.ArticleData
	if err := services.DB.Model(core.ArticleData{}).Where("`id` = ?", id).First(&articleData).Error; err != nil {
		return nil, err
	}
	article.Content = articleData.Content

	return &article, nil
}

func GetArticleSourceById(id int) (*core.ArticleSource, error) {
	var source core.ArticleSource
	if err := services.DB.Model(core.ArticleSource{}).Where("`id` = ?", id).First(&source).Error; err != nil {
		return nil, err
	}

	return &source, nil
}

func GetArticleSourceByUrl(uri string) (*core.ArticleSource, error) {
	var source core.ArticleSource
	if err := services.DB.Model(core.ArticleSource{}).Where("`url` = ?", uri).First(&source).Error; err != nil {
		return nil, err
	}

	return &source, nil
}
