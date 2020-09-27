package search

import (
	"fmt"
	"testing"

	"github.com/datainq/teryt"
	"github.com/datainq/teryt/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Fatalf interface {
	Fatalf(format string, args ...interface{})
}

func mustLoadNodes(t Fatalf) []*teryt.Location {
	tercData, simcData, err := utils.LoadData("../data/TERC_Urzedowy_2020-09-26.zip",
		"../data/SIMC_Urzedowy_2020-09-26.zip")
	if err != nil {
		t.Fatalf("loading data: %s", err)
	}
	root, err := teryt.BuildLocations(tercData, simcData)
	if err != nil {
		t.Fatalf("cannot build locations: %s", err)
	}
	return utils.EnlistLocations(root, teryt.Config{Separator: " "})
}

func BenchmarkSearch_SearchOld(b *testing.B) {
	s := NewSearchOld(mustLoadNodes(b))
	var names = []string{
		"Olsztyn", "Ełk", "Elbląg", "Giżycko", "Wrocław", "Berlin", "Reszel",
		"Kętrzyn", "Olsztynek", "Rzeszów", "Łódź", "Pupki", "Jonkowo", "Warkały",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val := s.Search(names[i%len(names)], 10)
		_ = val
	}
}

func BenchmarkSearch_Search(b *testing.B) {
	s := NewSearch(mustLoadNodes(b))
	var names = []string{
		"Olsztyn", "Ełk", "Elbląg", "Giżycko", "Wrocław", "Berlin", "Reszel",
		"Kętrzyn", "Olsztynek", "Rzeszów", "Łódź", "Pupki", "Jonkowo", "Warkały",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val := s.Search(names[i%len(names)], 10)
		_ = val
	}
}

func BenchmarkSearch_SearchV3(b *testing.B) {
	var names = []string{
		"Olsztyn", "Ełk", "Elbląg", "Giżycko", "Wrocław", "Berlin", "Reszel",
		"Kętrzyn", "Olsztynek", "Rzeszów", "Łódź", "Pupki", "Jonkowo", "Warkały",
	}
	for i := 1; i < 9; i++ {
		b.Run(fmt.Sprintf("parallel-%d", i), func(b *testing.B) {
			s := NewSearchV3(mustLoadNodes(b), i)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				val := s.Search(names[i%len(names)], 10)
				_ = val
			}
		})
	}
}

func TestSearch_SearchTop(t *testing.T) {
	s := NewSearch(mustLoadNodes(t))

	for _, text := range []string{
		"Olsztyn", "Ełk", "Elbląg", "Giżycko", "Wrocław", "Reszel",
		"Kętrzyn", "Olsztynek", "Rzeszów", "Łódź", "Pupki", "Jonkowo", "Warkały",
	} {
		res := s.Search(text, 1)
		require.Len(t, res, 1)
		r := res[0]
		assert.Equal(t, r.Location.Name, text)
	}
}

func TestSearch_Search3Top(t *testing.T) {
	s := NewSearchV3(mustLoadNodes(t), 4)

	for _, text := range []string{
		"Olsztyn", "Ełk", "Elbląg", "Giżycko", "Wrocław", "Reszel",
		"Kętrzyn", "Olsztynek", "Rzeszów", "Łódź", "Pupki", "Jonkowo", "Warkały",
	} {
		res := s.Search(text, 1)
		require.Len(t, res, 1)
		r := res[0]
		assert.Equal(t, r.Location.Name, text)
	}
}
