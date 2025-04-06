package service

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type BudgetService struct {
	bot *tgbotapi.BotAPI
}

func NewBudgetservice(bot *tgbotapi.BotAPI) *BudgetService {
	return &BudgetService{
		bot: bot,
	}
}
