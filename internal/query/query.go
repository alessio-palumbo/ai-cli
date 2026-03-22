package query

import (
	"regexp"
	"strings"
)

var (
	reIdentifier = regexp.MustCompile(`\b([a-zA-Z]+[A-Z][a-zA-Z0-9]*|[a-z][a-z0-9]*_[a-z][a-z0-9_]*)\b`)
	reFileRef    = regexp.MustCompile(`\b[\w\-]+\.(go|ts|js|py|rs|java|cpp|c|rb|cs|json|yaml|yml|toml|md)\b`)
	rePkgPath    = regexp.MustCompile(`[\w\-]+/[\w\-]+(/[\w\-]+)*`)
	reErrFrag    = regexp.MustCompile(`(?i)(panic:|error:|undefined:|cannot|failed to|nil pointer|index out of|no such file|unexpected token|syntax error)`)
)

var codeWords = []string{
	"function", "method", "struct", "interface",
	"command", "package", "module", "var",
	"handler", "client", "server", "config", "log",
}

// IsSearchableQuery determines whether a user query contains enough
// concrete signals (e.g. identifiers, file names, error messages, or
// code-related terms) to safely perform vector search over the indexed codebase.
//
// If it returns false, the query is considered too vague or exploratory,
// and additional context (e.g. repository summary or query expansion)
// should be used before attempting retrieval.
func IsSearchable(q string) bool {
	q = strings.TrimSpace(q)
	if q == "" {
		return false
	}

	return hasAnchor(q)
}

func hasAnchor(q string) bool {
	return reIdentifier.MatchString(q) ||
		reFileRef.MatchString(q) ||
		rePkgPath.MatchString(q) ||
		reErrFrag.MatchString(q) ||
		hasCodeWord(q)
}

func hasCodeWord(q string) bool {
	q = strings.ToLower(q)
	for _, w := range codeWords {
		if strings.Contains(q, w) {
			return true
		}
	}
	return false
}
