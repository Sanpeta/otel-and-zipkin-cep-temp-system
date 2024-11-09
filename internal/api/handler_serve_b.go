package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/Sanpeta/otel-and-zipkin-cep-temp-system/internal/config"
	"github.com/Sanpeta/otel-and-zipkin-cep-temp-system/internal/entity"
	"github.com/Sanpeta/otel-and-zipkin-cep-temp-system/internal/usecase"
	"github.com/Sanpeta/otel-and-zipkin-cep-temp-system/pkg/utils"
	"go.opentelemetry.io/otel"
)

func HandlerServeB(w http.ResponseWriter, r *http.Request) {
	config, err := config.LoadConfig("../../")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	tracer := otel.Tracer("servico_b")

	var request entity.CEPRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil || !utils.CheckCEP(request.CEP) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	_, span := tracer.Start(r.Context(), "find_location_and_temp")
	defer span.End()

	city, err := usecase.FetchCity(request.CEP)
	if err != nil {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	fmt.Println("City:", city)

	encodedCity := url.QueryEscape(city)

	resp, err := usecase.FetchTemperature(encodedCity, config.TOKEN_WEATHER_API)
	if err != nil {
		http.Error(w, "error fetching temperature", http.StatusInternalServerError)
		return
	}

	response := entity.WeatherResponse{
		City:  city,
		TempC: resp.TempC,
		TempF: resp.TempC*1.8 + 32,
		TempK: resp.TempC + 273,
	}

	fmt.Println("Response:", response)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
