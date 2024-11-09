package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Sanpeta/otel-and-zipkin-cep-temp-system/internal/entity"
	"github.com/Sanpeta/otel-and-zipkin-cep-temp-system/pkg/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

const (
	serviceBURL      = "http://localhost:8081/cep"
	otelCollectorURL = "localhost:4317" // URL do OTEL Collector
)

var tracer trace.Tracer

func initTracer() {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, otelCollectorURL, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to create gRPC connection to collector: %v", err)
	}

	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		log.Fatalf("failed to create OTLP trace exporter: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("ServiceA"),
		)),
	)

	otel.SetTracerProvider(tp)
	tracer = tp.Tracer("ServiceA")
}

func handler(w http.ResponseWriter, r *http.Request) {
	var cepReq entity.CEPRequest
	_, span := tracer.Start(r.Context(), "validate-cep")
	defer span.End()

	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &cepReq)

	if !utils.CheckCEP(cepReq.CEP) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`{"message": "invalid zipcode"}`))
		return
	}

	req, _ := json.Marshal(cepReq)
	resp, err := http.Post(serviceBURL, "application/json", bytes.NewBuffer(req))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "service unavailable"}`))
		return
	}
	defer resp.Body.Close()

	body, _ = ioutil.ReadAll(resp.Body)
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func main() {
	initTracer()

	mux := http.NewServeMux()

	mux.HandleFunc("POST /cep", handler)

	http.ListenAndServe(":8080", mux)
}
