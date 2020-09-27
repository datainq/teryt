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
	"archive/zip"
	"bufio"
	"container/heap"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/datainq/rwmc"
	"github.com/datainq/teryt"
	counters "github.com/orian/counters/global"
	"github.com/sirupsen/logrus"
)

// Warmińsko-Mazurskie / Olsztyn / Olsztyn
// województwo / powiat / gmina / miejscowosć
func main() {
	log := logrus.New()
	defer counters.WriteTo(os.Stderr)
	counters.LogrusOnSignal()

	tercFile := flag.String("terc", "", "TERC Urzędowy file")
	simcFile := flag.String("simc", "", "SIMC Urzędowy file")
	flag.Parse()

	log.Infof("terc: %s", *tercFile)
	log.Infof("simc: %s", *simcFile)
	tercData, simcData, err := loadData(*tercFile, *simcFile)
	if err != nil {
		log.Fatalf("loading data: %s", err)
	}
	root, err := teryt.BuildLocations(tercData, simcData)
	if err != nil {
		log.Fatalf("cannot build locations: %s", err)
	}

	searchNodes := enlistLocations(root, teryt.Config{Separator: " "})
	counters.Get("search:nodes").IncrementBy(len(searchNodes))
	runSearch(searchNodes, 10)
	// fmt.Printf("number of nodes: %d\n", printNodes(root))
}

func Open(filePath string) (io.ReadCloser, error) {
	var dataReader io.ReadCloser
	if strings.HasSuffix(filePath, ".zip") {
		r, err := zip.OpenReader(filePath)
		if err != nil {
			return nil, fmt.Errorf("cannot open: %w", err)
		}
		fileName := path.Base(filePath[:len(filePath)-4]) + ".csv"
		for _, f := range r.File {
			if strings.HasPrefix(f.Name, fileName) {
				dataReader, err = f.Open()
				if err != nil {
					_ = r.Close()
					return nil, fmt.Errorf("cannot open file: %filePath", err)
				}
				return rwmc.NewReadMultiCloser(dataReader, r), nil
			}
		}
		_ = r.Close()
	} else {
		return os.Open(filePath)
	}
	return nil, fmt.Errorf("file %q not found", filePath)
}

func loadData(tercFile, simcFile string) ([]teryt.SetTERC, []teryt.SetSIMC, error) {
	tercReader, err := Open(tercFile)
	if err != nil {
		return nil, nil, fmt.Errorf("open TERC file: %w", err)
	}
	defer tercReader.Close()
	tercData, err := teryt.TERCFromReader(tercReader)
	if err != nil {
		return nil, nil, fmt.Errorf("reading TERC: %w", err)
	}

	simcReader, err := Open(simcFile)
	if err != nil {
		return nil, nil, fmt.Errorf("open SIMC file: %w", err)
	}

	simcData, err := teryt.SIMCFromReader(simcReader)
	if err != nil {
		return nil, nil, fmt.Errorf("reading SIMC: %w", err)
	}
	return tercData, simcData, nil
}

func enlistLocations(root *teryt.Location, cfg teryt.Config) []*teryt.Location {
	var ret []*teryt.Location
	if root.Name != "" {
		root.Build(cfg)
		ret = append(ret, root)
	}
	if len(root.Children) > 1 {
		for _, v := range root.Children {
			ret = append(ret, enlistLocations(v, cfg)...)
		}
	}
	return ret
}

type LocationWrapper struct {
	*teryt.Location
	Score      int
	SearchText []rune
}

type Heap []*LocationWrapper

func (h Heap) Len() int {
	return len(h)
}

func (h Heap) Less(i, j int) bool {
	a, b := h[i], h[j]
	return a.Score > b.Score || a.Score == b.Score && (len(a.Parts) > len(b.Parts))
}

func (h Heap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *Heap) Push(x interface{}) {
	*h = append(*h, x.(*LocationWrapper))
}

func (h *Heap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*h = old[0 : n-1]
	return item
}

func runSearch(localities []*teryt.Location, maxResults int) {
	fmt.Print(">")
	scanner := bufio.NewScanner(os.Stdin)

	var nodes []*LocationWrapper
	for _, l := range localities {
		nodes = append(nodes, &LocationWrapper{
			Location:   l,
			Score:      0,
			SearchText: []rune(strings.ToLower(l.Parts[len(l.Parts)-1])),
		})
	}

	for scanner.Scan() {
		t := time.Now()
		text := strings.ToLower(strings.TrimSpace(scanner.Text()))
		fmt.Printf("search for %q\n", text)
		textRune := []rune(text)
		maxScore := 100000
		h := &Heap{}
		for i, v := range nodes {
			v.Score = levenshtein(v.SearchText, textRune)
			if i < maxResults {
				heap.Push(h, v)
				maxScore = v.Score
			} else if v.Score < maxScore {
				heap.Push(h, v)
				maxScore = heap.Pop(h).(*LocationWrapper).Score
			}
		}
		locs := *h
		sort.Sort(sort.Reverse(locs))

		fmt.Printf("Results: (%s)\n", time.Since(t))
		for idx, v := range locs {
			fmt.Printf("%d. dist %d: %v (%s)\n", idx, v.Score, v.FullName, v.Type)
		}
		fmt.Print(">")
	}
}

func levenshtein(str1, str2 []rune) int {
	s1len := len(str1)
	s2len := len(str2)
	column := make([]int, len(str1)+1)

	for y := 1; y <= s1len; y++ {
		column[y] = y
	}
	for x := 1; x <= s2len; x++ {
		column[0] = x
		lastkey := x - 1
		for y := 1; y <= s1len; y++ {
			oldkey := column[y]
			var incr int
			if str1[y-1] != str2[x-1] {
				incr = 1
			}

			column[y] = minimum3(column[y]+1, column[y-1]+1, lastkey+incr)
			lastkey = oldkey
		}
	}
	return column[s1len]
}

func minimum3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
	} else {
		if b < c {
			return b
		}
	}
	return c
}