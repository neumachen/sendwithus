package sendwithus

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

// Send ...
func (c *Client) Send(ctx context.Context, getter SendPayloadGetter) (*SendResponse, error) {
	if IsNil(getter) {
		return nil, errors.New("send email payload is nil")
	}

	payload := SendPayload{}
	if err := MarshalSendPayload(getter, &payload); err != nil {
		return nil, err
	}
	req, err := c.NewRequest(http.MethodPost, "/api/v1/send", &payload)
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
