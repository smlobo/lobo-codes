package internal

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitTracerProvider(serviceName string) (*trace.TracerProvider) {
	// The B3 HTTP header propagation
	b3Propagator := b3.New()
	otel.SetTextMapPropagator(b3Propagator)

	// Desired exporter
	var exp trace.SpanExporter
	var err error
	switch config["EXPORTER"] {
	case "jaeger":
		exp, err = jaegerExporter()
	case "file":
		file, err := os.Create(serviceName + "-traces.txt")
		if err != nil {
			log.Fatal(err)
		}
		exp, err = fileExporter(file)
		// Can't close file after this function. Don't bother closing.
		//defer file.Close()
	case "zipkin":
		exp, err = zipkinExporter()
	case "otlp":
		exp, err = otlpExporter()
	case "otlp-grpc":
		exp, err = otlpGrpcExporter()
	}
	if err != nil {
		log.Fatal(err)
	}

	// Set the service name - and any other attributes.
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		log.Fatal("Could not set resources: ", err)
	}

	// Set the main batched tracer provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(resources),
	)
	otel.SetTracerProvider(tp)

	return tp
}

// Console exporter.
func fileExporter(w io.Writer) (trace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		stdouttrace.WithPrettyPrint(),
		stdouttrace.WithoutTimestamps(),
	)
}

// Jaeger exporter
func jaegerExporter() (trace.SpanExporter, error) {
	jaegerServer := config["JAEGER_SERVER"]
	jaegerPort := config["JAEGER_PORT"]
	jaegerPath := config["JAEGER_PATH"]
	jaegerEndpoint := "http://" + jaegerServer + ":" + jaegerPort + jaegerPath
	return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)))
}

// Zipkin exporter
func zipkinExporter() (trace.SpanExporter, error) {
	return zipkin.New(config["ZIPKIN_ENDPOINT"])
}

// OTLP exporter
func otlpExporter() (trace.SpanExporter, error) {
	log.Printf("Setting up OTLP HTTP exporter: %s:%s %s\n", config["OTLP_SERVER"], config["OTLP_HTTP_PORT"],
		config["OTLP_URL"])
	client := otlptracehttp.NewClient(otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint(config["OTLP_SERVER"] + ":" + config["OTLP_HTTP_PORT"]),
		otlptracehttp.WithURLPath(config["OTLP_URL"]))
	return otlptrace.New(context.Background(), client)
}

// OTLP grpc exporter
func otlpGrpcExporter() (trace.SpanExporter, error) {
	log.Printf("Setting up OTLP GRPC exporter: %s:%s\n", config["OTLP_SERVER"], config["OTLP_GRPC_PORT"])
	conn, err := grpc.Dial(config["OTLP_SERVER"] + ":" + config["OTLP_GRPC_PORT"],
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}
	client := otlptracegrpc.NewClient(otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithGRPCConn(conn))
	return otlptrace.New(context.Background(), client)
}
