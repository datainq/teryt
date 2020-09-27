package search

import (
	"container/heap"
	"sort"
	"strings"

	"github.com/datainq/teryt"
)

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

type SearchResult struct {
	Location *teryt.Location
	Score    int
}

type Search struct {
	nodes []*LocationWrapper
}

func NewSearch(localities []*teryt.Location) *Search {
	var nodes []*LocationWrapper
	for _, l := range localities {
		nodes = append(nodes, &LocationWrapper{
			Location:   l,
			Score:      0,
			SearchText: []rune(strings.ToLower(l.Parts[len(l.Parts)-1])),
		})
	}

	return &Search{nodes}
}

func (s *Search) Search(text string, limit int) []*SearchResult {
	textRune := []rune(strings.ToLower(strings.TrimSpace(text)))
	maxScore := 100000
	h := &Heap{}
	for i, v := range s.nodes {
		v.Score = levenshtein(v.SearchText, textRune)
		if i < limit {
			heap.Push(h, v)
			maxScore = v.Score
		} else if v.Score < maxScore {
			heap.Push(h, v)
			maxScore = heap.Pop(h).(*LocationWrapper).Score
		}
	}
	locs := *h
	sort.Sort(sort.Reverse(locs))
	ret := make([]*SearchResult, 0, len(locs))
	for _, v := range locs {
		ret = append(ret, &SearchResult{
			Location: v.Location,
			Score:    v.Score,
		})
	}
	return ret
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
