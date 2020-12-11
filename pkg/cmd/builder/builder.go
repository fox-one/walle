package builder

import (
	"github.com/fox-one/pkg/store/db"
	"github.com/fox-one/walle/core"
)

type Builder interface {
	DB() *db.DB
	Brokers() core.BrokerStore
	Brokerz() core.BrokerService
	Render() core.Render
	Perm() core.Perm
	Executor() string
	TraceID() string
}
