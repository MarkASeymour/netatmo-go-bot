package handler

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func Main() {
	dg, err := discordgo.New("Bot " + getTokenFromConfig())
	if err != nil {
		fmt.Println("Error creating discord session: ", err)
		return
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection: ", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if strings.EqualFold(m.Content, "weather") {
		err, _ := s.ChannelMessageSend(m.ChannelID, WeatherPrint())
		if err != nil {
			fmt.Println("error sending message: ", err)
		}
	}
}

func getTokenFromConfig() string {
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

	token := viper.GetString("discord.token")
	return token
}
