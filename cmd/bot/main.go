package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/yegor2025/budgetManager/internal/bot"
	"log"
)

const (
	token = "7578319650:AAFGX3s4UTjZYP5Ii8tdhPI3hCd0Av9BAVc"
)

func main() {
	tgBot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	tgBot.Debug = true

	log.Printf("Authorized on account %s", tgBot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := tgBot.GetUpdatesChan(u)

	handler := bot.NewHandler(tgBot)
	for update := range updates {
		go handler.HandlerUpdate(update)
	}
}
