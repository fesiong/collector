package collector

import (
	"collector/config"
	"collector/library"
	"collector/service"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/parnurzeal/gorequest"
	"github.com/robfig/cron/v3"
	"golang.org/x/net/html/charset"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

var waitGroup sync.WaitGroup
var ch chan string

func Crond() {
	//一次使用几个通道
	ch = make(chan string, config.CollectorConfig.Channels)

	fmt.Println("collection")
	//publishData := map[string]string{}
	//publishData[config.ContentConfig.TitleField] = "这是一个标题"
	//publishData[config.ContentConfig.KeywordsField] = "这是一个关键词"
	//publishDataKeys := make([]string, len(publishData))
	//publishDataValues := make([]string, len(publishData))
	//j := 0
	//for k,v := range publishData {
	//	publishDataKeys[j] = k
	//	publishDataValues[j] = fmt.Sprintf("'%s'", v)
	//	j++
	//}
	//
	////var id struct{
	////	Id int
	////}
	//inserId := int64(0)
	//result, err := services.DB.DB().Exec(fmt.Sprintf("INSERT INTO `%s` (%s)VALUES(%s)", config.ContentConfig.TableName, strings.Join(publishDataKeys, ","), strings.Join(publishDataValues, ",")))
	//if err == nil {
	//	inserId, err = result.LastInsertId()
	//	if config.ContentConfig.TableName != config.ContentConfig.ContentTableName {
	//		dd := services.DB.Exec(fmt.Sprintf("INSERT INTO `%s` (%s, %s)VALUES(?, ?)", config.ContentConfig.ContentTableName, config.ContentConfig.ContentIdField, config.ContentConfig.ContentField), inserId, "这些一些content内容")
	//		fmt.Println(dd.RowsAffected)
	//	}
	//}
	//
	//fmt.Println(inserId,err, publishData, publishDataKeys, publishDataValues)
	//link := &Article{
	//	OriginUrl: "https://mbd.baidu.com/newspage/data/landingsuper?context=%7B%22nid%22%3A%22news_9278787920922000269%22%7D&n_type=0&p_from=1",
	//}
	//
	//_ = CollectDetail(link)
	//fmt.Println(link.Title, "--------", link.Content)
	//os.Exit(0)
	//1小时运行一次，采集地址，加入到地址池
	//每分钟运行一次，检查是否有需要采集的文章s
	crontab := cron.New(cron.WithSeconds())
	//10分钟抓一次列表
	crontab.AddFunc("1 */10 * * * *", CollectListTask)
	//1分钟抓一次详情
	crontab.AddFunc("1 */1 * * * *", CollectDetailTask)
	crontab.Start()
	//启动的时候，先执行一遍
	go CollectListTask()
	go CollectDetailTask()

	select {}
}

func CollectListTask() {
	db := services.DB
	var articleSources []ArticleSource
	err := db.Model(ArticleSource{}).Where("`error_times` < ?", config.CollectorConfig.ErrorTimes).Find(&articleSources).Error
	if err != nil {
		return
	}

	for _, v := range articleSources {
		//ch <- fmt.Sprintf("%d", i)
		//waitGroup.Add(1)
		GetArticleLinks(v)
	}

	//waitGroup.Wait()
}

func CollectDetailTask() {
	fmt.Println("collect detail")
	//检查article的地址
	var articleList []Article

	db := services.DB
	db.Debug().Model(Article{}).Where("status = 0").Order("id asc").Limit(config.CollectorConfig.Channels * 100).Scan(&articleList)
	for _, vv := range articleList {
		ch <- vv.OriginUrl
		waitGroup.Add(1)
		go GetArticleDetail(vv)
	}

	waitGroup.Wait()
}

func GetArticleLinks(v ArticleSource) {
	//defer func() {
	//	waitGroup.Done()
	//	<-ch
	//}()
	db := services.DB
	articleList, err := CollectLinks(v.Url)
	if err == nil {
		for _, article := range articleList {
			//先检查数据库里有没有，没有的话，就抓回来
			article.CreatedTime = int(time.Now().Unix())
			article.ArticleType = v.UrlType
			article.Status = 0
			db.Model(Article{}).Where(Article{OriginUrl: article.OriginUrl}).FirstOrCreate(&article)
		}
	} else {
		db.Model(&v).Update("error_times", v.ErrorTimes+1)
	}
}

func GetArticleDetail(v Article) {
	defer func() {
		waitGroup.Done()
		<-ch
	}()

	db := services.DB
	//标记当前为执行中
	db.Model(Article{}).Where("`id` = ?", v.Id).Update("status", 2)

	_ = CollectDetail(&v)

	//更新到数据库中
	status := int(1)
	if v.Content == "" {
		status = 3
	}
	if utf8.RuneCountInString(v.Title) < 10 {
		status = 3
	}
	urlArr := strings.Split(v.OriginUrl, "/")
	if len(urlArr) <= 3 {
		status = 3
	}
	if len(urlArr) <= 4 && strings.HasPrefix(v.OriginUrl, "/") {
		status = 3
	}

	if strings.Contains(v.Title, "法律声明") || strings.Contains(v.Title, "关于我们") || strings.Contains(v.Title, "站点地图") || strings.Contains(v.Title, "区长信箱") || strings.Contains(v.Title, "政务服务网") || strings.Contains(v.Title, "政务公开") || strings.Contains(v.Title, "人民政府网站") || strings.Contains(v.Title, "门户网站") || strings.Contains(v.Title, "领导介绍") || strings.Contains(v.Title, "403") || strings.Contains(v.Title, "404") || strings.Contains(v.Title, "Government") || strings.Contains(v.Title, "China") {
		status = 3
	}
	//小于500字 内容，不过审
	if utf8.RuneCountInString(v.ContentText) < 200 {
		status = 3
	}
	if strings.Contains(v.ContentText, "ICP备") || strings.Contains(v.ContentText, "政府网站标识码") || strings.Contains(v.ContentText, "以上版本浏览本站") || strings.Contains(v.ContentText, "版权声明") || strings.Contains(v.ContentText, "公网安备") {
		status = 3
	}

	db.Model(Article{}).Where("`id` = ?", v.Id).Update("status", status)

	timeTemplate1 := "2006-01-02 15:04:05"
	timestamp := int(time.Now().Unix())
	pubTime, _ := time.ParseInLocation(timeTemplate1, v.PubDate, time.Local)
	if pubTime.Unix() > 0 {
		timestamp = int(pubTime.Unix())
	}

	v.UpdatedTime = int(time.Now().Unix())
	v.CreatedTime = timestamp
	v.Status = status

	article := v
	fmt.Println(status, v.Title, v.OriginUrl)
	article.Save(db)
	//如果配置了自动发布，则将它发布到指定表中
	if config.ContentConfig.AutoPublish {
		publishData := map[string]string{
			config.ContentConfig.TitleField: article.Title,
		}
		if config.ContentConfig.KeywordsField != "" {
			publishData[config.ContentConfig.KeywordsField] = article.Keywords
		}
		if config.ContentConfig.DescriptionField != "" {
			publishData[config.ContentConfig.DescriptionField] = article.Description
		}
		if config.ContentConfig.CreatedTimeField != "" {
			publishData[config.ContentConfig.CreatedTimeField] = strconv.Itoa(article.CreatedTime)
		}
		if config.ContentConfig.AuthorField != "" {
			publishData[config.ContentConfig.AuthorField] = article.Author
		}
		if config.ContentConfig.ViewsField != "" {
			publishData[config.ContentConfig.ViewsField] = strconv.Itoa(article.Views)
		}
		if config.ContentConfig.TableName == config.ContentConfig.ContentTableName || config.ContentConfig.ContentTableName == "" {
			if config.ContentConfig.ContentField != "" {
				publishData[config.ContentConfig.ContentField] = article.Content
			}
		}

		publishDataKeys := make([]string, len(publishData))
		publishDataValues := make([]string, len(publishData))
		j := 0
		for k,v := range publishData {
			publishDataKeys[j] = k
			publishDataValues[j] = fmt.Sprintf("'%s'", v)
			j++
		}

		insertId := int64(0)
		result, err := services.DB.DB().Exec(fmt.Sprintf("INSERT INTO `%s` (%s)VALUES(%s)", config.ContentConfig.TableName, strings.Join(publishDataKeys, ","), strings.Join(publishDataValues, ",")))
		if err == nil {
			insertId, err = result.LastInsertId()
			if config.ContentConfig.ContentTableName != "" && config.ContentConfig.TableName != config.ContentConfig.ContentTableName {
				services.DB.Exec(fmt.Sprintf("INSERT INTO `%s` (%s, %s)VALUES(?, ?)", config.ContentConfig.ContentTableName, config.ContentConfig.ContentIdField, config.ContentConfig.ContentField), insertId, article.Content)
			}
		}
	}
}

func CollectLinks(link string) ([]Article, error) {
	resp, body, errs := gorequest.New().Timeout(5 * time.Second).Get(link).End()
	if errs != nil {
		fmt.Println(errs)
		return nil, errs[0]
	}
	defer resp.Body.Close()
	//检查网站编码，并转成uft8
	var htmlEncode string
	//先判断标识
	reg := regexp.MustCompile(`(?i:)charset="?([a-z0-9\-]*)`)
	match := reg.FindStringSubmatch(body)
	if len(match) > 1 {
		htmlEncode = strings.ToLower(match[1])
	}
	//再判断标题
	if htmlEncode == "" {
		reg := regexp.MustCompile(`(?i:)<title[^>]*>(.*?)<\/title>`)
		match := reg.FindStringSubmatch(body)
		if len(match) > 1 {
			title := match[1]
			_, htmlEncode, _ = charset.DetermineEncoding([]byte(title), "")
		}
	}
	if htmlEncode == "gb2312" || htmlEncode == "gbk" || htmlEncode == "big5" {
		body = library.ConvertToString(body, "gbk", "utf-8")
	}
	htmlR := strings.NewReader(body)
	doc, err := goquery.NewDocumentFromReader(htmlR)
	if err != nil {
		return nil, err
	}

	var articles []Article
	aLinks := doc.Find("a")
	//读取所有连接
	for i := range aLinks.Nodes {
		href, exists := aLinks.Eq(i).Attr("href")
		title := strings.TrimSpace(aLinks.Eq(i).Text())
		if exists {
			href = ParseLink(href, link)
		}
		if len(href) > 250 {
			href = string(href[:250])
		}
		//斜杠/结尾的抛弃
		//if strings.HasSuffix(href, "/") == false {
		articles = append(articles, Article{
			Title: title,
			OriginUrl:  href,
		})
		//}
	}

	return articles, nil
}

func ParseLink(link string, baseUrl string) string {
	if !strings.HasSuffix(baseUrl, "/") {
		baseUrl += "/"
	}
	if strings.Contains(link, "javascript") || strings.Contains(link, "void") || link == "#" || link == "./" || link == "../" || link == "../../" {
		return ""
	}

	link = replaceDot(link, baseUrl)

	return link
}

func replaceDot(currUrl string, baseUrl string) string {
	if strings.HasPrefix(currUrl, "//") {
		currUrl = fmt.Sprintf("https:%s", currUrl)
	}
	urlInfo, err := url.Parse(currUrl)
	if err != nil {
		return ""
	}
	if urlInfo.Scheme != "" {
		return currUrl
	}
	baseInfo, err := url.Parse(baseUrl)
	if err != nil {
		return ""
	}

	u := baseInfo.Scheme + "://" + baseInfo.Host
	var path string
	if strings.Index(urlInfo.Path, "/") == 0 {
		path = urlInfo.Path
	} else {
		path = filepath.Dir(baseInfo.Path) + "/" + urlInfo.Path
	}

	rst := make([]string, 0)
	pathArr := strings.Split(path, "/")

	// 如果path是已/开头，那在rst加入一个空元素
	if pathArr[0] == "" {
		rst = append(rst, "")
	}
	for _, p := range pathArr {
		if p == ".." {
			if len(rst) > 0 {
				if rst[len(rst)-1] == ".." {
					rst = append(rst, "..")
				} else {
					rst = rst[:len(rst)-1]
				}
			}
		} else if p != "" && p != "." {
			rst = append(rst, p)
		}
	}
	return u + strings.Join(rst, "/")
}

func CollectDetail(article *Article) error {
	resp, body, errs := gorequest.New().Timeout(5 * time.Second).Get(article.OriginUrl).End()
	if errs != nil {
		return errs[0]
	}
	defer resp.Body.Close()
	//检查网站编码，并转成uft8
	//_, htmlEncode, _ := charset.DetermineEncoding([]byte(body), "")
	//fmt.Println( htmlEncode)
	var htmlEncode string
	reg := regexp.MustCompile(`(?i:)<title[^>]*>(.*?)<\/title>`)
	match := reg.FindStringSubmatch(body)
	if len(match) > 1 {
		aa := match[1]
		_, htmlEncode, _ = charset.DetermineEncoding([]byte(aa), "")
		if htmlEncode != "utf-8" {
			body = library.ConvertToString(body, "gbk", "utf-8")
		}
	}
	//先删除一些不必要的标签
	re, _ := regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	body = re.ReplaceAllString(body, "")
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	body = re.ReplaceAllString(body, "")

	htmlR := strings.NewReader(body)
	doc, err := goquery.NewDocumentFromReader(htmlR)
	if err != nil {
		return err
	}

	//获取前缀
	article.GetDomain()

	//如果是百度百科地址，单独处理
	if strings.Contains(article.OriginUrl, "baike.baidu.com") {
		article.ParseBaikeDetail(doc, body)
	} else {
		article.ParseNormalDetail(doc, body)
	}
	nameRune := []rune(article.Description)
	curLen := len(nameRune)
	if curLen > 150 {
		article.Description = string(nameRune[:150])
	}

	return nil
}

func (article *Article) ParseBaikeDetail(doc *goquery.Document, body string) {
	//获取标题
	article.Title = doc.Find("h1").Text()
	//获取描述
	reg := regexp.MustCompile(`<meta\s+name="description"\s+content="([^"]+)">`)
	match := reg.FindStringSubmatch(body)
	if len(match) > 1 {
		article.Description = match[1]
	}
	//获取关键词
	reg = regexp.MustCompile(`<meta\s+name="keywords"\s+content="([^"]+)">`)
	match = reg.FindStringSubmatch(body)
	if len(match) > 1 {
		article.Keywords = match[1]
	}

	doc.Find(".edit-icon").Remove()
	contentList := doc.Find(".para-title,.para")
	content := ""
	for i := range contentList.Nodes {
		content += "<p>" + contentList.Eq(i).Text() + "</p>"
	}

	article.Content = content
}

func (article *Article) ParseNormalDetail(doc *goquery.Document, body string) {
	article.ParseTitle(doc, body)

	//尝试获取正文内容
	article.ParseContent(doc, body)

	//尝试获取作者
	reg := regexp.MustCompile(`<meta\s+name="Author"\s+content="(.*?)"[^>]*>`)
	match := reg.FindStringSubmatch(body)
	if len(match) > 1 {
		author := match[1]
		if author == "" {
			reg := regexp.MustCompile(`(?i)(来源|作者)\s*(:|：|\s)\s*([^\s]+)`)
			match := reg.FindStringSubmatch(body)
			if len(match) > 1 {
				author = match[3]
			}
		}
		article.Author = author
	}

	//尝试获取法布时间
	reg = regexp.MustCompile(`(?i)<meta\s+name="PubDate"\s+content="(.*?)"[^>]*>`)
	match = reg.FindStringSubmatch(body)
	if len(match) > 1 {
		pubDate := match[1]
		if pubDate == "" {
			reg = regexp.MustCompile(`(?i)([0-9]{4})\s*[\-|\/|年]\s*([0-9]{1,2})\s*[\-|\/|月]\s*([0-9]{1,2})\s*([\-|\/|日])?\s*(([0-9]{1,2})\s*[:|：|时]\s*([0-9]{1,2})\s*([:|：|分])?\s*([0-9]{1,2})?)?`)
			match = reg.FindStringSubmatch(body)
			if len(match) > 1 {
				if match[1] != "" {
					pubDate = match[1] + "-" + match[2] + "-" + match[3]
				}
				if match[5] != "" {
					pubDate += " " + match[6] + ":" + match[7]
					if match[9] != "" {
						pubDate += ":" + match[9]
					} else {
						pubDate += ":00"
					}
				} else {
					pubDate += " 12:00:00"
				}
			}
		}
		article.PubDate = pubDate
	}
}

func (article *Article) ParseTitle(doc *goquery.Document, body string) {
	//尝试获取标题
	//先尝试获取h1标签
	title := ""
	h1s := doc.Find("h1")
	if h1s.Length() > 0 {
		for i := range h1s.Nodes {
			item := h1s.Eq(i)
			item.Children().Remove()
			text := strings.TrimSpace(item.Text())
			textLen := utf8.RuneCountInString(text)
			if textLen >= config.CollectorConfig.TitleMinLength && textLen > utf8.RuneCountInString(title) && !HasContain(text, config.CollectorConfig.TitleExclude) && !HasPrefix(text, config.CollectorConfig.TitleExcludePrefix) && !HasSuffix(text, config.CollectorConfig.TitleExcludeSuffix) {
				title = text
			}
		}
	}
	if title == "" {
		//获取 政府网站的 <meta name='ArticleTitle' content='西城法院出台案件在线办理操作指南'>
		text, exist := doc.Find("meta[name=ArticleTitle]").Attr("content")
		if exist {
			text = strings.TrimSpace(text)
			if utf8.RuneCountInString(text) >= config.CollectorConfig.TitleMinLength && !HasContain(text, config.CollectorConfig.TitleExclude) && !HasPrefix(text, config.CollectorConfig.TitleExcludePrefix) && !HasSuffix(text, config.CollectorConfig.TitleExcludeSuffix) {
				title = text
			}
		}
	}
	if title == "" {
		//获取title标签
		text := doc.Find("title").Text()
		text = strings.ReplaceAll(text, "_", "-")
		sepIndex := strings.Index(text, "-")
		if sepIndex > 0 {
			text = text[:sepIndex]
		}
		text = strings.TrimSpace(text)
		if utf8.RuneCountInString(text) >= config.CollectorConfig.TitleMinLength && !HasContain(text, config.CollectorConfig.TitleExclude) && !HasPrefix(text, config.CollectorConfig.TitleExcludePrefix) && !HasSuffix(text, config.CollectorConfig.TitleExcludeSuffix) {
			title = text
		}
	}

	if title == "" {
		//获取title标签
		//title = doc.Find("#title,.title,.bt,.articleTit").First().Text()
		h2s := doc.Find("#title,.title,.bt,.articleTit,.right-xl>p,.biaoti")
		if h2s.Length() > 0 {
			for i := range h2s.Nodes {
				item := h2s.Eq(i)
				item.Children().Remove()
				text := strings.TrimSpace(item.Text())
				textLen := utf8.RuneCountInString(item.Text())
				if textLen >= config.CollectorConfig.TitleMinLength && textLen > utf8.RuneCountInString(title) && !HasContain(text, config.CollectorConfig.TitleExclude) && !HasPrefix(text, config.CollectorConfig.TitleExcludePrefix) && !HasSuffix(text, config.CollectorConfig.TitleExcludeSuffix) {
					title = text
				}
			}
		}
	}
	if title == "" {
		//如果标题为空，那么尝试h2
		h2s := doc.Find("h2,.name")
		if h2s.Length() > 0 {
			for i := range h2s.Nodes {
				item := h2s.Eq(i)
				item.Children().Remove()
				text := strings.TrimSpace(item.Text())
				textLen := utf8.RuneCountInString(text)
				if textLen >= config.CollectorConfig.TitleMinLength && textLen > utf8.RuneCountInString(title) && !HasContain(text, config.CollectorConfig.TitleExclude) && !HasPrefix(text, config.CollectorConfig.TitleExcludePrefix) && !HasSuffix(text, config.CollectorConfig.TitleExcludeSuffix) {
					title = text
				}
			}
		}
	}

	title = strings.Replace(strings.Replace(strings.TrimSpace(title), " ", "", -1), "\n", "", -1)
	title = strings.Replace(title, "<br>", "", -1)
	title = strings.Replace(title, "<br/>", "", -1)
	//只要第一个
	if utf8.RuneCountInString(title) > 50 {
		//减少误伤
		title = strings.ReplaceAll(title, "、", "-")
	}
	title = strings.ReplaceAll(title, "_", "-")
	sepIndex := strings.Index(title, "-")
	if sepIndex > 0 {
		title = title[:sepIndex]
	}

	article.Title = title
}


func(article *Article)ParseContent(doc *goquery.Document, body string) {
	content := ""
	contentText := ""
	description := ""
	contentLength := 0

	//对一些固定的内容，直接获取值
	contentItems := doc.Find("UCAPCONTENT,#mainText,.article-content,#article-content,#articleContnet,.entry-content,.the_body,.rich_media_content,#js_content,.word_content,.pages_content,.wendang_content,#content")
	if contentItems.Length() > 0 {
		for i := range contentItems.Nodes {
			contentItem := contentItems.Eq(i)
			content, _ = contentItem.Html()
			contentText = contentItem.Text()
			contentText = strings.Replace(contentText, " ", "", -1)
			contentText = strings.Replace(contentText, "\n", "", -1)
			contentText = strings.Replace(contentText, "\r", "", -1)
			contentText = strings.Replace(contentText, "\t", "", -1)
			nameRune := []rune(contentText)
			curLen := len(nameRune)
			if curLen > 150 {
				description = string(nameRune[:150])
			}
			//判断内容的真实性
			if curLen < config.CollectorConfig.ContentMinLength {
				contentText = ""
			}
			aCount := contentItem.Find("a").Length()
			if aCount > 5 {
				//太多连接了，直接放弃该内容
				contentText = ""
			}
			//查找内部div，如果存在，则使用它替代上一级
			divs := contentItem.Find("div")
			//只有内部没有div了或者内部div内容太少，才认为是真正的内容
			if divs.Length() > 0 {
				for i := range divs.Nodes {
					div := divs.Eq(i)
					if (div.Find("div").Length() == 0 || utf8.RuneCountInString(div.Find("div").Text()) < 100) && utf8.RuneCountInString(div.Text()) >= config.CollectorConfig.ContentMinLength {
						contentItem = div
						break
					}
				}
			}
			//排除一些不对的标签
			otherLength := contentItem.Find("input,textarea,form,button,footer,.footer").Length()
			if otherLength > 0 {
				contentText = ""
			}
			contentItem.Find("h1").Remove()
			//根据规则过滤
			if HasContain(contentText, config.CollectorConfig.ContentExclude) {
				contentText = ""
			}

			inner := contentItem.Find("*")
			for i := range inner.Nodes {
				item := inner.Eq(i)
				if HasContain(item.Text(), config.CollectorConfig.ContentExcludeLine) {
					item.Remove()
				}
			}

			if len(contentText) > 0 {
				break
			}
		}
	}

	if contentText == "" {
		content = ""
		//通用的获取方法
		divs := doc.Find("div")
		for i := range divs.Nodes {
			item := divs.Eq(i)
			pCount := item.ChildrenFiltered("p").Length()
			brCount := item.ChildrenFiltered("br").Length()
			aCount := item.Find("a").Length()
			if aCount > 5 {
				//太多连接了，直接放弃该内容
				continue
			}
			//排除一些不对的标签
			otherLength := item.Find("input,textarea,form,button,footer,.footer").Length()
			if otherLength > 0 {
				continue
			}
			if item.Find("div").Length() > 0 && utf8.RuneCountInString(item.Find("div").Text()) >= config.CollectorConfig.ContentMinLength {
				continue
			}
			if pCount > 0 || brCount > 0 {
				//表示查找到了一个p
				//移除空格和换行
				checkText := item.Text()
				checkText = strings.Replace(checkText, " ", "", -1)
				checkText = strings.Replace(checkText, "\n", "", -1)
				checkText = strings.Replace(checkText, "\r", "", -1)
				checkText = strings.Replace(checkText, "\t", "", -1)
				nameRune := []rune(checkText)
				curLen := len(nameRune)

				//根据规则过滤
				if HasContain(checkText, config.CollectorConfig.ContentExclude) {
					continue
				}
				if curLen <= config.CollectorConfig.ContentMinLength {
					continue
				}

				item.Find("h1,a").Remove()
				inner := item.Find("*")
				for i := range inner.Nodes {
					innerItem := inner.Eq(i)
					if HasContain(innerItem.Text(), config.CollectorConfig.ContentExcludeLine) {
						innerItem.Remove()
					}
				}

				if curLen > contentLength {
					contentLength = curLen
					content, _ = item.Html()
					contentText = checkText
					if curLen <= 150 {
						description = string(nameRune)
					} else {
						description = string(nameRune[:150])
					}
				}
			}
		}
	}
	//对内容进行处理
	re, _ := regexp.Compile("src=[\"']+?(.*?)[\"']+?[^>]+?>")
	content = re.ReplaceAllStringFunc(content, article.ReplaceSrc)

	re2, _ := regexp.Compile("href=[\"']+?(.*?)[\"']+?[^>]+?>")
	content = re2.ReplaceAllStringFunc(content, article.ReplaceHref)

	article.ContentText = contentText
	article.Description = strings.TrimSpace(description)
	article.Content = strings.TrimSpace(content)
}

func (article *Article) GetDomain() {
	baseUrlArr := strings.Split(article.OriginUrl, "/")
	pathUrlArr := baseUrlArr[:len(baseUrlArr)-1]
	baseUrlArr = baseUrlArr[:3]
	baseUrl := strings.Join(baseUrlArr, "/")
	article.OriginDomain = baseUrl
	article.OriginPath = strings.Join(pathUrlArr, "/")
}

func (article *Article) ReplaceSrc(src string) string {
	re, _ := regexp.Compile("src=[\"']+?(.*?)[\"']+?[^>]+?>")
	match := re.FindStringSubmatch(src)
	if len(match) < 1 {
		return src
	}

	if match[1] != "" {
		newSrc := ParseLink(match[1], article.OriginPath)
		src = strings.Replace(src, match[1], newSrc, -1)
	}
	return src
}

func (article *Article) ReplaceHref(src string) string {
	re, _ := regexp.Compile("href=[\"']+?(.*?)[\"']+?[^>]+?>")
	match := re.FindStringSubmatch(src)
	if len(match) < 1 {
		return src
	}


	if match[1] != "" {
		newSrc := ParseLink(match[1], article.OriginPath)
		src = strings.Replace(src, match[1], newSrc, -1)
	}
	return src
}

func InArray(need string, needArray []string) bool {
	for _, v := range needArray {
		if need == v {
			return true
		}
	}

	return false
}

func HasPrefix(need string, needArray []string) bool {
	for _, v := range needArray {
		if strings.HasPrefix(need, v) {
			return true
		}
	}

	return false
}

func HasSuffix(need string, needArray []string) bool {
	for _, v := range needArray {
		if strings.HasSuffix(need, v) {
			return true
		}
	}

	return false
}

func HasContain(need string, needArray []string) bool {
	for _, v := range needArray {
		if strings.Contains(need, v) {
			return true
		}
	}

	return false
}