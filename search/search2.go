package search

import (
	"container/heap"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/datainq/teryt"
)

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
		v.Score = LevenshteinRune(v.SearchText, textRune)
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

func LevenshteinStr(str1, str2 string) int {
	s1Len := utf8.RuneCountInString(str1)
	s2Len := utf8.RuneCountInString(str2)
	if s1Len > s2Len {
		str1, str2 = str2, str1
		s1Len = s2Len
	}
	column := make([]int, s1Len+2)

	for y := 1; y <= s1Len; y++ {
		column[y] = y
	}
	x := 1
	for _, str2R := range str2 {
		column[0] = x
		lastkey := x - 1
		y := 1
		for _, str1R := range str1 {
			oldkey := column[y]
			var incr int
			if str1R != str2R {
				incr = 1
			}

			column[y] = minimum3Int(column[y]+1, column[y-1]+1, lastkey+incr)
			lastkey = oldkey
			y++
		}
		x++
	}
	return column[s1Len]
}

// LevenshteinRune impelments levenstein distance.
// Optimizations taken from:
// https://github.com/agnivade/levenshtein/blob/master/levenshtein.go
func LevenshteinRune(str1, str2 []rune) int {
	if len(str1) > len(str2) {
		str1, str2 = str2, str1
	}
	s1len := len(str1)
	s2len := len(str2)
	column := make([]uint16, s1len+1)

	for i := 1; i <= s1len; i++ {
		column[i] = uint16(i)
	}
	var lastkey, oldkey uint16
	_ = column[s1len]
	var y int
	for x := 1; x <= s2len; x++ {
		lastkey = uint16(x)
		for y = 1; y <= s1len; y++ {
			oldkey = column[y-1]
			if str1[y-1] != str2[x-1] {
				oldkey = min(min(column[y-1]+1, column[y]+1), lastkey+1)
			}

			column[y-1] = lastkey
			lastkey = oldkey
		}
		column[s1len] = lastkey
	}
	return int(column[s1len])
}

func min(a, b uint16) uint16 {
	if a < b {
		return a
	}
	return b
}

func minimum3Int(a, b, c int) int {
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
