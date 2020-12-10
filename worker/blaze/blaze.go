package blaze

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"strings"
	"time"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/logger"
	"github.com/fox-one/walle/core"
	"github.com/fox-one/walle/pkg/cmd/root"
	"github.com/google/shlex"
)

func New(
	client *mixin.Client,
	brokers core.BrokerStore,
	brokerz core.BrokerService,
	messages core.MessageStore,
	render core.Render,
) *Blaze {
	return &Blaze{
		client:   client,
		brokers:  brokers,
		brokerz:  brokerz,
		messages: messages,
		render:   render,
	}
}

type Blaze struct {
	client   *mixin.Client
	brokers  core.BrokerStore
	brokerz  core.BrokerService
	messages core.MessageStore
	render   core.Render
}

func (w *Blaze) Run(ctx context.Context) error {
	log := logger.FromContext(ctx).WithField("worker", "blaze")
	ctx = logger.WithContext(ctx, log)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
			_ = w.client.LoopBlaze(ctx, w)
		}
	}
}

func (w *Blaze) OnAckReceipt(ctx context.Context, msg *mixin.MessageView, userID string) error {
	return nil
}

func (w *Blaze) OnMessage(ctx context.Context, msg *mixin.MessageView, userID string) error {
	if msg.Category != mixin.MessageCategoryPlainText {
		return nil
	}

	log := logger.FromContext(ctx).WithField("msg", msg.MessageID)

	data, err := base64.StdEncoding.DecodeString(msg.Data)
	if err != nil {
		return err
	}

	action := string(data)
	args := parseAction(action)

	b := &blazeBuilder{
		brokers: w.brokers,
		brokerz: w.brokerz,
		render:  w.render,
		msg:     msg,
	}

	cmd := root.NewCmd(b, "")
	cmd.SetArgs(args)
	writer := &messageWriter{msg: msg}
	cmd.SetOut(writer)
	cmd.SetErr(ioutil.Discard)

	if err := cmd.ExecuteContext(ctx); err != nil {
		log.WithError(err).Errorf("handle %q", action)
		return err
	}

	if err := w.messages.Create(ctx, writer.messages); err != nil {
		log.WithError(err).Errorln("messages.Create")
		return err
	}

	return nil
}

func parseAction(action string) []string {
	action = transformAction(action)
	args, _ := shlex.Split(action)
	return args
}

func transformAction(action string) string {
	keywords := map[string]string{
		"broker": "broker create",
	}

	for k, v := range keywords {
		if strings.EqualFold(k, action) {
			return v
		}
	}

	return action
}
