package sendwithus

// Attachment ...
type Attachment struct {
	ID   string `json:"id,omitempty"`
	Data string `json:"data,omitempty"`
}

func (a *Attachment) GetID() string {
	return a.ID
}

func (a *Attachment) GetData() string {
	return a.Data
}

func (a *Attachment) SetID(address string) error {
	a.ID = address
	return nil
}

func (a *Attachment) SetData(name string) error {
	a.Data = name
	return nil
}

var _ AttachmentAccessor = (*Attachment)(nil)

type AttachmentGetter interface {
	GetID() string
	GetData() string
}

type AttachmentGetters []AttachmentGetter

func (a AttachmentGetters) GetLength() int {
	return len(a)
}

type AttachmentSetter interface {
	SetID(address string) error
	SetData(name string) error
}

type AttachmentAccessor interface {
	AttachmentGetter
	AttachmentSetter
}

type AttachmentAccessors []AttachmentAccessor

func (a AttachmentAccessors) GetLength() int {
	return len(a)
}

func (a AttachmentAccessors) ToGetters() AttachmentGetters {
	if a.GetLength() < 1 {
		return nil
	}

	gg := make(AttachmentGetters, a.GetLength())
	for i := range a {
		gg[i] = a[i]
	}

	return gg
}

func MarshalAttachment(getter AttachmentGetter, setter AttachmentSetter) error {
	if OneIsNil(getter, setter) {
		return nil
	}
	if err := setter.SetID(getter.GetID()); err != nil {
		return err
	}
	return setter.SetData(getter.GetData())
}

type Attachments []Attachment

func (a Attachments) GetLength() int {
	return len(a)
}

func (a Attachments) AsAccessors() AttachmentAccessors {
	if a.GetLength() < 1 {
		return nil
	}

	aa := make(AttachmentAccessors, a.GetLength())
	for i := range a {
		aa[i] = &a[i]
	}
	return aa
}

func (a *Attachments) Append(getters ...AttachmentGetter) error {
	if len(getters) < 1 {
		return nil
	}

	if IsNil(a) {
		*a = make(Attachments, 0)
	}

	for i := range getters {
		rec := Attachment{}
		if err := MarshalAttachment(getters[i], &rec); err != nil {
			return err
		}

		*a = insertToAttachments(*a, i, rec)
	}

	return nil
}

func insertToAttachments(original Attachments, position int, value Attachment) Attachments {
	l := len(original)
	target := original
	if cap(original) == l {
		target = make(Attachments, l+1, l+10)
		copy(target, original[:position])
	} else {
		target = append(target, Attachment{})
	}
	copy(target[position+1:], original[position:])
	target[position] = value
	return target
}
