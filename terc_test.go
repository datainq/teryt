package teryt

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func mustT(t time.Time, err error) time.Time {
	if err != nil {
		panic(err)
	}
	return t
}

func TestTERCFromFile(t *testing.T) {
	f, err := os.Open("terc.csv")
	require.NoError(t, err)
	defer f.Close()

	got, err := TERCFromReader(f)
	require.NoError(t, err)

	want := []SetTERC{
		{
			Woj:    2,
			Pow:    0,
			Gmi:    0,
			Rodz:   0,
			Nazwa:  "DOLNOŚLĄSKIE",
			Nazdod: "województwo",
			StanNa: mustT(time.Parse("2006-01-02", "2019-01-01")),
		},
		{
			Woj:    2,
			Pow:    1,
			Gmi:    0,
			Rodz:   0,
			Nazwa:  "bolesławiecki",
			Nazdod: "powiat",
			StanNa: mustT(time.Parse("2006-01-02", "2019-01-01")),
		},
		{
			Woj:    2,
			Pow:    1,
			Gmi:    1,
			Rodz:   1,
			Nazwa:  "Bolesławiec",
			Nazdod: "gmina miejska",
			StanNa: mustT(time.Parse("2006-01-02", "2019-01-01")),
		},
	}
	require.Equal(t, want, got)
}
