package teryt

import (
	"encoding/csv"
	"io"
	"os"
	"time"

	"github.com/datainq/csvlib"
)

type SetSIMC struct {
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

	// RM - rodzaj miejscowości
	// 00 - część miejscowości
	// 01 - wieś
	// 02 - kolonia
	// 03 - przysiółek
	// 04 - osada
	// 05 - osada leśna
	// 06 - osiedle
	// 07 - schronisko turystyczne
	// 95 - dzielnica m. st. Warszawy
	// 96 - miasto
	// 98 - delegatura
	// 99 - część miasta
	// 2 zn. C
	Rm RodzMiejscowosci

	// MZ - występowanie nazwy zwyczajowej (0-tak,1-nie) - 1 zn. C
	Mz bool

	// NAZWA - nazwa miejscowości - 100 zn. C
	Nazwa string

	// SYM - identyfikator miejscowości - 7 zn. C
	Sym string

	// SYMPOD - identyfikator miejscowości podstawowej
	// - dla części miejscowości wiejskich - identyfikator miejscowości, do której dana część należy,
	// - dla części miast - identyfikator danego miasta (w miastach posiadających dzielnice/delegatury - identyfikator tej jednostki).
	// 7 zn. C
	SymPod string

	// STAN_NA - data aktualizacji danych w podsystemie SIMC w formacie RRRR-MM-DD. - 10 zn. C
	StanNa time.Time
}

func SIMCFromReader(reader io.Reader) ([]SetSIMC, error) {
	var ret []SetSIMC
	r := csv.NewReader(reader)
	r.Comma = ';'
	r.LazyQuotes = true
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
		p, err := SIMCParser.Parse(row)
		if err != nil {
			return nil, err
		}
		ret = append(ret, SetSIMC{
			Woj:     int(p[0].I32()),
			Pow:     int(p[1].I32()),
			Gmi:     int(p[2].I32()),
			RodzGmi: RodzGmi(p[3].I32()),
			Rm:      RodzMiejscowosci(p[4].I32()),
			Mz:      p[5].I32() == 1,
			Nazwa:   p[6].String(),
			Sym:     p[7].String(),
			SymPod:  p[8].String(),
			StanNa:  p[9].Time(),
		})
	}
	return ret, nil
}

func SIMCFromFile(file string) ([]SetSIMC, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return SIMCFromReader(f)
}

var SIMCParser = &csvlib.RowParser{
	P: []csvlib.Parser{
		csvlib.Int32Parser{Name: "WOJ", Optional: true, Default: 0},
		csvlib.Int32Parser{Name: "POW", Optional: true, Default: 0},
		csvlib.Int32Parser{Name: "GMI", Optional: true, Default: 0},
		csvlib.Int32Parser{Name: "RODZ_GMI", Optional: true, Default: 0},
		csvlib.Int32Parser{Name: "RM", Optional: true, Default: 0},
		csvlib.Int32Parser{Name: "MZ", Optional: true, Default: 0},
		csvlib.StringParser{Name: "NAZWA"},
		csvlib.StringParser{Name: "SYM"},
		csvlib.StringParser{Name: "SYMPOD"},
		csvlib.TimeParser{Name: "STAN_NA", Layout: "2006-01-02"},
	},
}
