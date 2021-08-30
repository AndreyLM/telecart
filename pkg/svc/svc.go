package svc

import (
	"context"
	"encoding/json"
	"log"

	"github.com/pkg/errors"
)

type (
	Message struct {
		ID      int
		Message string
		Name    string
	}

	MsgService struct {
		cli   MsqClient
		store Store
	}

	Store interface {
		Init() error
		Save(ctx context.Context, msg *Message) error
		Close() error
	}

	SubCallback func(ctx context.Context, topic string, payload []byte)

	MsqClient interface {
		Subscribe(topic string, callback SubCallback) error
		Unsubscribe(topic string) error
		Publish(topic string, payload []byte) error
		Connect() error
		Close() error
	}
)

func NewMsgService(cli MsqClient, store Store) *MsgService {
	return &MsgService{
		cli:   cli,
		store: store,
	}
}

func (s *MsgService) Init() error {
	if err := s.cli.Connect(); err != nil {
		return errors.Wrap(err, "connect")
	}

	return nil
}

func (s *MsgService) Subscribe(topic string) (err error) {
	err = s.cli.Subscribe(topic, func(ctx context.Context, topic string, payload []byte) {
		msg := &Message{}
		if err := json.Unmarshal(payload, &msg); err != nil {
			log.Println("decode payload", "Error:", err, "payload:", string(payload))
			return
		}
		if err := s.store.Save(ctx, msg); err != nil {
			log.Println("save data", err, topic)
		}
	})

	return err
}

func (s *MsgService) Close() (err error) {
	err = s.cli.Close()

	if dbErr := s.store.Close(); dbErr != nil {
		if err != nil {
			return errors.Wrap(err, dbErr.Error())
		}

		err = dbErr
	}

	return
}

func (s *MsgService) Wait(ctx context.Context) error {
	<-ctx.Done()
	return s.Close()
}
