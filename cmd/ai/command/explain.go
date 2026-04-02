package command

import (
	"ai-cli/pkg/ai"
	"ai-cli/pkg/spinner"
	"fmt"

	"github.com/urfave/cli/v2"
)

func ExplainCommand(client *ai.Client, sw *spinner.StreamWriter) *cli.Command {
	return &cli.Command{
		Name:  "explain",
		Usage: "explain a source file",
		Action: func(c *cli.Context) error {
			file := c.Args().First()

			result, err := spinner.Wrap(sw, func() (ai.TaskResult, error) {
				return client.Explain(c.Context, file)
			})
			if err != nil {
				return catchIndexError(err)
			}

			if result.Status.NoResults {
				fmt.Println("No relevant results found")
			}

			fmt.Println()
			return nil
		},
	}
}
