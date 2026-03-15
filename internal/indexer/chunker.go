package indexer

import (
	"path/filepath"
)

const defaultChunkCharacters = 400

func ChunkFile(path string, content string) []string {
	switch filepath.Ext(path) {
	case ".go":
		return ChunkGo(content)
	default:
		return ChunkText(content, 400)
	}
}

func ChunkText(text string, size int) []string {
	var chunks []string
	runes := []rune(text)

	for i := 0; i < len(runes); i += size {
		end := min(i+size, len(runes))
		chunks = append(chunks, string(runes[i:end]))
	}
	return chunks
}
