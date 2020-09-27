package teryt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestULICFromFile(t *testing.T) {
	got, err := ULICFromFile("ulic.csv")
	require.NoError(t, err)
	want := []SetULIC{
		{
			Woj:     4,
			Pow:     10,
			Gmi:     4,
			RodzGmi: 2,
			Sym:     "0094774",
			SymUl:   "19326",
			Cecha:   "ul.",
			Nazwa1:  "Rzeczna",
			Nazwa2:  "",
			StanNa:  mustT(time.Parse("2006-01-02", "2019-08-30")),
		},
		{
			Woj:     4,
			Pow:     10,
			Gmi:     4,
			RodzGmi: 2,
			Sym:     "0094870",
			SymUl:   "08646",
			Cecha:   "ul.",
			Nazwa1:  "Kmieciaka",
			Nazwa2:  "",
			StanNa:  mustT(time.Parse("2006-01-02", "2019-08-30")),
		},
		{
			Woj:     4,
			Pow:     10,
			Gmi:     5,
			RodzGmi: 4,
			Sym:     "0929612",
			SymUl:   "08153",
			Cecha:   "ul.",
			Nazwa1:  "Kasprowicza",
			Nazwa2:  "Jana",
			StanNa:  mustT(time.Parse("2006-01-02", "2019-08-30")),
		},
	}
	require.Equal(t, want, got)
}
