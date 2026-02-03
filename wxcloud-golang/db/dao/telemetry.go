package dao

import (
	"context"
	"fmt"

	"wxcloud-golang/db"
	"wxcloud-golang/telemetry"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

func execWithSpan(ctx context.Context, operation, table string, fn func(*gorm.DB) error) error {
	if ctx == nil {
		ctx = context.Background()
	}

	tracer := telemetry.Tracer()
	ctx, span := tracer.Start(ctx, fmt.Sprintf("%s.%s", table, operation), trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	span.SetAttributes(
		semconv.DBSystemMySQL,
		attribute.String("db.operation", operation),
		attribute.String("db.sql.table", table),
	)

	conn := db.Conn(ctx)
	if conn == nil {
		err := fmt.Errorf("gorm DB is not initialized")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	if err := fn(conn); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

// GetDB 返回绑定了 Context 的数据库连接
func GetDB(ctx context.Context) *gorm.DB {
	return db.Conn(ctx)
}
