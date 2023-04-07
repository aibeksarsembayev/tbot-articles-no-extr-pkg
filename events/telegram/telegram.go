package telegram

import (
	"errors"

	"github.com/aibeksarsembayev/tbot-articles-no-extr-pkg/clients/telegram"
	"github.com/aibeksarsembayev/tbot-articles-no-extr-pkg/events"
	"github.com/aibeksarsembayev/tbot-articles-no-extr-pkg/lib/e"
	"github.com/aibeksarsembayev/tbot-articles-no-extr-pkg/storage"
	"go.uber.org/zap"
)

type Processor struct {
	lg      *zap.Logger
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID            int
	Username          string
	Category          string
	Author            string
	CallbackQuery     string
	Status            string
	ChannelPost       string
	EditedChannelPost string
	PollID            string
}

type CallbackQuery struct {
	ID   string
	Data string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(logger *zap.Logger, client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		lg:      logger,
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
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
	case events.CallbackQuery:
		return p.processMessage(event)
	case events.MyChatMember:
		return p.processMessage(event)
	case events.ChannelPost:
		return p.processMessage(event)
	case events.EditedChannelPost:
		return p.processMessage(event)
	case events.EditedMessage:
		return p.processMessage(event)
	case events.Poll:
		return p.processMessage(event)
	default:
		return e.Wrap("can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}
	// fmt.Println("meta ", meta, "event", event)
	if err := p.doCmd(event.Text, meta); err != nil {
		return e.Wrap("can't process message", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	// in case of unknown message structure, use below return
	// return events.Event{}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	if updType == events.EditedMessage {
		res.Meta = Meta{
			ChatID:   upd.EditedMessage.Chat.ID,
			Username: upd.EditedMessage.From.Username,
		}
	}

	if updType == events.ChannelPost {
		if upd.ChannelPost.Text == "" {
			upd.ChannelPost.Text = "channelpost"
		}
		res.Meta = Meta{
			ChatID:      upd.ChannelPost.Chat.ID,
			Username:    upd.ChannelPost.Chat.Username,
			ChannelPost: upd.ChannelPost.Text,
		}
	}

	if updType == events.EditedChannelPost {
		if upd.EditedChannelPost.Text == "" {
			upd.EditedChannelPost.Text = "editedchannelpost"
		}
		res.Meta = Meta{
			ChatID:      upd.EditedChannelPost.Chat.ID,
			Username:    upd.EditedChannelPost.Chat.Username,
			ChannelPost: upd.EditedChannelPost.Text,
		}
	}

	if updType == events.Poll {
		res.Meta = Meta{
			PollID: upd.Poll.ID,
		}
	}

	if updType == events.CallbackQuery {
		res.Meta = Meta{
			ChatID:        upd.CallbackQuery.Message.Chat.ID,
			Username:      upd.CallbackQuery.From.Username,
			CallbackQuery: upd.CallbackQuery.Data,
		}
	}

	if updType == events.MyChatMember {
		res.Meta = Meta{
			ChatID:   upd.MyChatMember.Chat.ID,
			Username: upd.MyChatMember.From.Username,
			Status:   upd.MyChatMember.NewChatMember.Status,
		}
	}

	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil && upd.EditedMessage == nil {
		return ""
	}

	if upd.EditedMessage != nil {
		return upd.EditedMessage.Text
	}

	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message != nil {
		return events.Message
	}

	if upd.EditedMessage != nil {
		return events.EditedMessage
	}

	if upd.ChannelPost != nil {
		return events.ChannelPost
	}

	if upd.EditedChannelPost != nil {
		return events.EditedChannelPost
	}

	if upd.InlineQuery != nil {
		return events.InlineQuery
	}

	if upd.ChosenInlineResult != nil {
		return events.ChosenInlineResult
	}

	if upd.CallbackQuery != nil {
		return events.CallbackQuery
	}

	if upd.ShippingQuery != nil {
		return events.ShippingQuery
	}

	if upd.PreCheckoutQuery != nil {
		return events.PreCheckoutQuery
	}

	if upd.Poll != nil {
		return events.Poll
	}

	if upd.PollAnswer != nil {
		return events.PollAnswer
	}

	if upd.MyChatMember != nil {
		return events.MyChatMember
	}

	if upd.ChatMember != nil {
		return events.ChatMember
	}

	if upd.ChatJoinRequest != nil {
		return events.ChatJoinRequest
	}

	return events.Message
}
