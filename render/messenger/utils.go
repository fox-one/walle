package messenger

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"math/rand"

	"github.com/fox-one/mixin-sdk-go"
)

func write(w io.Writer, msg *mixin.MessageRequest) error {
	return json.NewEncoder(w).Encode(msg)
}

func writeText(w io.Writer, id, text string) error {
	msg := &mixin.MessageRequest{
		MessageID: id,
		Category:  mixin.MessageCategoryPlainText,
		Data:      base64.StdEncoding.EncodeToString([]byte(text)),
	}

	return write(w, msg)
}

func writeActions(w io.Writer, id string, args ...string) error {
	var buttons mixin.AppButtonGroupMessage

	for idx := 0; idx < len(args)-1; idx += 2 {
		buttons = append(buttons, mixin.AppButtonMessage{
			Label:  args[idx],
			Action: args[idx+1],
			Color:  randomHexColor(),
		})
	}

	data, _ := json.Marshal(buttons)
	msg := &mixin.MessageRequest{
		MessageID: id,
		Category:  mixin.MessageCategoryAppButtonGroup,
		Data:      base64.StdEncoding.EncodeToString(data),
	}

	return write(w, msg)
}

var hexColors = []string{
	"#FF93C9",
	"#FF4D00",
	"#0BAAFF",
	"#008080",
	"#5AC18E",
	"#0066CC",
	"#FD5392",
}

func randomHexColor() string {
	return hexColors[rand.Intn(len(hexColors))]
}
