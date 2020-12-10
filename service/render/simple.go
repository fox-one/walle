package render

import (
	"github.com/fox-one/walle/core"
)

func Simple() core.Render {
	return New(Templates{
		BrokerCreated: simpleBrokerCreated,
	})
}

const (
	simpleBrokerCreated = `wallet id: {{.WalletID}}`
)
