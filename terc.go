package teryt

import (
	"encoding/csv"
	"io"
	"os"
	"time"

	"github.com/datainq/csvlib"
)

type SetTERC struct {
	// WOJ - symbol województwa - 2 zn. C
	Woj int

	// POW - symbol powiatu - 2 zn. C
	Pow int

	// GMI - symbol gminy - 2 zn. C
	Gmi int

	// RODZ - symbol rodzaju jednostki
	// 1 - gmina miejska,
	// 2 - gmina wiejska,
	// 3 - gmina miejsko-wiejska,
	// 4 - miasto w gminie miejsko-wiejskiej,
	// 5 - obszar wiejski w gminie miejsko-wiejskiej,
	// 8 - dzielnica w m.st. Warszawa,
	// 9 - delegatury miast: Kraków, Łódź, Poznań i Wrocław
	// 1 zn. C
	Rodz int

	// NAZWA - nazwa województwa/ powiatu/ gminy - 100 zn. C
	Nazwa string

	// NAZDOD - określenie jednostki - (województwo; powiat; miasto na prawach powiatu; miasto stołeczne, na prawach powiatu; gmina miejska, miasto stołeczne; gmina miejska; gmina wiejska; gmina miejsko-wiejska; miasto; obszar wiejski; dzielnica; delegatura) -50 zn. C
	Nazdod string

	// STAN_NA - data aktualizacji danych w systemie TERC w formacie RRRR-MM-DD - 10 zn. C
	StanNa time.Time
}

func TERCFromReader(reader io.Reader) ([]SetTERC, error) {
	var ret []SetTERC
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
		p, err := TERCParser.Parse(row)
		if err != nil {
			return nil, err
		}
		ret = append(ret, SetTERC{
			Woj:    int(p[0].I32()),
			Pow:    int(p[1].I32()),
			Gmi:    int(p[2].I32()),
			Rodz:   int(p[3].I32()),
			Nazwa:  p[4].String(),
			Nazdod: p[5].String(),
			StanNa: p[6].Time(),
		})
	}
	return ret, nil
}

func TERCFromFile(file string) ([]SetTERC, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return TERCFromReader(f)
}

var TERCParser = &csvlib.RowParser{
	P: []csvlib.Parser{
		csvlib.Int32Parser{Name: "WOJ", Optional: true, Default: 0},
		csvlib.Int32Parser{Name: "POW", Optional: true, Default: 0},
		csvlib.Int32Parser{Name: "GMI", Optional: true, Default: 0},
		csvlib.Int32Parser{Name: "RODZ", Optional: true, Default: 0},
		csvlib.StringParser{Name: "NAZWA"},
		csvlib.StringParser{Name: "NAZDOD"},
		csvlib.TimeParser{Name: "STAN_NA", Layout: "2006-02-02"},
	},
}
