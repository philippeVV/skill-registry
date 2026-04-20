package telemetry

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

// Init sets up the OTEL TracerProvider.
//
//	endpoint="" + verbose=true  -> console exporter to stderr
//	endpoint="" + verbose=false -> no-op (silent)
//	endpoint set                -> OTLP HTTP exporter
func Init(ctx context.Context, endpoint string, verbose bool) (shutdown func(context.Context) error, err error) {
	noop := func(context.Context) error { return nil }

	var exporter sdktrace.SpanExporter

	switch {
	case endpoint != "":
		exporter, err = otlptracehttp.New(ctx,
			otlptracehttp.WithEndpoint(endpoint),
		)
		if err != nil {
			log.Debug().Err(err).Msg("failed to create OTLP exporter")
			return noop, nil
		}
	case verbose:
		exporter, err = stdouttrace.New(
			stdouttrace.WithWriter(os.Stderr),
			stdouttrace.WithPrettyPrint(),
		)
		if err != nil {
			log.Debug().Err(err).Msg("failed to create console exporter")
			return noop, nil
		}
	default:
		return noop, nil
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("skr"),
		),
	)
	if err != nil {
		log.Debug().Err(err).Msg("failed to create OTEL resource")
		return noop, nil
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter,
			sdktrace.WithBatchTimeout(1*time.Second),
		),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	tracer = tp.Tracer("skr")

	return func(ctx context.Context) error {
		return tp.Shutdown(ctx)
	}, nil
}

// EmitInstall emits an skr.package.install event.
func EmitInstall(ctx context.Context, name, version, pkgType, registryURL string) {
	if tracer == nil {
		return
	}
	_, span := tracer.Start(ctx, "skr.package.install",
		trace.WithAttributes(
			attribute.String("package.name", name),
			attribute.String("package.version", version),
			attribute.String("package.type", pkgType),
			attribute.String("registry.url", registryURL),
		),
	)
	span.End()
}

// EmitUninstall emits an skr.package.uninstall event.
func EmitUninstall(ctx context.Context, name, version string) {
	if tracer == nil {
		return
	}
	_, span := tracer.Start(ctx, "skr.package.uninstall",
		trace.WithAttributes(
			attribute.String("package.name", name),
			attribute.String("package.version", version),
		),
	)
	span.End()
}

// EmitUpdate emits an skr.package.update event.
func EmitUpdate(ctx context.Context, name, versionFrom, versionTo string) {
	if tracer == nil {
		return
	}
	_, span := tracer.Start(ctx, "skr.package.update",
		trace.WithAttributes(
			attribute.String("package.name", name),
			attribute.String("package.version_from", versionFrom),
			attribute.String("package.version_to", versionTo),
		),
	)
	span.End()
}
