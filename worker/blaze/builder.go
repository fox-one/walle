package blaze

import (
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/pkg/store/db"
	"github.com/fox-one/walle/core"
)

type blazeBuilder struct {
	brokers core.BrokerStore
	brokerz core.BrokerService
	render  core.Render
	perm    core.Perm

	msg *mixin.MessageView
}

func (b *blazeBuilder) DB() *db.DB {
	panic("implement me")
}

func (b *blazeBuilder) Brokers() core.BrokerStore {
	return b.brokers
}

func (b *blazeBuilder) Brokerz() core.BrokerService {
	return b.brokerz
}

func (b *blazeBuilder) Render() core.Render {
	return b.render
}

func (b *blazeBuilder) Executor() string {
	return b.msg.UserID
}

func (b *blazeBuilder) TraceID() string {
	return b.msg.MessageID
}

func (b *blazeBuilder) Perm() core.Perm {
	return b.perm
}
