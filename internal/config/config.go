package config

import (
	"fmt"
	"strings"
)

var defaultExtensions = []string{
	// core code
	".go", ".py", ".js", ".ts", ".rb", ".rs",

	// configs (HIGH VALUE in RAG)
	".json", ".yaml", ".yml", ".toml",

	// scripting
	".sh", ".bash", ".zsh",

	// docs
	".md", ".txt",

	// infra
	".tf", ".tfvars",
}

type Config struct {
	StoreDir    string
	ProjectRoot string
	DBName      string
	Extensions  map[string]struct{}

	LLM struct {
		Model          string
		EmbeddingModel string
		Temperature    float64
	}

	Index struct {
		IncludeExtensions []string
		IgnorePatterns    []string
	}
}

func (c *Config) Apply() error {
	c.resolveExtensions()
	return c.validate()
}

func (c *Config) resolveExtensions() {
	c.Extensions = make(map[string]struct{}, len(defaultExtensions))
	for _, e := range defaultExtensions {
		c.Extensions[e] = struct{}{}
	}
	for _, e := range c.Index.IncludeExtensions {
		if ext := normalizeExt(e); ext != "" {
			c.Extensions[ext] = struct{}{}
		}
	}
}

func (c *Config) validate() error {
	if c.LLM.Model == "" {
		return fmt.Errorf("llm model is required")
	}
	if c.LLM.EmbeddingModel == "" {
		return fmt.Errorf("embedding model is required")
	}
	if c.LLM.Temperature < 0 || c.LLM.Temperature > 1 {
		return fmt.Errorf("temperature must be between 0 and 1")
	}
	if c.ProjectRoot == "" {
		return fmt.Errorf("project root is required")
	}
	if c.StoreDir == "" {
		return fmt.Errorf("store directory is required to store vector DBs")
	}
	return nil
}

func normalizeExt(e string) string {
	e = strings.TrimSpace(e)
	if e == "" {
		return ""
	}

	if !strings.HasPrefix(e, ".") {
		e = "." + e
	}
	return strings.ToLower(e)
}
