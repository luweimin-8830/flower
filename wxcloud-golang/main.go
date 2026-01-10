package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"wxcloud-golang/db"
	"wxcloud-golang/handler"

	"github.com/gin-gonic/gin"

	// 引入 OpenTelemetry 相关包
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"
)

// 初始化 Tracer 的函数
func initTracer() (*sdktrace.TracerProvider, error) {
	ctx := context.Background()

	// 1. 获取上报地址 (从环境变量获取，不要写死)
	// 腾讯云 APM 通常要求 gRPC 端口 (4317)
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "pl.ap-shanghai.apm.tencentcs.com:4319" // 示例默认值，请按实际修改
	}

	// 2. 创建 gRPC 导出器
	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(), // 如果是内网或测试，通常用非加密；生产环境请确认文档要求
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithDialOption(grpc.WithBlock()), // 等待连接成功
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	// 3. 定义资源 (服务名等)
	res, err := resource.New(ctx,
		resource.WithAttributes(
			// 必填：你的服务名称，在 APM 面板上显示的名字
			semconv.ServiceNameKey.String("wxcloud-plant-service"),
			// 选填：其他标识
			attribute.String("environment", "production"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// 4. 创建 TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // 全量采集，生产环境可调整为 ParentBased(TraceIdRatioBased)
	)

	// 5. 设置全局 Tracer
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp, nil
}

func main() {
	if err := db.Init(); err != nil {
		panic(fmt.Sprintf("mysql init failed with %+v", err))
	}

	// 1. 初始化 Tracer
	tp, err := initTracer()
	if err != nil {
		// 如果 Tracer 初始化失败，不要 panic，打印错误即可，以免影响主业务
		log.Printf("Warning: OpenTelemetry init failed: %v", err)
	} else {
		// 确保程序退出时把剩下的数据发完
		defer func() {
			if err := tp.Shutdown(context.Background()); err != nil {
				log.Printf("Error shutting down tracer provider: %v", err)
			}
		}()
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
