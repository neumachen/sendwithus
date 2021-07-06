package sendwithus

// Recipient ...
type Recipient struct {
	Address string `json:"address,omitempty"`
	Name    string `json:"name,omitempty"`
}

func (r *Recipient) GetAddress() string {
	return r.Address
}

func (r *Recipient) GetName() string {
	return r.Name
}

func (r *Recipient) SetAddress(address string) error {
	r.Address = address
	return nil
}

func (r *Recipient) SetName(name string) error {
	r.Name = name
	return nil
}

var _ RecipientAccessor = (*Recipient)(nil)

type RecipientGetter interface {
	GetAddress() string
	GetName() string
}

type RecipientGetters []RecipientGetter

func (r RecipientGetters) GetLength() int {
	return len(r)
}

type RecipientSetter interface {
	SetAddress(address string) error
	SetName(name string) error
}

type RecipientAccessor interface {
	RecipientGetter
	RecipientSetter
}

type RecipientAccessors []RecipientAccessor

func (r RecipientAccessors) GetLength() int {
	return len(r)
}

func (r RecipientAccessors) ToGetters() RecipientGetters {
	if r.GetLength() < 1 {
		return nil
	}

	gg := make(RecipientGetters, r.GetLength())
	for i := range r {
		gg[i] = r[i]
	}

	return gg
}

func MarshalRecipient(getter RecipientGetter, setter RecipientSetter) error {
	if OneIsNil(getter, setter) {
		return nil
	}
	if err := setter.SetAddress(getter.GetAddress()); err != nil {
		return err
	}
	return setter.SetName(getter.GetName())
}

type Recipients []Recipient

func (r *Recipients) Append(getters ...RecipientGetter) error {
	if len(getters) < 1 {
		return nil
	}

	if IsNil(r) {
		*r = make(Recipients, 0)
	}

	for i := range getters {
		rec := Recipient{}
		if err := MarshalRecipient(getters[i], &rec); err != nil {
			return err
		}

		*r = insertToRecipients(*r, i, rec)
	}

	return nil
}

func (r Recipients) GetLength() int {
	return len(r)
}

func (r Recipients) AsAccessors() RecipientAccessors {
	if r.GetLength() < 1 {
		return nil
	}

	aa := make(RecipientAccessors, r.GetLength())
	for i := range r {
		aa[i] = &r[i]
	}
	return aa
}

func insertToRecipients(original Recipients, position int, value Recipient) Recipients {
	l := len(original)
	target := original
	if cap(original) == l {
		target = make(Recipients, l+1, l+10)
		copy(target, original[:position])
	} else {
		target = append(target, Recipient{})
	}
	copy(target[position+1:], original[position:])
	target[position] = value
	return target
}
