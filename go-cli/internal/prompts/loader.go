package prompts

import (
	"bytes"
	"embed"
	"text/template"
)

//go:embed templates/*.tmpl
var promptFS embed.FS

var templates = template.Must(
	template.ParseFS(promptFS, "templates/*.tmpl"),
)

type Content struct {
	Prompt string
}

func Render(templateName string, prompt string) (string, error) {
	var buf bytes.Buffer
	t := templates.Lookup(templateName + ".tmpl")
	if err := t.Execute(&buf, Content{Prompt: prompt}); err != nil {
		return "", err
	}

	return buf.String(), nil
}
