package render

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/fox-one/walle/core"
)

func Verbose() core.Render {
	return New(Templates{
		BrokerCreated: verboseBrokerCreated,
	})
}

var (
	verboseBrokerCreated = heredoc.Doc(`
		✅ New Broker Created
		
		Wallet ID: {{.WalletID}}
	`)
)
