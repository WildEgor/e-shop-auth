package adapters

import (
	"context"
	"github.com/WildEgor/e-shop-auth/internal/configs"
	"github.com/WildEgor/e-shop-gopack/pkg/libs/notifier"
	"log/slog"
)

type NotifierAdapter struct {
	client *notifier.NotifierClient
}

func NewNotifierAdapter(config *configs.NotifierConfig) *NotifierAdapter {

	client, err := notifier.NewNotifierClient(&notifier.NotifierConfig{
		DSN:      config.DSN,
		Exchange: config.Exchange,
	})
	if err != nil {
		panic(err) // TODO: handle error
	}

	return &NotifierAdapter{
		client,
	}
}

func (n *NotifierAdapter) Notify(payload *notifier.NotificationPayload) error {
	return n.client.Notify(context.TODO(), payload)
}

func (n *NotifierAdapter) Close() {
	if err := n.client.Close(); err != nil {
		slog.Error("unable to close notifier.", err)
		panic(err)
	}
}
