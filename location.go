package teryt

import (
	"fmt"
	"strings"
)

type Location struct {
	ID       string
	Type     string
	Name     string
	Parts    []string
	FullName string

	Parent   *Location
	Children []*Location
}

type Config struct {
	Separator string
}

func (l *Location) Build(cfg Config) {
	l.FullName = strings.Join(l.Parts, cfg.Separator)
}

type keyType struct {
	voivodeship, county, commune int
}

func BuildLocations(tercData []SetTERC, simcData []SetSIMC) (*Location, error) {
	root := &Location{}
	communeByKey := map[keyType]*Location{
		keyType{}: root,
	}
	var (
		id             string
		parentKey, key keyType
	)
	for _, v := range tercData {
		key = keyType{
			voivodeship: v.Woj,
			county:      v.Pow,
			commune:     v.Gmi,
		}

		if v.Pow == 0 { // województwo
			parentKey = keyType{}
			id = fmt.Sprintf("%02d", v.Woj)
		} else if v.Gmi == 0 { // powiat
			parentKey = keyType{voivodeship: v.Woj}
			id = fmt.Sprintf("%02d%02d", v.Woj, v.Pow)
		} else { // gmina
			parentKey = keyType{voivodeship: v.Woj, county: v.Pow}
			id = fmt.Sprintf("%02d%02d%02d", v.Woj, v.Pow, v.Gmi)
		}
		parent, ok := communeByKey[parentKey]
		id = fmt.Sprintf("%02d%02d", v.Woj, v.Pow)
		if !ok {
			return nil, fmt.Errorf("cannot find parent for: %s", id)
		}
		loc := &Location{
			ID:       id,
			Type:     v.Nazdod,
			Name:     v.Nazwa,
			Parts:    append(append([]string{}, parent.Parts...), v.Nazwa),
			Parent:   parent,
			Children: nil,
		}
		parent.Children = append(parent.Children, loc)
		communeByKey[key] = loc
	}

	simBySym := make(map[string]*Location)
	for i := 0; i < 2; i++ {
		for _, v := range simcData {
			if i == 0 && v.Sym != v.SymPod {
				continue
			} else if i == 1 && v.Sym == v.SymPod {
				continue
			}
			key := keyType{
				voivodeship: v.Woj,
				county:      v.Pow,
				commune:     v.Gmi,
			}
			var parent *Location
			if v.Sym != v.SymPod {
				continue
			}
			parent = communeByKey[key]
			if parent == nil {
				return nil, fmt.Errorf("cannot find parent! %s != %s", v.Sym, v.SymPod)
			}
			loc := &Location{
				ID:       v.Sym,
				Type:     v.Rm.Name(),
				Name:     v.Nazwa,
				Parts:    append(append([]string{}, parent.Parts...), v.Nazwa),
				Parent:   parent,
				Children: nil,
			}
			parent.Children = append(parent.Children, loc)
			simBySym[v.Sym] = loc
		}
	}
	return root, nil
}
