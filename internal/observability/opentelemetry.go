package observability

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	otelExporterOTLPEndpoint = "OTEL_EXPORTER_OTLP_ENDPOINT"
	otelExporterProtocol     = "OTEL_EXPORTER_PROTOCOL"
	insecureMode             = "INSECURE_MODE"
)

// InitProvider Initializes an OTLP exporter, and configures the corresponding trace and
// metric providers.
func InitProvider(appName string) (trace.TracerProvider, error) {
	ctx := context.Background()

	otelAgentEndpoint, ok := os.LookupEnv(otelExporterOTLPEndpoint)
	if !ok {
		otelAgentEndpoint = "0.0.0.0:4317"
	}

	otelAgentProtocol, ok := os.LookupEnv(otelExporterProtocol)
	if !ok {
		otelAgentProtocol = "grpc"
	}

	insecure, ok := os.LookupEnv(insecureMode)
	if !ok {
		insecure = "true"
	}

	exporter, err := traceExporter(ctx, insecure, otelAgentEndpoint, otelAgentProtocol)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String(appName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter),
	)

	otel.SetTracerProvider(tracerProvider)
	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tracerProvider, nil
}

func traceExporter(ctx context.Context, insecure, otelAgentEndpoint, otelAgentProtocol string) (*otlptrace.Exporter, error) {
	var exporter *otlptrace.Exporter
	var err error

	switch strings.ToLower(otelAgentProtocol) {
	case "http", "https":
		exporter, err = httpExporter(ctx, insecure, otelAgentEndpoint)
		if err != nil {
			return nil, err
		}
	case "grpc":
		exporter, err = grpcExpoter(ctx, insecure, otelAgentEndpoint)
		if err != nil {
			return nil, err
		}
	}
	return exporter, nil
}

func grpcExpoter(ctx context.Context, insecure, otelAgentEndpoint string) (*otlptrace.Exporter, error) {

	var secureOption otlptracegrpc.Option

	if strings.ToLower(insecure) == "false" || insecure == "0" || strings.ToLower(insecure) == "f" {
		secureOption = otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	} else {
		secureOption = otlptracegrpc.WithInsecure()
	}

	traceClient := otlptracegrpc.NewClient(
		secureOption,
		otlptracegrpc.WithEndpoint(otelAgentEndpoint),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
	)

	sctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	traceExporter, err := otlptrace.New(sctx, traceClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	return traceExporter, nil
}

func httpExporter(ctx context.Context, insecure, otelAgentEndpoint string) (*otlptrace.Exporter, error) {
	var secureOption otlptracehttp.Option

	if strings.ToLower(insecure) == "false" || insecure == "0" || strings.ToLower(insecure) == "f" {
		//secureOption = otlptracehttp.WithTLSClientConfig()
	} else {
		secureOption = otlptracehttp.WithInsecure()
	}

	traceClientHttp := otlptracehttp.NewClient(
		secureOption,
		otlptracehttp.WithEndpoint(otelAgentEndpoint),
	)
	otlptracehttp.WithCompression(1)

	traceExporter, err := otlptrace.New(ctx, traceClientHttp)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	return traceExporter, err
}
