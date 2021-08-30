package internal

import (
	"telecart/pkg/svc"

	"github.com/pkg/errors"
)

const brockerDNS = "tcp://localhost:1883"

func NewService(dbDns, topic string) (*svc.MsgService, error) {
	store, err := NewSqlLiteStore(dbDns)
	if err != nil {
		return nil, errors.Wrap(err, "new store")
	}

	if err := store.Init(); err != nil {
		return nil, errors.Wrap(err, "init store")
	}

	msqCli := NewMqttClient(brockerDNS)
	s := svc.NewMsgService(msqCli, store)

	if err := s.Init(); err != nil {
		return nil, errors.Wrap(err, "init service")
	}

	if err := s.Subscribe(topic); err != nil {
		return nil, errors.Wrap(err, "subscribe")
	}

	return s, nil
}
