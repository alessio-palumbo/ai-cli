package prompts

import (
	"bytes"
	"embed"
	"text/template"
)

type promptTemplate string

const (
	TemplateExplain   promptTemplate = "explain.tmpl"
	TemplateSummarize promptTemplate = "summarize.tmpl"
)

//go:embed templates/*.tmpl
var promptFS embed.FS

var templates = template.Must(
	template.ParseFS(promptFS, "templates/*.tmpl"),
)

type Content struct {
	Prompt string
}

func Render(tmpl promptTemplate, prompt string) (string, error) {
	var buf bytes.Buffer
	if err := templates.Lookup(string(tmpl)).
		Execute(&buf, Content{Prompt: prompt}); err != nil {
		return "", err
	}

	return buf.String(), nil
}
