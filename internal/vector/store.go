package vector

type Item struct {
	Text      string
	Embedding []float64
}

type Store struct {
	Items []Item
}

func (s *Store) Add(text string, emb []float64) {
	s.Items = append(s.Items, Item{
		Text:      text,
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

	// simple top-k selection
	for i := range results {
		for j := i + 1; j < len(results); j++ {
			if results[j].Score > results[i].Score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	var out []Item
	for i := 0; i < k && i < len(results); i++ {
		out = append(out, results[i].Item)
	}

	return out
}
