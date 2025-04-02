package telegram

import (
	"github.com/pkg/errors"
	"github.com/yegor2025/budgetManager/cilents/telegram"
	"log"
	"strings"
)

const (
	RndCmd   = "/rnd"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username)

	return p.saveRecord(chatID, text, username)

	//switch text {
	//case RndCmd:
	//	return p.saveRecord(chatID, text, username)
	//case StartCmd:
	//	return p.sendHello(chatID)
	//default:
	//	return p.tg.SendMessage(chatID, msgUnknownCommand)
	//}
}

func (p *Processor) saveRecord(chatID int, text string, username string) (err error) {
	defer func() { err = errors.Wrap(err, "can't do command add record") }()

	sendMsg := NewMessageSender(chatID, p.tg)

	if err = p.storage.InsertBeforeLast(text, false); err != nil {
		return err
	}

	//if err = p.storage.Append(text); err != nil {
	//	return err
	//}

	if err = sendMsg(msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func NewMessageSender(chatId int, tg *telegram.Client) func(string) error {
	return func(msg string) error {
		return tg.SendMessage(chatId, msg)
	}
}
