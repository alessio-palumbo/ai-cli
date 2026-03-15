package indexer

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func ChunkGo(src string) []string {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return ChunkText(src, 400)
	}

	var chunks []string

	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		start := fset.Position(fn.Pos()).Offset
		end := fset.Position(fn.End()).Offset
		chunks = append(chunks, src[start:end])
	}

	return chunks
}
