package internal

import (
	"context"
	"telecart/pkg/svc"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	s, err := NewSqlLiteStore("test1.db")
	require.NoError(t, err)
	err = s.Init()
	require.NoError(t, err)

	err = s.Save(context.Background(), &svc.Message{
		Name:    "s",
		Message: "s",
	})

	require.NoError(t, err)
}
