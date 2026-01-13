package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"wxcloud-golang/db"
	"wxcloud-golang/handler"

	"github.com/gin-gonic/gin"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	// 【修改】引入 HTTP 导出器
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

// ================= 配置区域 =================
// 这些默认值仅作本地开发兜底，生产环境请配置环境变量
const (
	DefaultServiceName = "plants"
	DefaultEndpoint    = "pl.ap-shanghai.apm.tencentcs.com:4319" // 你的 HTTP 端口
	DefaultToken       = "unZyTAhCEYVnKsKBSxGu"
)

// ===========================================

func CustomOtelMiddleware(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tracer := otel.Tracer(serviceName)
		ctx := otel.GetTextMapPropagator().Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

		opts := []trace.SpanStartOption{
			trace.WithAttributes(
				semconv.HTTPMethodKey.String(c.Request.Method),
				semconv.HTTPRouteKey.String(c.FullPath()),
				semconv.HTTPURLKey.String(c.Request.URL.String()),
			),
			trace.WithSpanKind(trace.SpanKindServer),
		}

		spanName := c.FullPath()
		if spanName == "" {
			spanName = fmt.Sprintf("HTTP %s route not found", c.Request.Method)
		}

		ctx, span := tracer.Start(ctx, spanName, opts...)
		defer span.End()

		c.Request = c.Request.WithContext(ctx)
		c.Next()

		status := c.Writer.Status()
		span.SetAttributes(semconv.HTTPStatusCodeKey.Int(status))
		if status >= 500 {
			span.SetStatus(codes.Error, fmt.Sprintf("HTTP %d", status))
		}
	}
}

func initTracer() (*sdktrace.TracerProvider, error) {
	ctx := context.Background()

	// 优先读取环境变量
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = DefaultEndpoint
	}

	token := os.Getenv("OTEL_TOKEN")
	if token == "" {
		token = DefaultToken
	}

	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = DefaultServiceName
	}

	// 【修改】配置 HTTP 导出器
	opts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(endpoint),
		// 因为是 HTTP (4319) 而不是 HTTPS (4320)，所以必须开启 Insecure
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithHeaders(map[string]string{
			"Authentication": token,
		}),
		otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
	}

	// 【修改】使用 otlptracehttp.New
	traceExporter, err := otlptracehttp.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			attribute.String("environment", "production"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp, nil
}

func main() {
	if err := db.Init(); err != nil {
		log.Printf("mysql init failed with %+v", err)
	}

	tp, err := initTracer()
	if err != nil {
		log.Printf("Warning: OpenTelemetry init failed: %v", err)
	} else {
		log.Println("✅ OpenTelemetry initialized (HTTP Mode)")
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := tp.Shutdown(ctx); err != nil {
				log.Printf("Error shutting down tracer provider: %v", err)
			}
		}()
	}

	r := gin.Default()

	svcName := os.Getenv("OTEL_SERVICE_NAME")
	if svcName == "" {
		svcName = DefaultServiceName
	}
	r.Use(CustomOtelMiddleware(svcName))

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
