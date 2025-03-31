package telegram

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/yegor2025/budgetManager/cilents/telegram"
	"github.com/yegor2025/budgetManager/events"
	"github.com/yegor2025/budgetManager/storage"
	"log"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		log.Printf("Events.Fetch.Updates: %v", err)
		return nil, errors.Wrap(err, "can't get events")
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("list updates empty")
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return errors.Wrap(ErrUnknownEventType, "can't process message")
	}

}

func (p *Processor) processMessage(event events.Event) error {
	m, err := meta(event)
	if err != nil {
		return errors.Wrap(err, "can't process message")
	}

	if err = p.doCmd(event.Text, m.ChatID, m.Username); err != nil {
		return errors.Wrap(err, "can't process message")
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, errors.Wrap(ErrUnknownMetaType, "can't get meta")
	}

	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	return res
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}
