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
	voivodeshipNames := make(map[int]string)
	countyNames := make(map[int]string)
	communeByKey := make(map[keyType]*Location)
	root := &Location{}
	for _, v := range tercData {
		key := keyType{
			voivodeship: v.Woj,
			county:      v.Pow,
			commune:     v.Gmi,
		}
		if v.Pow == 0 {
			voivodeshipNames[v.Woj] = v.Nazwa
			loc := &Location{
				ID:       "",
				Type:     v.Nazdod,
				Name:     v.Nazwa,
				Parts:    []string{v.Nazwa},
				Parent:   root,
				Children: nil,
			}
			root.Children = append(root.Children, loc)
			communeByKey[key] = loc
		} else if v.Gmi == 0 {
			countyNames[v.Pow] = v.Nazwa

			parent, ok := communeByKey[keyType{voivodeship: v.Woj}]
			if !ok {
				panic("not found")
			}
			loc := &Location{
				ID:       "",
				Type:     v.Nazdod,
				Name:     v.Nazwa,
				Parts:    append(append([]string{}, parent.Parts...), v.Nazwa),
				Parent:   parent,
				Children: nil,
			}
			parent.Children = append(parent.Children, loc)
			communeByKey[key] = loc
		} else {
			parent, ok := communeByKey[keyType{
				voivodeship: v.Woj,
				county:      v.Pow,
			}]
			if !ok {
				panic("not found")
			}
			loc := &Location{
				ID:       "",
				Type:     v.Nazdod,
				Name:     v.Nazwa,
				Parts:    append(append([]string{}, parent.Parts...), v.Nazwa),
				Parent:   parent,
				Children: nil,
			}
			parent.Children = append(parent.Children, loc)
			communeByKey[key] = loc
		}
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
