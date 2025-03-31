package main

import (
	"context"
	"github.com/yegor2025/budgetManager/cilents/telegram"
	event_consumer "github.com/yegor2025/budgetManager/consumer/event-consumer"
	eventTelegram "github.com/yegor2025/budgetManager/events/telegram"
	"github.com/yegor2025/budgetManager/storage/googleSheets"
	"log"
	"os"
)

const (
	tgBotHost = "api.telegram.org"
	batchSize = 100
)

func main() {
	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("token is empty")
	}

	tgClient := telegram.New(tgBotHost, token)

	ctx := context.Background()
	storage := googleSheets.New(ctx, "./tokens.json")

	eventProcessor := eventTelegram.New(tgClient, storage)

	log.Printf("service started")

	consumer := event_consumer.New(eventProcessor, eventProcessor, 100)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}
