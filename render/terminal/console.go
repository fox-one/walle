package terminal

import (
	"io"
	"text/template"

	"github.com/asaskevich/govalidator"
	"github.com/fox-one/walle/core"
)

type Templates struct {
	BrokerCreated string `valid:"required"`
}

func New(templates Templates) core.Render {
	if _, err := govalidator.ValidateStruct(templates); err != nil {
		panic(err)
	}

	return &console{
		brokerCreated: parseTemplate(templates.BrokerCreated),
	}
}

type console struct {
	brokerCreated *template.Template
}

func (c *console) BrokerCreated(w io.Writer, broker *core.Broker) error {
	return c.brokerCreated.Execute(w, broker)
}
