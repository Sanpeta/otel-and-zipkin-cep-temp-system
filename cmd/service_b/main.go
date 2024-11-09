package main

import (
	"net/http"

	"github.com/Sanpeta/otel-and-zipkin-cep-temp-system/internal/api"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /buscar", api.HandlerServeB)

	http.ListenAndServe(":8081", mux)
}
