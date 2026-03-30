package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApply(t *testing.T) {
	testCases := map[string]struct {
		cfg         Config
		wantErr     bool
		errContains string
		check       func(t *testing.T, cfg Config)
	}{
		"valid config with default extensions": {
			cfg: Config{
				StoreDir:    "/tmp/store",
				ProjectRoot: "/project",
				LLM: struct {
					Model          string
					EmbeddingModel string
					Temperature    float64
				}{
					Model:          "gpt",
					EmbeddingModel: "embed",
					Temperature:    0.5,
				},
			},
			check: func(t *testing.T, cfg Config) {
				for _, ext := range defaultExtensions {
					_, ok := cfg.Extensions[ext]
					assert.True(t, ok, "missing default extension %s", ext)
				}
			},
		},

		"include additional extensions with normalization": {
			cfg: Config{
				StoreDir:    "/tmp/store",
				ProjectRoot: "/project",
				LLM: struct {
					Model          string
					EmbeddingModel string
					Temperature    float64
				}{
					Model:          "gpt",
					EmbeddingModel: "embed",
					Temperature:    0.5,
				},
				Index: struct {
					IncludeExtensions []string
					IgnorePatterns    []string
				}{
					IncludeExtensions: []string{"log", ".CFG", "  txt  "},
				},
			},
			check: func(t *testing.T, cfg Config) {
				expected := []string{".log", ".cfg", ".txt"}
				for _, ext := range expected {
					_, ok := cfg.Extensions[ext]
					assert.True(t, ok, "expected extension %s to be present", ext)
				}
			},
		},

		"error when model missing": {
			cfg: Config{
				StoreDir:    "/tmp/store",
				ProjectRoot: "/project",
				LLM: struct {
					Model          string
					EmbeddingModel string
					Temperature    float64
				}{
					EmbeddingModel: "embed",
					Temperature:    0.5,
				},
			},
			wantErr:     true,
			errContains: "llm model is required",
		},

		"error when embedding model missing": {
			cfg: Config{
				StoreDir:    "/tmp/store",
				ProjectRoot: "/project",
				LLM: struct {
					Model          string
					EmbeddingModel string
					Temperature    float64
				}{
					Model:       "gpt",
					Temperature: 0.5,
				},
			},
			wantErr:     true,
			errContains: "embedding model is required",
		},

		"error when temperature out of range": {
			cfg: Config{
				StoreDir:    "/tmp/store",
				ProjectRoot: "/project",
				LLM: struct {
					Model          string
					EmbeddingModel string
					Temperature    float64
				}{
					Model:          "gpt",
					EmbeddingModel: "embed",
					Temperature:    1.5,
				},
			},
			wantErr:     true,
			errContains: "temperature must be between 0 and 1",
		},

		"error when project root missing": {
			cfg: Config{
				StoreDir: "/tmp/store",
				LLM: struct {
					Model          string
					EmbeddingModel string
					Temperature    float64
				}{
					Model:          "gpt",
					EmbeddingModel: "embed",
					Temperature:    0.5,
				},
			},
			wantErr:     true,
			errContains: "project root is required",
		},

		"error when store dir missing": {
			cfg: Config{
				ProjectRoot: "/project",
				LLM: struct {
					Model          string
					EmbeddingModel string
					Temperature    float64
				}{
					Model:          "gpt",
					EmbeddingModel: "embed",
					Temperature:    0.5,
				},
			},
			wantErr:     true,
			errContains: "store directory is required",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			cfg := tc.cfg
			err := cfg.Apply()

			if tc.wantErr {
				assert.Error(t, err)
				if tc.errContains != "" {
					assert.Contains(t, err.Error(), tc.errContains)
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, cfg.Extensions)

			if tc.check != nil {
				tc.check(t, cfg)
			}
		})
	}
}
