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

	ParentID string
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
	symbol                       string
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
	communeBySymbol := make(map[string]*Location)
	for _, v := range tercData {
		key = keyType{
			voivodeship: v.Woj,
			county:      v.Pow,
			commune:     v.Gmi,
		}

		if v.Pow == 0 { // województwo
			parentKey = keyType{}
			id = fmt.Sprintf("terc:%02d", v.Woj)
		} else if v.Gmi == 0 { // powiat
			parentKey = keyType{voivodeship: v.Woj}
			id = fmt.Sprintf("terc:%02d%02d", v.Woj, v.Pow)
		} else { // gmina
			parentKey = keyType{voivodeship: v.Woj, county: v.Pow}
			id = fmt.Sprintf("terc:%02d%02d%02d", v.Woj, v.Pow, v.Gmi)
		}
		parent, ok := communeByKey[parentKey]
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
		communeBySymbol[loc.ID] = loc
	}

	for _, v := range simcData {
		key := keyType{
			voivodeship: v.Woj,
			county:      v.Pow,
			commune:     v.Gmi,
		}
		if v.Sym != v.SymPod {
			key.symbol = "simc:" + v.SymPod
		}
		parent, ok := communeByKey[key]
		loc := &Location{
			ID:       "simc:" + v.Sym,
			Type:     v.Rm.Name(),
			Name:     v.Nazwa,
			Children: nil,
		}
		if ok {
			loc.Parts = append(append([]string{}, parent.Parts...), v.Nazwa)
			loc.Parent = parent
			loc.ParentID = parent.ID
			parent.Children = append(parent.Children, loc)
		} else if v.Sym != v.SymPod {
			loc.ParentID = key.symbol
		}
		key.symbol = loc.ID
		communeBySymbol[loc.ID] = loc
	}
	for _, v := range communeByKey {
		if v.ParentID == "" || v.Parent != nil {
			continue
		}
		parent, ok := communeBySymbol[v.ParentID]
		if !ok {
			return nil, fmt.Errorf("cannot find parent: %s", v.ParentID)
		}
		v.Parts = append(append([]string{}, parent.Parts...), v.Name)
		v.Parent = parent
		parent.Children = append(parent.Children, v)
	}
	return root, nil
}
