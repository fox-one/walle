package terminal

import "text/template"

func parseTemplate(tpl string) *template.Template {
	return template.Must(template.New("-").Parse(tpl))
}
