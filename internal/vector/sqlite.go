package vector

import (
	"bytes"
	"encoding/binary"
)

func (s *Store) Clear() error {
	_, err := s.db.Exec(`DELETE FROM embeddings`)
	if err != nil {
		return err
	}

	s.Items = nil
	return nil
}

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
	    INSERT INTO embeddings(filepath, content, embedding)
	    VALUES(?, ?, ?)
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

		_, err = stmt.Exec(item.FilePath, item.Content, blob)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *Store) Load() error {
	rows, err := s.db.Query(`
	    SELECT filepath, content, embedding
	    FROM embeddings
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var items []Item

	for rows.Next() {
		var (
			filepath string
			content  string
			blob     []byte
		)

		err := rows.Scan(&filepath, &content, &blob)
		if err != nil {
			return err
		}

		vec, err := decodeEmbedding(blob)
		if err != nil {
			return err
		}

		items = append(items, Item{
			FilePath:  filepath,
			Content:   content,
			Embedding: vec,
		})
	}

	s.Items = items
	return nil
}

func (s *Store) init() error {
	query := `
	    CREATE TABLE IF NOT EXISTS embeddings (
	        id INTEGER PRIMARY KEY,
	        filepath TEXT,
	        content TEXT,
	        embedding BLOB
	    );
	`

	_, err := s.db.Exec(query)
	return err
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
