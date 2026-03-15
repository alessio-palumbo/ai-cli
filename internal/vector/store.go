package vector

import (
	"cmp"
	"database/sql"
	"os"
	"slices"

	_ "github.com/mattn/go-sqlite3"
)

type Item struct {
	Content   string
	FilePath  string
	Embedding []float64
}

type Store struct {
	db    *sql.DB
	Items []Item
}

func OpenDefaultDB() (*sql.DB, error) {
	err := os.MkdirAll(".ai", 0755)
	if err != nil {
		return nil, err
	}

	return sql.Open("sqlite3", ".ai/index.db")
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Init() error {
	query := `
	    CREATE TABLE IF NOT EXISTS embeddings (
	        id INTEGER PRIMARY KEY,
	        content TEXT,
	        filepath TEXT,
	        embedding BLOB
	    );
	`

	_, err := s.db.Exec(query)
	return err
}

func (s *Store) Add(chunk, path string, emb []float64) {
	s.Items = append(s.Items, Item{
		Content:   chunk,
		FilePath:  path,
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
		return cmp.Compare(i.Score, j.Score)
	})

	var out []Item
	for i := 0; i < k && i < len(results); i++ {
		out = append(out, results[i].Item)
	}

	return out
}
