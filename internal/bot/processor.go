package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/yegor2025/budgetManager/internal/service"
	"log"
)

type Processor struct {
	bot           *tgbotapi.BotAPI
	budgetService *service.BudgetService
}

func NewProcessor(bot *tgbotapi.BotAPI) *Processor {
	return &Processor{
		bot:           bot,
		budgetService: service.NewBudgetservice(bot),
	}
}

func (p *Processor) ProcessMessage(msg *tgbotapi.Message) {
	switch msg.Text {
	case "/start":
		p.sendStartMenu(msg.Chat.ID)
	case "📥 Новое сообщение":
		p.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Напиши новое сообщение!"))
	case "⚙ Настройки":
		p.bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Здесь будут настройки."))
	case "❌ Скрыть клавиатуру":
		remove := tgbotapi.NewRemoveKeyboard(true)
		hideMsg := tgbotapi.NewMessage(msg.Chat.ID, "Клавиатура скрыта.")
		hideMsg.ReplyMarkup = remove
		p.bot.Send(hideMsg)
	default:
		// По умолчанию — эхо
	}
}

func (p *Processor) sendStartMenu(chatID int64) {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📥 Новое сообщение"),
			tgbotapi.NewKeyboardButton("⚙ Настройки"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("❌ Скрыть клавиатуру"),
		),
	)
	keyboard.ResizeKeyboard = true // Автоматически подгоняет размер

	msg := tgbotapi.NewMessage(chatID, "Добро пожаловать!")
	msg.ReplyMarkup = keyboard

	p.bot.Send(msg)
}

func (p *Processor) ProcessCallback(cb *tgbotapi.CallbackQuery) {
	log.Printf("Received callback: %s", cb.Data)
	msg := tgbotapi.NewMessage(cb.Message.Chat.ID, "Ты нажал кнопку: "+cb.Data)
	p.bot.Send(msg)
}
