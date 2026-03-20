package benchmark

import (
	"ai-cli/internal/vector"
	"fmt"
	"log"
	"strings"
)

var tests = []TestCase{
	{
		Name:  "ask command discovery",
		Query: "how to ask a question about indexed code",
		ExpectedFiles: []string{
			"ask.go",
			"query.go",
			"search.go",
		},
	},
	{
		Name:  "indexing flow",
		Query: "how indexing works",
		ExpectedFiles: []string{
			"index.go",
			"chunk.go",
		},
	},
}

type TestCase struct {
	Name          string
	Query         string
	ExpectedFiles []string
}

func Run(store *vector.Store, embed func(string) ([]float64, error)) {
	for _, tc := range tests {
		fmt.Println("Test:", tc.Name)

		for _, mode := range []string{vector.SearchModeFast, vector.SearchModeBalanced, vector.SearchModeDeep} {
			queryVec, err := embed(tc.Query)
			if err != nil {
				log.Fatal(err)
			}
			results, err := store.SearchForMode(mode, queryVec)
			if err != nil {
				log.Fatal(err)
			}
			s := score(results, tc.ExpectedFiles)
			fmt.Printf(" %s: %.2f\n", mode, s)
			for _, r := range results {
				fmt.Println(" -", r.FilePath)
			}
			fmt.Println()
		}
	}
}

func score(results []vector.Result, expected []string) float64 {
	hits := 0
	for _, e := range expected {
		for _, r := range results {
			if strings.Contains(r.FilePath, e) {
				hits++
				break
			}
		}
	}
	return float64(hits) / float64(len(expected))
}
