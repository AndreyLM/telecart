package internal

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSvc(t *testing.T) {
	sqlDNS := "test"
	topic := "test"

	svc, err := NewService(sqlDNS, "topic")
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		pub := NewMqttClient(brockerDNS)

		str := struct {
			Name    string
			Message string
		}{
			Name:    "user",
			Message: "message",
		}

		buf, err := json.Marshal(&str)
		require.NoError(t, err)
		pub.Publish(topic, buf)
		time.Sleep(1 * time.Second)
		cancel()
	}()

	if err := svc.Wait(ctx); err != nil {
		t.Fatal(err)
	}

	t.Log("Success....")
}
