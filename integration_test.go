package sendwithus

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSend(t *testing.T) {
	client := NewClient(os.Getenv("SENDWITHUS_TEST_API_KEY"), nil)

	email := Email{}
	email.Template = os.Getenv("SENDWITHUS_TEST_TEMPLATE")
	email.Recipient = &Recipient{
		Name:    "Test Recipient",
		Address: os.Getenv("SENDWITHUS_TEST_SENDER"),
	}
	email.CC = []Recipient{
		Recipient{
			Name:    "Test CC",
			Address: "kareem@joinpara.com",
		},
	}
	email.BCC = []Recipient{
		Recipient{
			Name:    "Test BCC",
			Address: os.Getenv("SENDWITHUS_TEST_RECEIPIENT"),
		},
	}
	email.Sender = &Sender{
		Recipient: Recipient{
			Name:    "Test Sender",
			Address: os.Getenv("SENDWITHUS_TEST_RECEIPIENT"),
		},
	}

	td := map[string]interface{}{
		"test": "test",
	}

	b, err := json.Marshal(td)
	assert.NoError(t, err)
	assert.NotNil(t, b)

	resp, err := client.Send(context.Background(), &email)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
