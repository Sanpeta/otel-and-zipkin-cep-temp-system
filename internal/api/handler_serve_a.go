package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/Sanpeta/otel-and-zipkin-cep-temp-system/internal/entity"
	"github.com/Sanpeta/otel-and-zipkin-cep-temp-system/pkg/utils"
	"go.opentelemetry.io/otel"
)

func HandlerServeA(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("servico_a")

	var request entity.CEPRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil || !utils.CheckCEP(request.CEP) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	_, span := tracer.Start(r.Context(), "send_cep_to_service_b")
	defer span.End()

	jsonBody, err := json.Marshal(request)
	if err != nil {
		http.Error(w, "error parsing request", http.StatusInternalServerError)
		return
	}

	// Enviar para o Serviço B
	resp, err := http.Post("http://localhost:8081/buscar", "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		http.Error(w, "error calling service B", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, `{"error": "error reading response from service B"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}
