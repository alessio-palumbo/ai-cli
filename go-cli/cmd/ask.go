package cmd

import (
	"fmt"
	"strings"

	"ai-cli/ai"

	"github.com/urfave/cli/v2"
)

func AskCommand() *cli.Command {
	return &cli.Command{
		Name:  "ask",
		Usage: "ask the AI a question",
		Action: func(c *cli.Context) error {
			prompt := strings.Join(c.Args().Slice(), " ")
			err := ai.NewClient().GenerateStream(prompt)
			if err != nil {
				return err
			}

			fmt.Println()
			return nil
		},
	}
}
