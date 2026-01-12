package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"wxcloud-golang/db"
	"wxcloud-golang/handler"

	"github.com/gin-gonic/gin"

	// OpenTelemetry 依赖
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"
)

// ================= 配置区域 (请修改这里) =================
const (
	// 1. 你的服务名称 (在 APM 列表里显示的名字)
	ServiceName = "plants"

	// 2. 接入点 (Endpoint)
	// 去腾讯云 APM 控制台 -> 应用监控 -> 接入设置 -> 接入点 复制 gRPC 的地址
	// 注意：不要带 http:// 前缀，例如 "ap-guangzhou.apm.tencentcs.com:4317"
	APMEndpoint = "pl.ap-shanghai.apm.tencentcs.com:4319" // TODO: 替换为你的真实接入点

	// 3. 鉴权 Token
	// 去腾讯云 APM 控制台 -> 应用监控 -> 接入设置 -> Token 复制
	APMToken = "unZyTAhCEYVnKsKBSxGu" // TODO: 替换为你的真实 Token
)

// ======================================================

func initTracer() (*sdktrace.TracerProvider, error) {
	ctx := context.Background()

	// 配置 gRPC 选项
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithInsecure(), // 使用非加密连接 (腾讯云 APM gRPC 通常支持非加密)
		otlptracegrpc.WithEndpoint(APMEndpoint),
		otlptracegrpc.WithHeaders(map[string]string{
			"Authentication": APMToken, // 鉴权头部
		}),
		otlptracegrpc.WithDialOption(grpc.WithBlock()), // 等待连接成功
		otlptracegrpc.WithTimeout(5 * time.Second),     // 连接超时时间
	}

	// 创建 Exporter
	exporter, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("创建 Exporter 失败: %w", err)
	}

	// 创建资源属性
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(ServiceName),
			attribute.String("deploy.environment", "production"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("创建 Resource 失败: %w", err)
	}

	// 创建 TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // 全量采集 (生产环境建议调整采样率)
	)

	// 设置全局 Tracer
	otel.SetTracerProvider(tp)
	// 设置上下文传播 (用于微服务链路透传)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp, nil
}

func main() {
	// 1. 初始化数据库
	if err := db.Init(); err != nil {
		// 数据库连不上是致命错误，直接 panic
		panic(fmt.Sprintf("mysql init failed: %v", err))
	}

	// 2. 初始化 OpenTelemetry
	tp, err := initTracer()
	if err != nil {
		// APM 初始化失败不应该影响主业务，打印日志即可
		log.Printf("⚠️  OpenTelemetry 初始化失败 (监控将不可用): %v", err)
	} else {
		log.Println("✅ OpenTelemetry 初始化成功")
		// 程序退出时关闭 Tracer
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := tp.Shutdown(ctx); err != nil {
				log.Printf("Tracer Shutdown error: %v", err)
			}
		}()
	}

	// 3. 初始化 Gin
	r := gin.Default()

	// 4. 添加 OpenTelemetry 中间件 (关键步骤)
	// 这行代码会自动记录所有 HTTP 请求的耗时和链路
	r.Use(otelgin.Middleware(ServiceName))

	// 5. 注册路由
	r.GET("/", handler.IndexHandler)

	api := r.Group("/api")
	{
		api.POST("/login", handler.UserLoginHandler)

		plant := api.Group("/plant")
		{
			plant.POST("/detail", handler.GetPlantHandler) // 建议改为 detail
			plant.POST("/list", handler.GetPlantsHandler)
			plant.POST("/add", handler.CreatePlantHandler)
			plant.POST("/delete", handler.DeletePlantHandler)
			plant.POST("/update", handler.UpdatePlantHandler)
		}

		tag := api.Group("/tag")
		{
			tag.POST("/add", handler.CreateTagHandler)
			tag.POST("/delete", handler.DeleteTagHandler)
			tag.POST("/update", handler.UpdateTagHandler)
			tag.POST("/list", handler.GetTagListHandler)
		}
	}

	// 6. 启动服务
	log.Println("🚀 Server starting on :80")
	log.Fatal(r.Run(":80"))
}
