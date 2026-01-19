package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"wxcloud-golang/db"
	"wxcloud-golang/handler"
	"wxcloud-golang/telemetry"

	"github.com/gin-gonic/gin"
	"github.com/tencentyun/cos-go-sdk-v5"
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
		api.POST("/upload", handler.UploadHandler)
	}
	ForceOpenPublicRead()
	log.Fatal(r.Run(":80"))
}

func envOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func ForceOpenPublicRead() {
	// 1. 填入你的桶 URL
	u, _ := url.Parse("https://7072-prod-0gr2o3qpe533f1fb-1352691102.cos.ap-shanghai.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}

	// 2. 填入你的 SecretId 和 SecretKey (必须是腾讯云的，不是微信的)
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("OS_TEMP_SECRET_ID"),  // 你的 SecretId
			SecretKey: os.Getenv("OS_TEMP_SECRET_KEY"), // 你的 SecretKey
		},
	})

	// 3. 强行设置为“公有读”
	opt := &cos.BucketPutACLOptions{
		Header: &cos.ACLHeaderOptions{
			XCosACL: "public-read", // 关键代码
		},
	}

	_, err := client.Bucket.PutACL(context.Background(), opt)
	if err != nil {
		fmt.Println("修改失败:", err)
	} else {
		fmt.Println("修改成功！现在是公有读了。")
	}
}
