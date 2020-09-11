package provider

import (
	"collector"
	"collector/services"
)

func GetArticleSourceList(currentPage int, pageSize int) ([]collector.ArticleSource, int, error) {
	var sources []collector.ArticleSource
	offset := (currentPage - 1) * pageSize
	var total int

	builder := services.DB.Model(collector.ArticleSource{}).Order("id desc")
	if err := builder.Count(&total).Limit(pageSize).Offset(offset).Find(&sources).Error; err != nil {
		return nil, 0, err
	}

	return sources, total, nil
}

func GetArticleList(currentPage int, pageSize int) ([]collector.Article, int, error) {
	var articles []collector.Article
	offset := (currentPage - 1) * pageSize
	var total int

	builder := services.DB.Model(collector.Article{}).Order("id desc")
	if err := builder.Count(&total).Limit(pageSize).Offset(offset).Find(&articles).Error; err != nil {
		return nil, 0, err
	}
	if len(articles) > 0 {
		for i, v := range articles {
			var articleData collector.ArticleData
			if err := services.DB.Model(collector.ArticleData{}).Where("`id` = ?", v.Id).First(&articleData).Error; err == nil {
				articles[i].Content = articleData.Content
			}
		}
	}
	return articles, total, nil
}

func GetArticleById(id int) (*collector.Article, error) {
	var article collector.Article
	if err := services.DB.Model(collector.Article{}).Where("`id` = ?", id).First(&article).Error; err != nil {
		return nil, err
	}
	var articleData collector.ArticleData
	if err := services.DB.Model(collector.ArticleData{}).Where("`id` = ?", id).First(&articleData).Error; err != nil {
		return nil, err
	}
	article.Content = articleData.Content

	return &article, nil
}

func GetArticleSourceById(id int) (*collector.ArticleSource, error) {
	var source collector.ArticleSource
	if err := services.DB.Model(collector.ArticleSource{}).Where("`id` = ?", id).First(&source).Error; err != nil {
		return nil, err
	}

	return &source, nil
}

func GetArticleSourceByUrl(uri string) (*collector.ArticleSource, error) {
	var source collector.ArticleSource
	if err := services.DB.Model(collector.ArticleSource{}).Where("`url` = ?", uri).First(&source).Error; err != nil {
		return nil, err
	}

	return &source, nil
}