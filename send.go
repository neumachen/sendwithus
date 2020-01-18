package sendwithus

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

// Email ...
type Email struct {
	Template     string           `json:"template,omitempty"`
	Recipient    *Recipient       `json:"recipient,omitempty"`
	CC           []Recipient      `json:"cc,omitempty"`
	BCC          []Recipient      `json:"bcc,omitempty"`
	Sender       *Sender          `json:"sender,omitempty"`
	TemplateData *json.RawMessage `json:"template_data,omitempty"`
	Tags         []string         `json:"tags,omitempty"`
	Inline       *Attachment      `json:"inline,omitempty"`
	Files        []Attachment     `json:"files,omitempty"`
	ESPAccount   string           `json:"esp_account,omitempty"`
	VersionName  string           `json:"version_name,omitempty"`
}

// SendResponse ...
type SendResponse struct {
	Success   bool   `json:"success"`
	Status    string `json:"status"`
	ReceiptID string `json:"receipt_id"`
	Email     struct {
		Name        string `json:"name"`
		VersionName string `json:"version_name"`
		Locale      string `json:"locale"`
	} `json:"email"`
}

// Marshal ...
func (e *Email) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

// Recipient ...
type Recipient struct {
	Address string `json:"address,omitempty"`
	Name    string `json:"name,omitempty"`
}

// Sender ...
type Sender struct {
	Recipient
	ReplyTo string `json:"reply_to,omitempty"`
}

// Attachment ...
type Attachment struct {
	ID   string `json:"id,omitempty"`
	Data string `json:"data,omitempty"`
}

// Send ...
func (c *Client) Send(ctx context.Context, email *Email) (*SendResponse, error) {
	req, err := c.NewRequest("POST", "/api/v1/send", email)
	if err != nil {
		return nil, err
	}
	b, err := c.Do(ctx, http.StatusOK, req)
	if err != nil {
		return nil, err
	}
	sr := SendResponse{}
	err = json.Unmarshal(b, &sr)
	if err != nil {
		return nil, err
	}

	if !sr.Success {
		return nil, errors.New("response returned false")
	}
	return &sr, nil
}
