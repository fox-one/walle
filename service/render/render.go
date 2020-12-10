package render

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

	return &render{
		brokerCreated: parseTemplate(templates.BrokerCreated),
	}
}

type render struct {
	brokerCreated *template.Template
}

func (r *render) BrokerCreated(w io.Writer, broker *core.Broker) error {
	return r.brokerCreated.Execute(w, broker)
}

func parseTemplate(tpl string) *template.Template {
	return template.Must(template.New("-").Parse(tpl))
}
