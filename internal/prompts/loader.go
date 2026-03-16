package prompts

import (
	"ai-cli/internal/vector"
	"bytes"
	"embed"
	"text/template"
)

type promptTemplate string

const (
	TemplateExplain   promptTemplate = "explain.tmpl"
	TemplateSummarize promptTemplate = "summarize.tmpl"
	TemplateQuery     promptTemplate = "query.tmpl"
)

//go:embed templates/*.tmpl
var promptFS embed.FS

var templates = template.Must(
	template.ParseFS(promptFS, "templates/*.tmpl"),
)

type Content struct {
	Prompt  string
	Context string
}

func Render(tmpl promptTemplate, prompt string, context ...vector.Result) (string, error) {
	content := Content{Prompt: prompt}
	for _, r := range context {
		content.Context += r.Content + "\n---\n"
	}

	var buf bytes.Buffer
	if err := templates.Lookup(string(tmpl)).
		Execute(&buf, content); err != nil {
		return "", err
	}

	return buf.String(), nil
}
