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

func (s *SearchV3) search(idx int, wg *sync.WaitGroup, text string, limit int, result chan<- *LocationWrapper) {
	defer wg.Done()
	textRune := []rune(strings.ToLower(strings.TrimSpace(text)))
	maxScore := 100000
	h := &Heap{}
	for i, v := range s.chunks[idx] {
		v.Score = LevenshteinRune(v.SearchText, textRune)
		if i < limit {
			heap.Push(h, v)
			maxScore = v.Score
		} else if v.Score < maxScore {
			heap.Push(h, v)
			maxScore = heap.Pop(h).(*LocationWrapper).Score
		}
	}
	for _, v := range *h {
		result <- v
	}
}

func (s *SearchV3) Search(text string, limit int) []*SearchResult {
	wg := &sync.WaitGroup{}
	result := make(chan *LocationWrapper, 10)
	for i := 0; i < s.parallelism; i++ {
		wg.Add(1)
		go s.search(i, wg, text, limit, result)
	}
	go func() {
		wg.Wait()
		close(result)
	}()

	h := &ReverseHeap{}
	for v := range result {
		heap.Push(h, v)
	}
	ret := make([]*SearchResult, 0, limit)
	for i := 0; i < limit; i++ {
		v := heap.Pop(h).(*LocationWrapper)
		ret = append(ret, &SearchResult{
			Location: v.Location,
			Score:    v.Score,
		})
	}

	return ret
}
