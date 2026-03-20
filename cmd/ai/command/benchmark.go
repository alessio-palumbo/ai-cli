package command

import (
	"ai-cli/internal/benchmark"
	"ai-cli/internal/llm"
	"ai-cli/internal/vector"

	"github.com/urfave/cli/v2"
)

func BenchmarkCommand(llmClient *llm.Client, store *vector.Store) *cli.Command {
	return &cli.Command{
		Name:  "benchmark",
		Usage: "benchmark harness to systematically evaluate changes",
		Action: func(c *cli.Context) error {
			benchmark.Run(store, llmClient.Embed)
			return nil
		},
	}
}
