package handler

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
	"wxcloud-golang/db"
	"wxcloud-golang/db/model"
	"wxcloud-golang/response"

	"github.com/gin-gonic/gin"
)

// GetWeatherHandler 处理获取天气请求
func GetWeatherHandler(c *gin.Context) {
	amapKey := os.Getenv("A_MAP_KEY")
	if amapKey == "" {
		// 如果环境变量为空，可以设置一个默认值或返回错误
		amapKey = "a98bb85a364792f820ae0703f76852f2"
	}

	var req WeatherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	openID := c.GetHeader("X-WX-OPENID")
	ctx := c.Request.Context()
	database := db.Conn(ctx)

	var user model.User
	var city string
	var adcode string

	// 1. 前端传入经纬度,先去和user表存入对比,如果3位小数内都一致,直接取出市名
	if openID != "" {
		if err := database.Where("open_id = ?", openID).First(&user).Error; err == nil {
			// 对比经纬度，保留3位小数比较
			if user.City != "" && user.Adcode != "" &&
				math.Abs(req.Longitude-user.Longitude) < 0.001 &&
				math.Abs(req.Latitude-user.Latitude) < 0.001 {
				city = user.City
				adcode = user.Adcode
			}
		}
	}

	// 2. 如果不一致或没有值,则请求amap获取市名并存入user表
	if city == "" {
		regeoUrl := fmt.Sprintf("https://restapi.amap.com/v3/geocode/regeo?location=%.6f,%.6f&key=%s", req.Longitude, req.Latitude, amapKey)
		resp, err := http.Get(regeoUrl)
		if err == nil {
			defer resp.Body.Close()
			var regeoData AmapRegeoResponse
			if err := json.NewDecoder(resp.Body).Decode(&regeoData); err == nil && regeoData.Status == "1" {
				addressComp := regeoData.Regeocode.AddressComponent
				var rawCity string
				switch v := addressComp.City.(type) {
				case string:
					if len(v) > 0 {
						rawCity = v
					} else {
						rawCity = addressComp.Province
					}
				default:
					rawCity = addressComp.Province
				}
				city = strings.TrimSuffix(rawCity, "市")
				adcode = addressComp.Adcode

				// 存入 user 表
				if openID != "" {
					database.Model(&model.User{}).Where("open_id = ?", openID).Updates(map[string]interface{}{
						"longitude": req.Longitude,
						"latitude":  req.Latitude,
						"city":      city,
						"adcode":    adcode,
					})
				}
			}
		}
	}

	if city == "" {
		response.Fail(c, "获取地理位置失败")
		return
	}

	// 3. 先用市名和今日日期去天气表中查询
	today := time.Now().Format("2006-01-02")
	var cache model.WeatherCache
	if err := database.Where("city = ? AND date = ?", city, today).First(&cache).Error; err == nil {
		response.Success(c, gin.H{
			"city":    cache.City,
			"adcode":  cache.Adcode,
			"temp":    cache.Temperature,
			"icon":    cache.Icon,
			"weather": cache.Weather,
		})
		return
	}

	// 4. 如没有,则调用天气接口查询后存入天气表
	weatherUrl := fmt.Sprintf("https://restapi.amap.com/v3/weather/weatherInfo?city=%s&key=%s&extensions=all", adcode, amapKey)
	wResp, err := http.Get(weatherUrl)
	if err != nil {
		response.Fail(c, "获取天气接口失败")
		return
	}
	defer wResp.Body.Close()

	var weatherData AmapWeatherResponse
	if err := json.NewDecoder(wResp.Body).Decode(&weatherData); err != nil || weatherData.Status != "1" || len(weatherData.Lives) == 0 {
		response.Fail(c, "解析天气数据失败")
		return
	}

	live := weatherData.Lives[0]
	icon := mapWeatherToIcon(live.Weather)

	// 存入天气缓存表
	newCache := model.WeatherCache{
		City:        city,
		Adcode:      adcode,
		Date:        today,
		Temperature: live.Temperature,
		Weather:     live.Weather,
		Icon:        icon,
	}
	database.Create(&newCache)

	response.Success(c, gin.H{
		"city":    city,
		"adcode":  adcode,
		"temp":    live.Temperature,
		"icon":    icon,
		"weather": live.Weather,
	})
}

type WeatherRequest struct {
	Longitude float64 `json:"longitude" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required"`
}

type AmapRegeoResponse struct {
	Status    string `json:"status"`
	Regeocode struct {
		AddressComponent struct {
			City     interface{} `json:"city"`
			Province string      `json:"province"`
			Adcode   string      `json:"adcode"`
		} `json:"addressComponent"`
	} `json:"regeocode"`
}

type AmapWeatherResponse struct {
	Status string `json:"status"`
	Lives  []struct {
		Province    string `json:"province"`
		City        string `json:"city"`
		Adcode      string `json:"adcode"`
		Weather     string `json:"weather"`
		Temperature string `json:"temperature"`
	} `json:"lives"`
}

func mapWeatherToIcon(weather string) string {
	if strings.Contains(weather, "晴") {
		return "sun-filled"
	}
	if strings.Contains(weather, "云") || strings.Contains(weather, "阴") {
		return "cloud-filled"
	}
	if strings.Contains(weather, "雨") {
		return "rain-filled"
	}
	if strings.Contains(weather, "雪") {
		return "snow-filled"
	}
	return "location-filled"
}
