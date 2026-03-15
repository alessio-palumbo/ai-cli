package vector

import (
	"cmp"
	"database/sql"
	"os"
	"slices"

	_ "github.com/mattn/go-sqlite3"
)

type Item struct {
	FilePath  string
	Content   string
	Embedding []float64
}

type Store struct {
	db    *sql.DB
	Items []Item
}

func NewStore() (*Store, error) {
	db, err := openDefaultDB()
	if err != nil {
		return nil, err
	}

	s := &Store{db: db}
	if err := s.init(); err != nil {
		db.Close()
		return nil, err
	}
	if err := s.Load(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) Add(chunk, path string, emb []float64) {
	s.Items = append(s.Items, Item{
		FilePath:  path,
		Content:   chunk,
		Embedding: emb,
	})
}

func (s *Store) Search(query []float64, k int) []Item {
	type result struct {
		Item
		Score float64
	}

	var results []result
	for _, item := range s.Items {
		score := Cosine(query, item.Embedding)
		results = append(results, result{
			Item:  item,
			Score: score,
		})
	}

	slices.SortFunc(results, func(i, j result) int {
		return cmp.Compare(j.Score, i.Score)
	})

	var out []Item
	for i := 0; i < k && i < len(results); i++ {
		out = append(out, results[i].Item)
	}

	return out
}

func openDefaultDB() (*sql.DB, error) {
	err := os.MkdirAll(".ai", 0755)
	if err != nil {
		return nil, err
	}

	return sql.Open("sqlite3", ".ai/index.db")
}
