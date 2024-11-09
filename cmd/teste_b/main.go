package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Sanpeta/otel-and-zipkin-cep-temp-system/internal/config"
	"github.com/Sanpeta/otel-and-zipkin-cep-temp-system/internal/usecase"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

const (
	viaCepURL        = "https://viacep.com.br/ws/%s/json/"
	weatherAPIURL    = "https://api.weatherapi.com/v1/current.json?key=YOUR_API_KEY&q=%s"
	otelCollectorURL = "localhost:4317" // URL do OTEL Collector
)

var tracer trace.Tracer

type CepRequest struct {
	CEP string `json:"cep"`
}

type WeatherResponse struct {
	Location struct {
		Name string `json:"name"`
	} `json:"location"`
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

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
			semconv.ServiceNameKey.String("ServiceB"),
		)),
	)

	otel.SetTracerProvider(tp)
	tracer = tp.Tracer("ServiceB")
}

func handler(w http.ResponseWriter, r *http.Request) {
	config, err := config.LoadConfig("../../")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	var cepReq CepRequest
	_, span := tracer.Start(r.Context(), "get-city-weather")
	defer span.End()

	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &cepReq)

	city, err := usecase.FetchCity(cepReq.CEP)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "can not find zipcode"}`))
		return
	}

	resp, err := usecase.FetchTemperature(city, config.TOKEN_WEATHER_API)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "weather service unavailable"}`))
		return
	}

	result := resp
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func main() {
	initTracer()
	mux := http.NewServeMux()

	mux.HandleFunc("POST /cep", handler)

	log.Fatal(http.ListenAndServe(":8081", mux))
}
