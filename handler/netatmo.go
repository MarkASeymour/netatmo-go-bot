package handler

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	netatmo "github.com/exzz/netatmo-api-go"
	"github.com/spf13/viper"
)

func WeatherPrint() string {

	viper.SetConfigName("appconfig")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./appconfig/")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Error reading config \n", err)
		os.Exit(1)
	}

	n, err := netatmo.NewClient(netatmo.Config{
		ClientID:     viper.GetString("netatmo.clientID"),
		ClientSecret: viper.GetString("netatmo.clientSecret"),
		Username:     viper.GetString("netatmo.username"),
		Password:     viper.GetString("netatmo.password"),
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dc, err := n.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ct := time.Now().UTC().Unix()

	var message string = ""

	for _, station := range dc.Stations() {
		message += fmt.Sprintf("Station: %s\n", "My Home - Westerville")

		for _, module := range station.Modules() {
			if module.ModuleName != "" {
				message += fmt.Sprintf("\tModule : %s\n", module.ModuleName)
			} else {
				message += fmt.Sprintf("\tModule : %s\n", "Indoor")
			}

			{
				if module.DashboardData.LastMeasure == nil {
					fmt.Printf("\t\tSkipping %s, no measurement data available.\n", module.ModuleName)
					continue
				}
				ts, data := module.Info()
				for dataName, value := range data {
					message += fmt.Sprintf("\t\t%s : %v (updated %ds ago)\n", dataName, value, ct-ts)

				}
			}

			{
				ts, data := module.Data()
				for dataName, value := range data {
					if strings.EqualFold(dataName, "Temperature") {
						fmt.Println("dropped into if-block")
						v := fmt.Sprintf("%v", value)
						newVal, err := strconv.ParseFloat(v, 64)
						if err != nil {
							fmt.Println("error parsing temp", err)
						}
						temp := (newVal * 1.8) + 32

						message += fmt.Sprintf("\t\t%s : %v °F (updated %ds ago)\n", dataName, math.Round(temp), ct-ts)
					} else {
						message += fmt.Sprintf("\t\t%s : %v (updated %ds ago)\n", dataName, value, ct-ts)
					}

				}
			}
		}
	}
	return message
}
