package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type Handler struct {
	bot       *tgbotapi.BotAPI
	processor *Processor
}

func NewHandler(bot *tgbotapi.BotAPI) *Handler {
	return &Handler{
		bot:       bot,
		processor: NewProcessor(bot),
	}
}

func (h *Handler) HandlerUpdate(update tgbotapi.Update) {
	switch {
	case update.Message != nil:
		h.processor.ProcessMessage(update.Message)
	case update.CallbackQuery != nil:
		h.processor.ProcessCallback(update.CallbackQuery)
	default:
		log.Println("Unknown update type received")
	}
}
