package telemetry

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

var tracerName = "wxcloud-golang"

// Config 用于初始化 OTel TracerProvider。
type Config struct {
	ServiceName string
	Environment string
	Endpoint    string
	Token       string
	Timeout     time.Duration
}

func (c *Config) normalize() {
	if c.ServiceName == "" {
		c.ServiceName = tracerName
	}
	if c.Environment == "" {
		c.Environment = "production"
	}
	if c.Endpoint == "" {
		c.Endpoint = "pl.ap-shanghai.apm.tencentcs.com:4320"
	}
	if c.Timeout <= 0 {
		c.Timeout = 5 * time.Second
	}
	tracerName = c.ServiceName
}

// Init 配置全局 TracerProvider，并返回 provider 以便调用方优雅关闭。
func Init(ctx context.Context, cfg Config) (*sdktrace.TracerProvider, error) {
	cfg.normalize()

	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(cfg.Endpoint),
	}

	exporter, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("telemetry exporter init failed: %w", err)
	}

	res, err := resource.New(ctx, []resource.Option{
		resource.WithAttributes(
			attribute.String("token", "unZyTAhCEYVnKsKBSxGu"),
			attribute.String("service.name", "plants"),
			attribute.String("host.name", "127.0.0.1"),
		),
	}...)
	if err != nil {
		return nil, fmt.Errorf("telemetry resource init failed: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp, nil
}

// Shutdown 优雅关闭 tracer provider。
func Shutdown(ctx context.Context, tp *sdktrace.TracerProvider) error {
	if tp == nil {
		return nil
	}
	return tp.Shutdown(ctx)
}

// Tracer 暴露统一的 tracer 供内部使用。
func Tracer() trace.Tracer {
	return otel.Tracer(tracerName)
}

// GinMiddleware 返回可直接注册到 gin.Use 的 OTel 中间件。
func GinMiddleware(serviceName string) gin.HandlerFunc {
	if serviceName == "" {
		serviceName = tracerName
	}

	return func(c *gin.Context) {
		tracer := otel.Tracer(serviceName)
		ctx := otel.GetTextMapPropagator().Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

		spanName := c.FullPath()
		if spanName == "" {
			spanName = fmt.Sprintf("HTTP %s route not found", c.Request.Method)
		}

		ctx, span := tracer.Start(ctx, spanName,
			trace.WithAttributes(
				semconv.HTTPMethodKey.String(c.Request.Method),
				semconv.HTTPRouteKey.String(c.FullPath()),
				semconv.HTTPURLKey.String(c.Request.URL.String()),
			),
			trace.WithSpanKind(trace.SpanKindServer),
		)
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

// UseMiddleware 采用 useMiddleware 的方式接入 Gin。
func UseMiddleware(router gin.IRoutes, serviceName string) {
	if router == nil {
		return
	}
	router.Use(GinMiddleware(serviceName))
}
