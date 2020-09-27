// search.go implements a super simple location search by name.
//
// The main elements are:
// 1. Load TERC and SIMC data.
// 2. Use data to create a tree like structure.
//    root ─┬─ PODKARPACKIE ─┬─ Przemyśl
//          │                └─ Rzeszów
//          └  WARMIŃSKO-MAZURSKIE ─┬─ Elbląg
//                                  ├─ Ełk
//                                  └─ Olsztyn
//
// 3. From the tree like structure it creates a list of all nodes.
// 4. It starts a search loop.
//
// Search is simple process:
// 1. Read text from stdin, strip from whitespace
// 2. Iterate over ALL nodes and compute a score: Levenstein distance between
//    the last component of location and the input text.
// 3. Sort all nodes based on the score.
// 4. Print 10 results with lowest score.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/datainq/teryt"
	"github.com/datainq/teryt/pkg/search"
	"github.com/datainq/teryt/utils"
	counters "github.com/orian/counters/global"
	"github.com/sirupsen/logrus"
)

// Warmińsko-Mazurskie / Olsztyn / Olsztyn
// województwo / powiat / gmina / miejscowosć
func main() {
	log := logrus.New()
	defer counters.WriteTo(os.Stderr)
	counters.LogrusOnSignal()

	var tercFile, simcFile string
	flag.StringVar(&tercFile, "terc", "", "TERC Urzędowy file")
	flag.StringVar(&simcFile, "simc", "", "SIMC Urzędowy file")
	flag.Parse()

	log.Infof("terc: %s", tercFile)
	log.Infof("simc: %s", simcFile)
	tercData, simcData, err := utils.LoadData(tercFile, simcFile)
	if err != nil {
		log.Fatalf("loading data: %s", err)
	}
	root, err := teryt.BuildLocations(tercData, simcData)
	if err != nil {
		log.Fatalf("cannot build locations: %s", err)
	}

	searchNodes := utils.EnlistLocations(root, teryt.Config{Separator: " "})
	counters.Get("search:nodes").IncrementBy(len(searchNodes))
	runSearch(10, search.NewSearchV1(searchNodes))
}

func printResults(results []*search.SearchResult, d time.Duration) {
	fmt.Printf("\r\nResults: (%s)", d)
	for idx, v := range results {
		fmt.Printf("\r\n%d. dist %d: %v (%s)", idx, v.Score,
			v.Location.FullName, v.Location.Type)
	}
}

func runSearch(maxResults int, s search.Interface) {
	fmt.Print("\n>")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Printf("search for %q\n", text)
		t := time.Now()
		results := s.Search(text, maxResults)
		printResults(results, time.Since(t))
		fmt.Print("\n>")
	}
}
