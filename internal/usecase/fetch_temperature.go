package usecase

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/Sanpeta/otel-and-zipkin-cep-temp-system/internal/entity"
	"github.com/Sanpeta/otel-and-zipkin-cep-temp-system/pkg/utils"
)

func FetchTemperature(city string, token string) (entity.TemperatureResponse, error) {
	encodedCity := url.QueryEscape(city)
	weatherAPI := "https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s"
	weatherAPI = fmt.Sprintf(weatherAPI, encodedCity, token)

	resp, err := http.Get(weatherAPI)
	if err != nil {
		fmt.Println("Error:", err)
		return entity.TemperatureResponse{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return entity.TemperatureResponse{}, err
	}

	var weather entity.Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		fmt.Println("Error:", err)
		return entity.TemperatureResponse{}, err
	}

	return entity.TemperatureResponse{
		TempC: utils.RoundToTwoDecimals(weather.Main.Temp - 273.15),
		TempF: utils.RoundToTwoDecimals((weather.Main.Temp-273.15)*1.8 + 32),
		TempK: utils.RoundToTwoDecimals(weather.Main.Temp),
	}, nil
}
