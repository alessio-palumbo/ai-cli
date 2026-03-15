package vector

import "math"

func Cosine(a, b []float64) float64 {
	var dot, normA, normB float64
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}
