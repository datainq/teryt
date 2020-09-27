package search

import (
	"testing"

	"github.com/datainq/teryt"
	"github.com/datainq/teryt/utils"
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
	s := NewSearchV1(mustLoadNodes(b))
	var names = []string{
		"Olsztyn", "Ełk", "Elbląg", "Giżycko", "Wrocław", "Berlin", "Reszel",
		"Kętrzyn", "Olsztynek", "Rzeszów", "Łódź", "Pupki", "Jonkowo", "Warkały",
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val := s.Search(names[i%len(names)], 10)
		result = val
	}
}

var result []*SearchResult
