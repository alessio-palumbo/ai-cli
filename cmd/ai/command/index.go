package command

import (
	"ai-cli/internal/indexer"
	"ai-cli/internal/llm"
	"ai-cli/internal/vector"
	"fmt"

	"github.com/urfave/cli/v2"
)

func IndexCommand() *cli.Command {
	return &cli.Command{
		Name:  "index",
		Usage: "index the current repository",

		Action: func(c *cli.Context) error {
			store, err := vector.NewStore()
			if err != nil {
				return err
			}
			defer store.Close()

			if err := store.Clear(); err != nil {
				return err
			}

			if err := indexer.Build(".", store, llm.NewClient()); err != nil {
				return err
			}
			if err := store.Save(); err != nil {
				return err
			}

			fmt.Printf("Indexed %d chunks\n", len(store.Items))
			return nil
		},
	}
}
