package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
	"wxcloud-golang/db"
	"wxcloud-golang/db/model"
	"wxcloud-golang/response"
)

type PriceChartRequest struct {
	FamilyID  uint   `json:"familyId" binding:"required"`
	StartDate string `json:"startDate" binding:"required"`
	EndDate   string `json:"endDate" binding:"required"`
}
type PriceChartResponse struct {
	TotalCents    int64       `json:"totalCents"`
	PricedCount   int64       `json:"pricedCount"`
	UnpricedCount int64       `json:"unpricedCount"`
	Buckets       []ChartItem `json:"buckets"`
}
type priceCount struct {
	PriceCents *int64
	Count      int64
}

func priceDateRange(start, end string) (time.Time, time.Time, error) {
	loc := time.FixedZone("UTC+8", 8*60*60)
	a, err := time.ParseInLocation("2006-01-02", start, loc)
	if err != nil {
		return a, a, err
	}
	b, err := time.ParseInLocation("2006-01-02", end, loc)
	if err != nil {
		return a, b, err
	}
	if a.After(b) {
		return a, b, fmt.Errorf("开始日期不能晚于结束日期")
	}
	return a, b.AddDate(0, 0, 1), nil
}

func summarizePrices(rows []priceCount) PriceChartResponse {
	result := PriceChartResponse{Buckets: []ChartItem{
		{Name: "0～50元"}, {Name: "50～100元"}, {Name: "100～200元"},
		{Name: "200～500元"}, {Name: "500～1000元"}, {Name: "1000元及以上"},
	}}
	bounds := []int64{5000, 10000, 20000, 50000, 100000}
	for _, row := range rows {
		if row.PriceCents == nil {
			result.UnpricedCount += row.Count
			continue
		}
		result.PricedCount += row.Count
		result.TotalCents += *row.PriceCents * row.Count
		bucket := 0
		for bucket < len(bounds) && *row.PriceCents >= bounds[bucket] {
			bucket++
		}
		result.Buckets[bucket].Value += row.Count
	}
	return result
}

func GetPriceChartHandler(c *gin.Context) {
	var req PriceChartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	start, end, err := priceDateRange(req.StartDate, req.EndDate)
	if err != nil {
		response.Fail(c, "请选择有效的起止日期")
		return
	}
	openID := c.GetHeader("X-WX-OPENID")
	if openID == "" {
		response.FailWithCode(c, 403, "无权查看该家庭统计")
		return
	}
	var members int64
	err = db.DB.WithContext(c.Request.Context()).Model(&model.FamilyMember{}).
		Where("family_id = ? AND open_id = ?", req.FamilyID, openID).Count(&members).Error
	if err != nil {
		response.Fail(c, "获取统计数据失败")
		return
	}
	if members == 0 {
		response.FailWithCode(c, 403, "无权查看该家庭统计")
		return
	}
	var rows []priceCount
	err = db.DB.WithContext(c.Request.Context()).Model(&model.Plant{}).
		Select("price_cents, COUNT(*) AS count").
		Where("family_id = ? AND birthday >= ? AND birthday < ?", req.FamilyID, start, end).
		Group("price_cents").Scan(&rows).Error
	if err != nil {
		response.Fail(c, "获取价格统计失败")
		return
	}
	response.Success(c, summarizePrices(rows))
}
