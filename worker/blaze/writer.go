package blaze

import (
	"encoding/json"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/uuid"
	"github.com/fox-one/walle/core"
)

type messageWriter struct {
	msg      *mixin.MessageView
	messages []*core.Message
}

func (m *messageWriter) Write(p []byte) (n int, err error) {
	var msg mixin.MessageRequest
	if err := json.Unmarshal(p, &msg); err != nil {
		return 0, err
	}

	msg.ConversationID = m.msg.ConversationID
	msg.MessageID = uuid.Modify(m.msg.MessageID, msg.MessageID)
	msg.RecipientID = m.msg.UserID

	m.messages = append(m.messages, core.BuildMessage(&msg))
	return 0, nil
}
