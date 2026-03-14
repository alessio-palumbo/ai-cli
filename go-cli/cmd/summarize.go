package cmd

import (
	"fmt"
	"os"

	"ai-cli/ai"

	"github.com/urfave/cli/v2"
)

func SummarizeCommand() *cli.Command {
	return &cli.Command{
		Name:  "summarize",
		Usage: "summarize a file",
		Action: func(c *cli.Context) error {
			file := c.Args().First()
			data, err := os.ReadFile(file)
			if err != nil {
				return err
			}

			prompt := fmt.Sprintf(
				"Summarize the following text:\n\n%s",
				string(data),
			)

			err = ai.NewClient().GenerateStream(prompt)
			if err != nil {
				return err
			}

			fmt.Println()
			return nil
		},
	}
}
