package search

import (
	"sort"
	"strings"

	"github.com/datainq/teryt"
)

type SearchOld struct {
	nodes []*LocationWrapper
}

func NewSearchOld(localities []*teryt.Location) *SearchOld {
	var nodes []*LocationWrapper
	for _, l := range localities {
		nodes = append(nodes, &LocationWrapper{
			Location: l,
			Score:    0,
		})
	}

	return &SearchOld{nodes}
}

func (s *SearchOld) Search(text string, limit int) []*SearchResult {
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
