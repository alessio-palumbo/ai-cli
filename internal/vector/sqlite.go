package vector

import (
	"bytes"
	"encoding/binary"
)

func (s *Store) Save() error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	stmt, err := tx.Prepare(`
	    INSERT INTO embeddings(content, embedding)
	    VALUES(?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, item := range s.Items {
		blob, err := encodeEmbedding(item.Embedding)
		if err != nil {
			return err
		}

		_, err = stmt.Exec(item.Content, blob)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *Store) Load() ([]Item, error) {
	rows, err := s.db.Query(`
	    SELECT content, embedding
	    FROM embeddings
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Item

	for rows.Next() {
		var (
			content string
			blob    []byte
		)

		err := rows.Scan(&content, &blob)
		if err != nil {
			return nil, err
		}

		vec, err := decodeEmbedding(blob)
		if err != nil {
			return nil, err
		}

		items = append(items, Item{
			Content:   content,
			Embedding: vec,
		})
	}

	return items, nil
}

func encodeEmbedding(vec []float64) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, vec)
	return buf.Bytes(), err
}

func decodeEmbedding(data []byte) ([]float64, error) {
	count := len(data) / 8
	vec := make([]float64, count)

	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, &vec)
	return vec, err
}
