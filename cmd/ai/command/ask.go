package command

import (
	"ai-cli/pkg/ai"
	"ai-cli/pkg/spinner"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

func AskCommand(client *ai.Client, sw *spinner.StreamWriter) *cli.Command {
	return &cli.Command{
		Name:  "ask",
		Usage: "ask the AI a question",
		Action: func(c *cli.Context) error {
			prompt := strings.Join(c.Args().Slice(), " ")

			return spinner.WrapError(sw, func() error {
				if err := client.Ask(c.Context, prompt); err != nil {
					return catchIndexError(err)
				}
				fmt.Println()
				return nil
			})
		},
	}
}
