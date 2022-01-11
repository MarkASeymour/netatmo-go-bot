package handler_test

import (
	"testing"

	"github.com/markaseymour/netatmo-go-bot/handler"
)

func TestWeatherPrint(t *testing.T) {
	output, err := handler.WeatherPrint()
	if err != nil {
		t.Error("returned error", err)
	}
	if output == "" {
		t.Error("returns empty string")
	}
}
