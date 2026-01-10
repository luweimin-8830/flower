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
		plant := api.Group("/plant")
		{
			plant.POST("/", handler.GetPlantHandler)
			plant.POST("/list", handler.GetPlantsHandler)
			plant.POST("/add", handler.CreatePlantHandler)
			plant.POST("/delete", handler.DeletePlantHandler)
			plant.POST("/update", handler.UpdatePlantHandler)
		}
		tag := api.Group("/tag")
		{
			tag.POST("/", handler.GetTagListHandler)
			tag.POST("/add", handler.CreateTagHandler)
			tag.POST("/delete", handler.DeleteTagHandler)
			tag.POST("/update", handler.UpdateTagHandler)
		}
	}

	log.Fatal(r.Run(":80"))
}
