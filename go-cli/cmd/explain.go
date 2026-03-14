package cmd

import (
	"fmt"
	"os"

	"ai-cli/ai"
	"ai-cli/internal/prompts"

	"github.com/urfave/cli/v2"
)

func ExplainCommand() *cli.Command {
	return &cli.Command{
		Name:  "explain",
		Usage: "explain a source file",
		Action: func(c *cli.Context) error {
			file := c.Args().First()
			data, err := os.ReadFile(file)
			if err != nil {
				return err
			}

			prompt, err := prompts.Render("explain", string(data))
			if err != nil {
				return err
			}
			if err := ai.NewClient().GenerateStream(prompt); err != nil {
				return err
			}

			fmt.Println()
			return nil
		},
	}
}
