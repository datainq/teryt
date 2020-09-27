package search

import (
	"sort"
	"strings"

	"github.com/datainq/teryt"
)

type Interface interface {
	Search(text string, limit int) []*SearchResult
}

type SearchV1 struct {
	nodes []*LocationWrapper
}

func NewSearchV1(localities []*teryt.Location) *SearchV1 {
	var nodes []*LocationWrapper
	for _, l := range localities {
		nodes = append(nodes, &LocationWrapper{
			Location: l,
			Score:    0,
		})
	}

	return &SearchV1{nodes}
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

func (s *SearchV1) Search(text string, limit int) []*SearchResult {
	text = strings.ToLower(strings.TrimSpace(text))
	for _, v := range s.nodes {
		v.Score = levenshtein([]rune(strings.ToLower(v.Parts[len(v.Parts)-1])), []rune(text))
	}
	sort.Slice(s.nodes, func(i, j int) bool {
		a, b := s.nodes[i], s.nodes[j]
		return a.Score < b.Score || a.Score == b.Score && (len(a.Parts) < len(b.Parts))
	})
	locs := s.nodes
	if len(locs) > 10 {
		locs = locs[:10]
	}

	ret := make([]*SearchResult, 0, len(locs))
	for _, v := range locs {
		ret = append(ret, &SearchResult{
			Location: v.Location,
			Score:    v.Score,
		})
	}
	return ret
}
