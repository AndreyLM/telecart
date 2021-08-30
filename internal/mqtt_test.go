package internal

import (
	"context"
	"encoding/json"
	"telecart/pkg/svc"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMqtt(t *testing.T) {
	topic := "test"

	c := NewMqttClient("tcp://localhost:1883")
	if err := c.Connect(); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	count := 5
	j := 0
	err := c.Subscribe(topic, func(ctx context.Context, topic string, payload []byte) {
		require.Equal(t, topic, topic)
		require.Equal(t, []byte("msg"), payload, string(payload))
		j++
	})
	require.NoError(t, err)
	defer c.Unsubscribe(topic)

	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < count; i++ {
			c.Publish(topic, []byte("msg"))
			time.Sleep(1 * time.Second)
		}
	}()

	<-done
	require.NoError(t, err)
	if j < count {
		t.Fatal("didn't receive all messages", j, count)
	}
}

func TestMqttPublish(t *testing.T) {
	topic := "topic"

	c := NewMqttClient("tcp://localhost:1883", "test-cli")
	if err := c.Connect(); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	msg := svc.Message{
		Name:    "user2",
		Message: "message2",
	}
	buf, err := json.Marshal(&msg)
	require.NoError(t, err)

	err = c.Publish(topic, buf)

	require.NoError(t, err)
	time.Sleep(1 * time.Second)
}
