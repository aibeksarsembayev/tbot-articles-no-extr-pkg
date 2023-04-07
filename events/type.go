package events

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	Process(e Event) error
}

type Type int

const (
	Unknown Type = iota
	Message
	EditedMessage
	ChannelPost
	EditedChannelPost
	InlineQuery
	ChosenInlineResult
	CallbackQuery
	ShippingQuery
	PreCheckoutQuery
	Poll
	PollAnswer
	MyChatMember
	ChatMember
	ChatJoinRequest
)

type Event struct {
	Type Type
	Text string
	Meta interface{}
}
