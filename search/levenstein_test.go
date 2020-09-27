package search

import (
	"testing"

	agnivade "github.com/agnivade/levenshtein"
	arbovm "github.com/arbovm/levenshtein"
	dgryski "github.com/dgryski/trifles/leven"
	m1ome "github.com/m1ome/leven"
)

func TestLevenshteinSanity(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"", "hello", 5},
		{"hello", "", 5},
		{"hello", "hello", 0},
		{"ab", "aa", 1},
		{"ab", "ba", 2},
		{"ab", "aaa", 2},
		{"bbb", "a", 3},
		{"kitten", "sitting", 3},
		{"distance", "difference", 5},
		{"levenshtein", "frankenstein", 6},
		{"resume and cafe", "resumes and cafes", 2},
	}
	for i, d := range tests {
		n := LevenshteinRune([]rune(d.a), []rune(d.b))
		if n != d.want {
			t.Errorf("Test[%d]: ComputeDistance(%q,%q) returned %v, want %v",
				i, d.a, d.b, n, d.want)
		}
	}
}

func TestLevenshteinStrSanity(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"", "hello", 5},
		{"hello", "", 5},
		{"hello", "hello", 0},
		{"ab", "aa", 1},
		{"ab", "ba", 2},
		{"ab", "aaa", 2},
		{"bbb", "a", 3},
		{"kitten", "sitting", 3},
		{"distance", "difference", 5},
		{"levenshtein", "frankenstein", 6},
		{"resume and cafe", "resumes and cafes", 2},
	}
	for i, d := range tests {
		n := LevenshteinStr(d.a, d.b)
		if n != d.want {
			t.Errorf("Test[%d]: ComputeDistance(%q,%q) returned %v, want %v",
				i, d.a, d.b, n, d.want)
		}
	}
}

func TestUnicode(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		// Testing acutes and umlauts
		{"resumé and café", "resumés and cafés", 2},
		{"resume and cafe", "resumé and café", 2},
		{"Hafþór Júlíus Björnsson", "Hafþor Julius Bjornsson", 4},
		// Only 2 characters are less in the 2nd string
		{"།་གམ་འས་པ་་མ།", "།་གམའས་པ་་མ", 2},
	}
	for i, d := range tests {
		n := LevenshteinRune([]rune(d.a), []rune(d.b))
		if n != d.want {
			t.Errorf("Test[%d]: ComputeDistance(%q,%q) returned %v, want %v",
				i, d.a, d.b, n, d.want)
		}
	}
}

// Benchmarks
// ----------------------------------------------
var sink int

func BenchmarkAll(b *testing.B) {
	tests := []struct {
		a, b string
		name string
	}{
		// ASCII
		{"LevenshteinRune", "frankenstein", "ASCII"},
		{"ŁÓźęĄ ŚŹła", "Żółć ąęŹÓń", "Polish"},
		// Testing acutes and umlauts
		{"resumé and café", "resumés and cafés", "French"},
		{"Hafþór Júlíus Björnsson", "Hafþor Julius Bjornsson", "Nordic"},
		// Only 2 characters are less in the 2nd string
		{"།་གམ་འས་པ་་མ།", "།་གམའས་པ་་མ", "Tibetan"},
	}
	tmp := 0
	b.ReportAllocs()
	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			b.Run("agniva", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					tmp = agnivade.ComputeDistance(test.a, test.b)
				}
			})
			b.Run("m1ome", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					tmp = m1ome.Distance(test.a, test.b)
				}
			})
			b.Run("arbovm", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					tmp = arbovm.Distance(test.a, test.b)
				}
			})
			b.Run("dgryski", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					tmp = dgryski.Levenshtein([]rune(test.a), []rune(test.b))
				}
			})
			b.Run("datainq", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					tmp = LevenshteinRune([]rune(test.a), []rune(test.b))
				}
			})
			b.Run("datainq2", func(b *testing.B) {
				l := &Levenshteiner{}
				b.ResetTimer()
				for n := 0; n < b.N; n++ {
					tmp = l.LevenshteinRune([]rune(test.a), []rune(test.b))
				}
			})
			// b.Run("datainq2", func(b *testing.B) {
			// 	for n := 0; n < b.N; n++ {
			// 		tmp = LevenshteinStr(test.a, test.b)
			// 	}
			// })
		})
	}
	sink = tmp
}
