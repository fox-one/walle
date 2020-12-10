package core

import (
	"io"
)

type Render interface {
	BrokerCreated(w io.Writer, broker *Broker) error
}
