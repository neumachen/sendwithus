package sendwithus

// Sender ...
type Sender struct {
	Recipient
	ReplyTo string `json:"reply_to,omitempty"`
}

func (s *Sender) GetRecipient() RecipientAccessor {
	return &s.Recipient
}

func (s *Sender) GetReplyTo() string {
	return s.ReplyTo
}

func (s *Sender) SetRecipient(getter RecipientGetter) error {
	return MarshalRecipient(getter, &s.Recipient)
}

func (s *Sender) SetReplyTo(replyTo string) error {
	s.ReplyTo = replyTo
	return nil
}

var _ SenderAccessor = (*Sender)(nil)

type SenderGetter interface {
	GetRecipient() RecipientAccessor
	GetReplyTo() string
}

type SenderSetter interface {
	SetRecipient(getter RecipientGetter) error
	SetReplyTo(name string) error
}

type SenderAccessor interface {
	SenderGetter
	SenderSetter
}

func MarshalSender(getter SenderGetter, setter SenderSetter) error {
	if OneIsNil(getter, setter) {
		return nil
	}
	if err := setter.SetRecipient(getter.GetRecipient()); err != nil {
		return err
	}
	return setter.SetReplyTo(getter.GetReplyTo())
}

type Senders []Sender

func (s Senders) GetLength() int {
	return len(s)
}
