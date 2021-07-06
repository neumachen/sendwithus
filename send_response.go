package sendwithus

import "go.uber.org/zap/zapcore"

type SendResponseEmail struct {
	Name        string `json:"name,omitempty"`
	VersionName string `json:"version_name,omitempty"`
	Locale      string `json:"locale,omitempty"`
}

func (s *SendResponseEmail) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if v := s.GetName(); !StringIsEmpty(v) {
		enc.AddString("name", v)
	}
	if v := s.GetVersionName(); !StringIsEmpty(v) {
		enc.AddString("version_name", v)
	}
	if v := s.GetLocale(); !StringIsEmpty(v) {
		enc.AddString("locale", v)
	}
	return nil
}

func (s *SendResponseEmail) GetName() string {
	return s.Name
}

func (s *SendResponseEmail) GetVersionName() string {
	return s.VersionName
}

func (s *SendResponseEmail) GetLocale() string {
	return s.Locale
}

func (s *SendResponseEmail) SetName(name string) error {
	s.Name = name
	return nil
}

func (s *SendResponseEmail) SetVersionName(versionName string) error {
	s.VersionName = versionName
	return nil
}

func (s *SendResponseEmail) SetLocale(locale string) error {
	s.Locale = locale
	return nil
}

var _ SendResponseEmailAccessor = (*SendResponseEmail)(nil)

type SendResponseEmailGetter interface {
	zapcore.ObjectMarshaler
	GetLocale() string
	GetName() string
	GetVersionName() string
}

type SendResponseEmailSetter interface {
	SetLocale(locale string) error
	SetName(name string) error
	SetVersionName(versionName string) error
}

type SendResponseEmailAccessor interface {
	SendResponseEmailGetter
	SendResponseEmailSetter
}

func MarshalSendResponseEmail(getter SendResponseEmailGetter, setter SendResponseEmailSetter) error {
	if OneIsNil(getter, setter) {
		return nil
	}

	if err := setter.SetName(getter.GetName()); err != nil {
		return err
	}
	if err := setter.SetVersionName(getter.GetVersionName()); err != nil {
		return err
	}

	return setter.SetLocale(getter.GetLocale())
}

// SendResponse ...
type SendResponse struct {
	Success   bool               `json:"success,omitempty"`
	Status    string             `json:"status,omitempty"`
	ReceiptID string             `json:"receipt_id,omitempty"`
	Email     *SendResponseEmail `json:"email,omitempty"`
}

func (s *SendResponse) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddBool("success", s.GetSuccess())
	if v := s.GetStatus(); !StringIsEmpty(v) {
		enc.AddString("status", v)
	}
	if v := s.GetReceiptID(); !StringIsEmpty(v) {
		enc.AddString("receipt_id", v)
	}
	if v := s.GetEmail(); !IsNil(v) {
		enc.AddObject("email", v)
	}
	return nil
}

func (s *SendResponse) GetSuccess() bool {
	return s.Success
}

func (s *SendResponse) GetStatus() string {
	return s.Status
}

func (s *SendResponse) GetReceiptID() string {
	return s.ReceiptID
}

func (s *SendResponse) GetEmail() SendResponseEmailAccessor {
	return s.Email
}

func (s *SendResponse) SetSuccess(success bool) error {
	s.Success = success
	return nil
}

func (s *SendResponse) SetStatus(status string) error {
	s.Status = status
	return nil
}

func (s *SendResponse) SetReceiptID(receiptID string) error {
	s.ReceiptID = receiptID
	return nil
}

func (s *SendResponse) SetEmail(getter SendResponseEmailGetter) error {
	if IsNil(getter) {
		return nil
	}
	if s.Email == nil {
		s.Email = &SendResponseEmail{}
	}

	return MarshalSendResponseEmail(getter, s.Email)
}

var _ SendResponseAccessor = (*SendResponse)(nil)

type SendResponseGetter interface {
	zapcore.ObjectMarshaler
	GetEmail() SendResponseEmailAccessor
	GetReceiptID() string
	GetStatus() string
	GetSuccess() bool
}

type SendResponseSetter interface {
	SetEmail(getter SendResponseEmailGetter) error
	SetReceiptID(receiptID string) error
	SetStatus(status string) error
	SetSuccess(success bool) error
}

type SendResponseAccessor interface {
	SendResponseGetter
	SendResponseSetter
}

func MarshalSendResponse(getter SendResponseGetter, setter SendResponseSetter) error {
	if OneIsNil(getter, setter) {
		return nil
	}
	if err := setter.SetEmail(getter.GetEmail()); err != nil {
		return err
	}
	if err := setter.SetReceiptID(getter.GetReceiptID()); err != nil {
		return err
	}
	if err := setter.SetStatus(getter.GetStatus()); err != nil {
		return err
	}
	return setter.SetSuccess(getter.GetSuccess())
}
