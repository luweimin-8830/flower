package main

import (
	"fmt"
	"log"

	"wxcloud-golang/db"
	"wxcloud-golang/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := db.Init(); err != nil {
		panic(fmt.Sprintf("mysql init failed with %+v", err))
	}

	// 2. 初始化 Gin 引擎
	r := gin.Default()

	r.GET("/", handler.IndexHandler)

	api := r.Group("/api")
	{
		api.POST("/login", handler.UserLoginHandler)
		api.POST("/plant", handler.GetPlantHandler)
		api.POST("/plant/list", handler.GetPlantsHandler)
		api.POST("/plant/add", handler.CreatePlantHandler)
		api.POST("/plant/delete", handler.DeletePlantHandler)
		api.POST("/plant/update", handler.UpdatePlantHandler)
		api.POST("/tag/add", handler.CreateTagHandler)
		api.POST("/tag/delete", handler.DeleteTagHandler)
		api.POST("/tag/update", handler.UpdateTagHandler)
		api.POST("/tag", handler.GetTagListHandler)
	}

	log.Fatal(r.Run(":80"))
}
