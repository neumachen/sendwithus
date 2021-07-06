package sendwithus

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	client, err := NewClient(os.Getenv("SENDWITHUS_TEST_API_KEY"), nil)
	require.NoError(t, err)

	sendPayload := SendPayload{}
	sendPayload.Template = os.Getenv("SENDWITHUS_TEST_TEMPLATE")
	sendPayload.Recipient = &Recipient{
		Name:    "Test Recipient",
		Address: os.Getenv("SENDWITHUS_TEST_SENDER"),
	}
	sendPayload.CC = Recipients{
		{
			Name:    "Test CC",
			Address: "kareem@joinpara.com",
		},
	}
	sendPayload.BCC = Recipients{
		{
			Name:    "Test BCC",
			Address: os.Getenv("SENDWITHUS_TEST_RECEIPIENT"),
		},
	}
	sendPayload.Sender = &Sender{
		Recipient: Recipient{
			Name:    "Test Sender",
			Address: os.Getenv("SENDWITHUS_TEST_RECEIPIENT"),
		},
	}

	td := map[string]interface{}{
		"test": "test",
	}

	b, err := json.Marshal(td)
	require.NoError(t, err)
	require.NotNil(t, b)

	resp, err := client.Send(context.Background(), &sendPayload)
	require.NoError(t, err)
	require.NotNil(t, resp)
}
