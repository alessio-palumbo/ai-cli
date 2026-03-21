package indexer

import (
	"fmt"
	"path/filepath"
	"strings"
)

const (
	defaultChunkCharacters = 400
	minChunkSize           = 50
	chunkPathPrefix        = "file: "
)

type Chunk struct {
	Text      string
	StartLine int
	EndLine   int
}

func ChunkFile(path string, content string) []Chunk {
	switch filepath.Ext(path) {
	case ".go":
		return ChunkGo(path, content)
	default:
		return ChunkText(path, content)
	}
}

// ChunkText splits arbitrary text into fixed-size chunks.
// The file path and line numbers are included to preserve context for embeddings.
func ChunkText(path string, content string) []Chunk {
	var chunks []Chunk
	startByte := 0
	startLine := 1
	charCount := 0
	line := 1

	for i := range len(content) {
		charCount++
		if content[i] == '\n' {
			line++
		}

		if charCount >= defaultChunkCharacters {
			endByte := i + 1
			// trim trailing newline
			if endByte < len(content) && content[endByte-1] == '\n' {
				endByte--
			}

			body := content[startByte:endByte]
			chunks = append(chunks, Chunk{
				StartLine: startLine,
				EndLine:   line,
				Text:      formatChunk(path, startLine, line, body),
			})

			startByte = endByte
			startLine = line
			charCount = 0
		}
	}

	if startByte < len(content) {
		body := content[startByte:]
		if len(body) >= minChunkSize {
			chunks = append(chunks, Chunk{
				StartLine: startLine,
				EndLine:   line,
				Text:      formatChunk(path, startLine, line, body),
			})
		}
	}

	return chunks
}

func formatChunk(path string, startLine, endLine int, text ...string) string {
	var sb strings.Builder
	lines := fmt.Sprintf("(lines %d-%d)", startLine, endLine)
	totalSize := len(chunkPathPrefix) + len(path) + 1
	totalSize += len(lines) + 1
	for _, t := range text {
		totalSize += len(t) + 1
	}
	sb.Grow(totalSize)

	sb.WriteString(chunkPathPrefix)
	sb.WriteString(path)
	sb.WriteString(" ")
	sb.WriteString(lines)
	sb.WriteByte('\n')

	for _, t := range text {
		sb.WriteByte('\n')
		sb.WriteString(t)
	}
	return sb.String()
}
