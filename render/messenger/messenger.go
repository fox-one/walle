package messenger

import (
	"io"

	"github.com/fox-one/walle/core"
	"github.com/oxtoacart/bpool"
)

func Wrap(text core.Render) core.Render {
	return &messenger{
		text: text,
		pool: bpool.NewSizedBufferPool(16, 1024),
	}
}

type messenger struct {
	text core.Render
	pool *bpool.SizedBufferPool
}

func (m *messenger) BrokerCreated(w io.Writer, broker *core.Broker) error {
	b := m.pool.Get()
	defer m.pool.Put(b)

	if err := m.text.BrokerCreated(b, broker); err != nil {
		return err
	}

	return writeText(w, broker.WalletID, b.String())
}
