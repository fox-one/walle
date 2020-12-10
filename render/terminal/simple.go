package terminal

import (
	"github.com/fox-one/walle/core"
)

const simpleBrokerCreated = `wallet id: {{.WalletID}}`

func Simple() core.Render {
	return New(Templates{
		BrokerCreated: simpleBrokerCreated,
	})
}
