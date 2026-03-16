package command

import (
	"ai-cli/internal/llm"
	"ai-cli/internal/prompts"
	"ai-cli/internal/vector"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func QueryCommand() *cli.Command {
	return &cli.Command{
		Name:  "query",
		Usage: "ask a question over your indexed code",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "k",
				Value: 5,
				Usage: "number of top chunks to retrieve",
			},
		},
		Action: func(c *cli.Context) error {
			query := strings.Join(c.Args().Slice(), " ")
			if query == "" {
				return fmt.Errorf("query required")
			}

			store, err := vector.NewStore()
			if err != nil {
				return err
			}
			defer store.Close()

			client := llm.NewClient()
			queryVec, err := client.Embed(query)
			if err != nil {
				return err
			}

			results := store.Search(queryVec, c.Int("k"))
			if len(results) == 0 {
				fmt.Println("No relevant results found")
				return nil
			}

			prompt, err := prompts.Render(prompts.TemplateQuery, query, results...)
			if err != nil {
				return err
			}

			if err := client.GenerateStream(prompt, os.Stdout); err != nil {
				return err
			}

			fmt.Println()
			return nil
		},
	}
}
