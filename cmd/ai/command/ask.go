package command

import (
	"ai-cli/internal/llm"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func AskCommand() *cli.Command {
	return &cli.Command{
		Name:  "ask",
		Usage: "ask the AI a question",
		Action: func(c *cli.Context) error {
			prompt := strings.Join(c.Args().Slice(), " ")
			if err := llm.NewClient().GenerateStream(prompt, os.Stdout); err != nil {
				return err
			}

			fmt.Println()
			return nil
		},
	}
}
