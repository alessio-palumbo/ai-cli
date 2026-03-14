package main

import (
	"ai-cli/cmd"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "ai",
		Usage: "AI powered developer assistant",
		Commands: []*cli.Command{
			cmd.AskCommand(),
			cmd.SummarizeCommand(),
			cmd.ExplainCommand(),
		},
	}

	app.Run(os.Args)
}
