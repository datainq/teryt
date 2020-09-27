package search

import (
	"testing"

	"github.com/datainq/teryt"
	"github.com/datainq/teryt/utils"
)

func BenchmarkSearch_Search(b *testing.B) {
	tercData, simcData, err := utils.LoadData("../data/TERC_Urzedowy_2020-09-26.zip",
		"../data/SIMC_Urzedowy_2020-09-26.zip")
	if err != nil {
		b.Fatalf("loading data: %s", err)
	}
	root, err := teryt.BuildLocations(tercData, simcData)
	if err != nil {
		b.Fatalf("cannot build locations: %s", err)
	}
	searchNodes := utils.EnlistLocations(root, teryt.Config{Separator: " "})
	s := NewSearch(searchNodes)
	var names = []string{
		"Olsztyn", "Ełk", "Elbląg", "Giżycko", "Wrocław", "Berlin", "Reszel",
		"Kętrzyn", "Olsztynek", "Rzeszów", "Łódź", "Pupki", "Jonkowo", "Warkały",
	}
	for i := 0; i < b.N; i++ {
		val := s.Search(names[i%len(names)], 10)
		_ = val
	}
}
