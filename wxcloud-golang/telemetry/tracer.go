package telemetry

import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
)

const tracerName = "wxcloud-golang"

func Tracer() trace.Tracer {
    return otel.Tracer(tracerName)
}
