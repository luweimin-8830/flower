package main

import (
	"context"
	"log"
	"os"
	"time"

	"wxcloud-golang/db"
	"wxcloud-golang/handler"
	"wxcloud-golang/telemetry"

	"github.com/gin-gonic/gin"
)

const (
	defaultServiceName = "plants"
	defaultEndpoint    = "pl.ap-shanghai.apm.tencentcs.com:4319"
	defaultToken       = "unZyTAhCEYVnKsKBSxGu"
	defaultEnvironment = "production"
)

func main() {
	cfg := telemetry.Config{
		ServiceName: envOrDefault("OTEL_SERVICE_NAME", defaultServiceName),
		Endpoint:    envOrDefault("OTEL_EXPORTER_OTLP_ENDPOINT", defaultEndpoint),
		Token:       envOrDefault("OTEL_TOKEN", defaultToken),
		Environment: envOrDefault("OTEL_ENVIRONMENT", defaultEnvironment),
	}

	tp, err := telemetry.Init(context.Background(), cfg)
	if err != nil {
		log.Printf("Warning: OpenTelemetry init failed: %v", err)
	} else {
		log.Println("✅ OpenTelemetry initialized via telemetry.Init")
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := telemetry.Shutdown(ctx, tp); err != nil {
				log.Printf("Error shutting down tracer provider: %v", err)
			}
		}()
	}

	if err := db.Init(); err != nil {
		log.Fatalf("mysql init failed: %+v", err)
	}

	r := gin.Default()
	telemetry.UseMiddleware(r, cfg.ServiceName)

	r.GET("/", handler.IndexHandler)

	api := r.Group("/api")
	{
		api.POST("/login", handler.UserLoginHandler)
		family := api.Group("/family")
		{
			family.POST("/", handler.GetFamilyHandler)
			family.POST("/sort", handler.SortFamilyHandler)
			family.POST("/switch", handler.SwitchFamilyHandler)
			family.POST("/delete", handler.DeleteFamilyHandler)
			family.POST("/update", handler.UpdateFamilyHandler)
		}

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
			tag.POST("/sort", handler.SortTagHandler)
		}
		image := api.Group("/image")
		{
			image.POST("/check", handler.CheckImageHandler)
			image.POST("/add", handler.SaveImageHandler)
		}
	}
	log.Fatal(r.Run(":80"))
}

func envOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
