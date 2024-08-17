package main

import (
	"encoding/json"
	"math"
	"net/http"
	"os"
	"strings"
)

type apiConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Celsius float64 `json:"temp"`
	} `json:"main"`
}

func loadApiConfig(filename string) (apiConfigData, error) {
	bytes, err := os.ReadFile(filename)

	if err != nil {
		return apiConfigData{}, err
	}

	var data apiConfigData

	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return apiConfigData{}, err
	}
	return data, nil

}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello! I'm tempTracker!"))
}

func query(city string) (weatherData, error) {
	apiConfig, err := loadApiConfig((".apiConfig"))
	if err != nil {
		return weatherData{}, err
	}
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?q=" + city + "&APPID=" + apiConfig.OpenWeatherMapApiKey)
	if err != nil {
		return weatherData{}, err
	}
	defer resp.Body.Close()
	var data weatherData
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return weatherData{}, err
	}
	data.Main.Celsius = math.Round((data.Main.Celsius-273.15)*100) / 100
	return data, nil
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/weather/",
		func(w http.ResponseWriter, r *http.Request) {
			city := strings.SplitN(r.URL.Path, "/", 3)[2]
			data, err := query(city)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(data)
		})

	http.ListenAndServe(":8080", nil)
}
