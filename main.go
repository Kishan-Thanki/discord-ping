package main

import (
	"fmt"

	"github.com/Kishan-Thanki/discord-ping/bot"
	"github.com/Kishan-Thanki/discord-ping/config"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found. Ensure environment variables are set.")
	}

	err = config.ReadConfig()
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}

	bot.Start()

	<-make(chan struct{})
}
