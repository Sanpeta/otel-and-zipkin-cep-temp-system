package main

import (
	"net/http"

	"github.com/Sanpeta/otel-and-zipkin-cep-temp-system/internal/api"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /cep", api.HandlerServeA)

	http.ListenAndServe(":8080", mux)
}
