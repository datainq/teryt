package search

import (
	"container/heap"
	"strings"
	"sync"

	"github.com/datainq/teryt"
)

type SearchV3 struct {
	nodes       []*LocationWrapper
	chunks      [][]*LocationWrapper
	parallelism int
}

func NewSearchV3(localities []*teryt.Location, parallelism int) *SearchV3 {
	var nodes []*LocationWrapper
	for _, l := range localities {
		nodes = append(nodes, &LocationWrapper{
			Location:   l,
			Score:      0,
			SearchText: []rune(strings.ToLower(l.Parts[len(l.Parts)-1])),
		})
	}

	var chunks [][]*LocationWrapper
	for i := 0; i < parallelism; i++ {
		start := i * len(nodes) / parallelism
		end := (i + 1) * len(nodes) / parallelism
		if end > len(nodes) {
			end = len(nodes)
		}
		chunks = append(chunks, nodes[start:end])
	}
	return &SearchV3{nodes, chunks, parallelism}
}

func (s *SearchV3) search(idx int, wg *sync.WaitGroup, text []rune, limit int, result []*LocationWrapper) {
	maxScore := 100000
	h := &Heap{}
	for _, v := range s.chunks[idx][:limit] {
		v.Score = LevenshteinRune(v.SearchText, text)
		heap.Push(h, v)
		maxScore = max(maxScore, v.Score)
	}
	for _, v := range s.chunks[idx][limit:] {
		v.Score = LevenshteinRune(v.SearchText, text)
		if v.Score <= maxScore {
			heap.Push(h, v)
			maxScore = heap.Pop(h).(*LocationWrapper).Score
		}
	}
	for i, v := range *h {
		result[i] = v
	}
	wg.Done()
}

func (s *SearchV3) Search(text string, limit int) []*SearchResult {
	wg := &sync.WaitGroup{}
	result := make([]*LocationWrapper, limit*s.parallelism)
	textRune := []rune(strings.ToLower(strings.TrimSpace(text)))
	for i := 0; i < s.parallelism; i++ {
		wg.Add(1)
		go s.search(i, wg, textRune, limit, result[i*limit:(i+1)*limit])
	}
	wg.Wait()
	h := ReverseHeap(result)
	heap.Init(&h)
	ret := make([]*SearchResult, 0, limit)
	for i := 0; i < limit; i++ {
		v := heap.Pop(&h).(*LocationWrapper)
		ret = append(ret, &SearchResult{
			Location: v.Location,
			Score:    v.Score,
		})
	}

	return ret
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
