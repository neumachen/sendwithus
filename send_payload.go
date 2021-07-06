package sendwithus

import "encoding/json"

// SendPayload ...
type SendPayload struct {
	Template     string           `json:"template,omitempty"`
	Recipient    *Recipient       `json:"recipient,omitempty"`
	CC           Recipients       `json:"cc,omitempty"`
	BCC          Recipients       `json:"bcc,omitempty"`
	Sender       *Sender          `json:"sender,omitempty"`
	TemplateData *json.RawMessage `json:"template_data,omitempty"`
	Tags         Tags             `json:"tags,omitempty"`
	Inline       *Attachment      `json:"inline,omitempty"`
	Files        Attachments      `json:"files,omitempty"`
	ESPAccount   string           `json:"esp_account,omitempty"`
	VersionName  string           `json:"version_name,omitempty"`
}

// Marshal ...
func (s *SendPayload) Marshal() ([]byte, error) {
	return json.Marshal(s)
}

func (s *SendPayload) GetTemplate() string {
	return s.Template
}

func (s *SendPayload) GetRecipient() RecipientAccessor {
	return s.Recipient
}

func (s *SendPayload) GetCC() RecipientAccessors {
	if s.CC.GetLength() < 1 {
		return nil
	}
	return s.CC.AsAccessors()
}

func (s *SendPayload) GetBCC() RecipientAccessors {
	if s.BCC.GetLength() < 1 {
		return nil
	}
	return s.BCC.AsAccessors()
}

func (s *SendPayload) GetSender() SenderAccessor {
	return s.Sender
}

func (s *SendPayload) GetTemplateData() *json.RawMessage {
	return s.TemplateData
}

func (s *SendPayload) GetTags() Tags {
	return s.Tags
}

func (s *SendPayload) GetInline() AttachmentAccessor {
	return s.Inline
}

func (s *SendPayload) GetFiles() AttachmentAccessors {
	if s.Files.GetLength() < 1 {
		return nil
	}
	return s.Files.AsAccessors()
}

func (s *SendPayload) GetESPAccount() string {
	return s.ESPAccount
}

func (s *SendPayload) GetVersionName() string {
	return s.VersionName
}

func (s *SendPayload) SetTemplate(template string) error {
	s.Template = template
	return nil
}

func (s *SendPayload) SetRecipient(getter RecipientGetter) error {
	if IsNil(getter) {
		return nil
	}
	if s.Recipient == nil {
		s.Recipient = &Recipient{}
	}

	return MarshalRecipient(getter, s.Recipient)
}

func (s *SendPayload) SetCC(getters ...RecipientGetter) error {
	if len(getters) < 1 {
		return nil
	}
	return s.CC.Append(getters...)
}

func (s *SendPayload) SetBCC(getters ...RecipientGetter) error {
	if len(getters) < 1 {
		return nil
	}
	return s.BCC.Append(getters...)
}

func (s *SendPayload) SetSender(getter SenderGetter) error {
	if IsNil(getter) {
		return nil
	}
	if s.Sender == nil {
		s.Sender = &Sender{}
	}

	return MarshalSender(getter, s.Sender)
}

func (s *SendPayload) SetTemplateData(templateData *json.RawMessage) error {
	s.TemplateData = templateData
	return nil
}

func (s *SendPayload) SetTags(tags ...string) error {
	if len(tags) < 1 {
		return nil
	}
	s.Tags.Append(tags...)
	return nil
}

func (s *SendPayload) SetInline(getter AttachmentGetter) error {
	if IsNil(getter) {
		return nil
	}
	if s.Inline == nil {
		s.Inline = &Attachment{}
	}

	return MarshalAttachment(getter, s.Inline)
}

func (s *SendPayload) SetFiles(getters ...AttachmentGetter) error {
	if len(getters) < 1 {
		return nil
	}
	return s.Files.Append(getters...)
}

func (s *SendPayload) SetESPAccount(espAccount string) error {
	s.ESPAccount = espAccount
	return nil
}

func (s *SendPayload) SetVersionName(versionName string) error {
	s.VersionName = versionName
	return nil
}

var _ SendPayloadAccessor = (*SendPayload)(nil)

type SendPayloadGetter interface {
	Marshaler
	GetBCC() RecipientAccessors
	GetCC() RecipientAccessors
	GetESPAccount() string
	GetFiles() AttachmentAccessors
	GetInline() AttachmentAccessor
	GetRecipient() RecipientAccessor
	GetSender() SenderAccessor
	GetTags() Tags
	GetTemplate() string
	GetTemplateData() *json.RawMessage
	GetVersionName() string
}

type SendPayloadSetter interface {
	SetBCC(getters ...RecipientGetter) error
	SetCC(getters ...RecipientGetter) error
	SetESPAccount(espAccount string) error
	SetFiles(getters ...AttachmentGetter) error
	SetInline(getter AttachmentGetter) error
	SetRecipient(getter RecipientGetter) error
	SetSender(getter SenderGetter) error
	SetTags(tags ...string) error
	SetTemplate(template string) error
	SetTemplateData(templateData *json.RawMessage) error
	SetVersionName(versionName string) error
}

type SendPayloadAccessor interface {
	SendPayloadGetter
	SendPayloadSetter
}

func MarshalSendPayload(getter SendPayloadGetter, setter SendPayloadSetter) error {
	if OneIsNil(getter, setter) {
		return nil
	}

	if bcc := getter.GetBCC(); !IsNil(bcc) && bcc.GetLength() > 0 {
		if err := setter.SetBCC(bcc.ToGetters()...); err != nil {
			return err
		}
	}
	if cc := getter.GetCC(); !IsNil(cc) && cc.GetLength() > 0 {
		if err := setter.SetCC(cc.ToGetters()...); err != nil {
			return err
		}
	}
	if err := setter.SetESPAccount(getter.GetESPAccount()); err != nil {
		return err
	}
	if files := getter.GetFiles(); !IsNil(files) && files.GetLength() > 0 {
		if err := setter.SetFiles(files.ToGetters()...); err != nil {
			return err
		}
	}
	if err := setter.SetInline(getter.GetInline()); err != nil {
		return err
	}
	if err := setter.SetRecipient(getter.GetRecipient()); err != nil {
		return err
	}
	if err := setter.SetSender(getter.GetSender()); err != nil {
		return err
	}
	if tags := getter.GetTags(); !IsNil(tags) && tags.GetLength() > 0 {
		if err := setter.SetTags(tags.ToStrings()...); err != nil {
			return err
		}
	}
	if err := setter.SetTemplate(getter.GetTemplate()); err != nil {
		return err
	}
	if err := setter.SetTemplateData(getter.GetTemplateData()); err != nil {
		return err
	}

	return setter.SetVersionName(getter.GetVersionName())
}
