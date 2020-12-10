package terminal

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/fox-one/walle/core"
)

var verboseBrokerCreated = heredoc.Doc(`
		✅ New Broker Created
		
		Wallet ID: {{.WalletID}}
	`)

func Verbose() core.Render {
	return New(Templates{
		BrokerCreated: verboseBrokerCreated,
	})
}
