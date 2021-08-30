package internal

import (
	"context"
	"fmt"
	"log"
	"telecart/pkg/svc"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
)

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("Undefined subscription. TOPIC: %s\nMSG: %s\n", msg.Topic(), msg.Payload())
}

type MqttClient struct {
	cli MQTT.Client
}

func NewMqttClient(dns string, params ...string) svc.MsqClient {
	cliID := "go-telecart"
	if len(params) > 0 {
		cliID = params[0]
	}
	opts := MQTT.NewClientOptions().
		AddBroker(dns).
		SetUsername("user").
		SetPassword("password").
		SetClientID(cliID).
		SetDefaultPublishHandler(f).
		SetCleanSession(true)

	return &MqttClient{
		cli: MQTT.NewClient(opts),
	}
}

func (c *MqttClient) Connect() error {
	if token := c.cli.Connect(); token.Wait() && token.Error() != nil {
		return errors.Wrap(token.Error(), "create connection")
	}

	return nil
}

func (c *MqttClient) Unsubscribe(topic string) error {
	if token := c.cli.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		return errors.Wrap(token.Error(), "Unsubscribe "+topic)
	}

	return nil
}

func (c *MqttClient) Close() error {
	c.cli.Disconnect(100)
	return nil
}

func (c *MqttClient) Subscribe(topic string, callback svc.SubCallback) error {
	token := c.cli.Subscribe(topic, 0, func(client MQTT.Client, msg MQTT.Message) {
		log.Println("Payload", string(msg.Payload()))
		ctx := context.Background()
		callback(ctx, msg.Topic(), msg.Payload())
	})

	if token.Wait() && token.Error() != nil {
		return errors.Wrap(token.Error(), "subscribe to "+topic)
	}

	return nil
}

func (c *MqttClient) Publish(topic string, payload []byte) error {
	if token := c.cli.Publish(topic, 0, true, payload); token.Wait() && token.Error() != nil {
		return errors.Wrap(token.Error(), "publish")
	}

	return nil
}
