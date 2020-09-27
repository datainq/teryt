package teryt

import (
	"encoding/csv"
	"io"
	"os"
	"time"

	"github.com/datainq/csvlib"
)

type SetULIC struct {
	// WOJ - symbol województwa - 2 zn. C
	Woj int

	// POW - symbol powiatu - 2 zn. C
	Pow int

	// GMI - symbol gminy - 2 zn. C
	Gmi int

	// RODZ_GMI - symbol rodzaju jednostki
	// 1 - gmina miejska,
	// 2 - gmina wiejska,
	// 3 - gmina miejsko-wiejska,
	// 4 - miasto w gminie miejsko-wiejskiej,
	// 5 - obszar wiejski w gminie miejsko-wiejskiej,
	// 8 - dzielnica w m.st. Warszawa,
	// 9 - delegatury miast: Kraków, Łódź, Poznań i Wrocław
	// 1 zn. C
	RodzGmi RodzGmi

	// SYM - identyfikator miejscowości - 7 zn. C
	Sym string

	// SYM_UL - identyfikator nazwy ulicy - 5 zn. C
	SymUl string

	// CECHA - określenie rodzaju ulicy - 5 zn. C
	// (ul., al., pl., skwer, bulw., rondo, park, rynek, szosa, droga, os.,
	// ogród, wyspa, wyb., inne)
	Cecha string

	// NAZWA_1 - część nazwy począwszy od słowa, które decyduje o pozycji
	// ulicy w układzie alfabetycznym, aż do końca nazwy - 100 zn. C
	Nazwa1 string

	// NAZWA_2 - pozostała część nazwy lub pole puste - 100 zn. C
	//
	// W przypadku, gdy pole Nazwa_2 nie jest puste, aby otrzymać nazwę ulicy
	// w pełnym brzmieniu, człony nazwy należy ułożyć w kolejności:
	// Nazwa_2, Nazwa_1.
	Nazwa2 string

	// STAN_NA - data aktualizacji danych w podsystemie ULIC w formacie RRRR-MM-DD. - 10 zn. C
	StanNa time.Time
}

type SetWMRODZ struct {
	// RM - symbol rodzaju miejscowości - 2 zn. C
	Rm int

	// NAZWA_RM - nazwa rodzaju miejscowości - 100 zn. C
	NazwaRm string

	// STAN_NA - data aktualizacji danych w podsystemie WMRODZ w formacie RRRR-MM-DD. - 10 zn. C
	StanNa time.Time
}

func ULICFromReader(reader io.Reader) ([]SetULIC, error) {
	var ret []SetULIC
	r := csv.NewReader(reader)
	r.Comma = ';'
	for idx := 0; ; idx++ {
		row, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		if idx == 0 {
			continue
		}
		p, err := ULICParser.Parse(row)
		if err != nil {
			return nil, err
		}
		ret = append(ret, SetULIC{
			Woj:     int(p[0].I32()),
			Pow:     int(p[1].I32()),
			Gmi:     int(p[2].I32()),
			RodzGmi: RodzGmi(p[3].I32()),
			Sym:     p[4].String(),
			SymUl:   p[5].String(),
			Cecha:   p[6].String(),
			Nazwa1:  p[7].String(),
			Nazwa2:  p[8].String(),
			StanNa:  p[9].Time(),
		})
	}
	return ret, nil
}

func ULICFromFile(file string) ([]SetULIC, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ULICFromReader(f)
}

var ULICParser = &csvlib.RowParser{
	P: []csvlib.Parser{
		csvlib.Int32Parser{Name: "WOJ"},
		csvlib.Int32Parser{Name: "POW"},
		csvlib.Int32Parser{Name: "GMI"},
		csvlib.Int32Parser{Name: "RODZ_GMI"},
		csvlib.StringParser{Name: "SYM"},
		csvlib.StringParser{Name: "SYM_UL"},
		csvlib.StringParser{Name: "CECHA"},
		csvlib.StringParser{Name: "NAZWA_1"},
		csvlib.StringParser{Name: "NAZWA_2"},
		csvlib.TimeParser{Name: "STAN_NA", Layout: "2006-01-02"},
	},
}
