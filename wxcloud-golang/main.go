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
		api.POST("/plantList", handler.PlantListHandler)
		api.POST("/plant/add", handler.AddPlantHandler)
		api.POST("/plant/delete", handler.DeletePlantHandler)
		api.POST("/plant/update", handler.UpdatePlantHandler)
	}

	log.Fatal(r.Run(":80"))
}
