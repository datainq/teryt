package teryt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSIMCFromFile(t *testing.T) {
	got, err := SIMCFromFile("simc.csv")
	require.NoError(t, err)
	want := []SetSIMC{
		{
			Woj:     2,
			Pow:     24,
			Gmi:     4,
			RodzGmi: 2,
			Rm:      1,
			Mz:      true,
			Nazwa:   "Jemna",
			Sym:     "0855434",
			SymPod:  "0855434",
			StanNa:  mustT(time.Parse("2006-01-02", "2019-01-01"))},
		{
			Woj:     2,
			Pow:     24,
			Gmi:     4,
			RodzGmi: 2,
			Rm:      1,
			Mz:      true,
			Nazwa:   "Różana",
			Sym:     "0855486",
			SymPod:  "0855486",
			StanNa:  mustT(time.Parse("2006-01-02", "2019-01-01"))},
		{
			Woj:     2,
			Pow:     15,
			Gmi:     2,
			RodzGmi: 2,
			Rm:      1,
			Mz:      true,
			Nazwa:   "Piskorzów",
			Sym:     "0874414",
			SymPod:  "0874414",
			StanNa:  mustT(time.Parse("2006-01-02", "2019-01-01"))},
	}
	require.Equal(t, want, got)
}
