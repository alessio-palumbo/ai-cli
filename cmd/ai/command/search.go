package command

import (
	"ai-cli/internal/llm"
	"ai-cli/internal/vector"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

func SearchCommand(llmClient *llm.Client, store *vector.Store) *cli.Command {
	return &cli.Command{
		Name:  "search",
		Usage: "semantic search in indexed code",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "k",
				Value: 5,
				Usage: "number of results",
			},
			&cli.BoolFlag{
				Name:  "mmr",
				Value: false,
				Usage: "use Max Marginal Relevance",
			},
		},
		Action: func(c *cli.Context) error {
			prompt := strings.Join(c.Args().Slice(), " ")
			if prompt == "" {
				return fmt.Errorf("query required")
			}

			queryVec, err := llmClient.Embed(prompt)
			if err != nil {
				return err
			}

			results, err := store.Search(queryVec, c.Int("k"), c.Bool("mmr"))
			if err != nil {
				return err
			}

			fmt.Println(vector.JoinResults(results...))
			return nil
		},
	}
}
