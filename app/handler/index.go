package handler

import (
	"collector"
	"collector/config"
	"collector/services"
	"github.com/kataras/iris/v12"
	"log"
	"time"
)

func Index(ctx iris.Context) {
	ctx.View("index.html")
}

func IndexEchartsApi(ctx iris.Context) {
	//今天数据
	hours := []string{
		"01:00","02:00","03:00","04:00","05:00","06:00","07:00","08:00","09:00","10:00","11:00","12:00","13:00","14:00","15:00","16:00","17:00","18:00","19:00","20:00","21:00","22:00","23:00","23:59",
	}
	loc, _ := time.LoadLocation("Local")
	nowTime := time.Now()
	days := []string{
		nowTime.Format("01-02"),
		nowTime.AddDate(0, 0, -1).Format("01-02"),
		nowTime.AddDate(0, 0, -2).Format("01-02"),
		nowTime.AddDate(0, 0, -3).Format("01-02"),
		nowTime.AddDate(0, 0, -4).Format("01-02"),
		nowTime.AddDate(0, 0, -5).Format("01-02"),
		nowTime.AddDate(0, 0, -6).Format("01-02"),
	}
	var articleHourCounts []int
	var articleDayCounts []int
	var sourceHourCounts []int
	var sourceDayCounts []int
	for _, v := range hours {
		var articleCount int
		var sourceCount int
		log.Println(time.Now().Format("2006-01-02") + " " + v)
		countTime, _ := time.ParseInLocation("2006-01-02 15:04", time.Now().Format("2006-01-02") + " " + v, loc)
		endTime := countTime.Unix()
		startTime := endTime - 3600
		services.DB.Model(&collector.Article{}).Where("`status` = 1").Where("`created_time` >= ?", startTime).Where("`created_time` < ?", endTime).Count(&articleCount)
		articleHourCounts = append(articleHourCounts, articleCount)

		services.DB.Model(&collector.Article{}).Where("`status` = 1").Where("`created_time` >= ?", startTime).Where("`created_time` < ?", endTime).Group("source_id").Count(&sourceCount)
		sourceHourCounts = append(sourceHourCounts, sourceCount)
	}
	for i, _ := range days {
		var articleCount int
		var sourceCount int

		countDay, _ := time.ParseInLocation("2006-01-02 15:04", nowTime.AddDate(0, 0, -i).Format("2006-01-02 00:00"), loc)
		startTime := countDay.Unix()
		endTime := startTime +86400
		services.DB.Model(&collector.Article{}).Where("`status` = 1").Where("`created_time` >= ?", startTime).Where("`created_time` < ?", endTime).Count(&articleCount)
		articleDayCounts = append(articleDayCounts, articleCount)

		services.DB.Model(&collector.Article{}).Where("`status` = 1").Where("`created_time` >= ?", startTime).Where("`created_time` < ?", endTime).Group("source_id").Count(&sourceCount)

		sourceDayCounts = append(sourceDayCounts, sourceCount)
	}

	ctx.JSON(iris.Map{
		"code": config.StatusOK,
		"msg":  "",
		"data": iris.Map{
			"articleHourCounts": articleHourCounts,
			"articleDayCounts": articleDayCounts,
			"sourceHourCounts": sourceHourCounts,
			"sourceDayCounts": sourceDayCounts,
			"hours": hours,
			"days": days,
		},
	})
}