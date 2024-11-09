package usecase

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Sanpeta/otel-and-zipkin-cep-temp-system/internal/entity"
)

func FetchCity(cep string) (string, error) {
	cepAPI := "https://viacep.com.br/ws/%s/json/"
	cepAPI = fmt.Sprintf(cepAPI, cep)

	resp, err := http.Get(cepAPI)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	var city entity.City
	err = json.Unmarshal(body, &city)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	return city.Localidade, nil
}
