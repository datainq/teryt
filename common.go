package teryt

type RodzGmi int

var (
	RodzGmi_miejska                = RodzGmi(1)
	RodzGmi_wiejska                = RodzGmi(2)
	RodzGmi_miejskoWiejska         = RodzGmi(3)
	RodzGmi_miastoWGminieMW        = RodzGmi(4)
	RodzGmi_obszarWiejskiWGminieMW = RodzGmi(5)
	RodzGmi_dzielnicaWarszawy      = RodzGmi(8)
	RodzGmi_delegatura             = RodzGmi(9)
)

func (r RodzGmi) Name() string {
	switch r {
	case 1:
		return "miejska"
	case 2:
		return "wiejska"
	case 3:
		return "miejsko-wiejska"
	case 4:
		return "miasto w gminie miejsko-wiejskiej"
	case 5:
		return "obszar wiejski w gminie miejsko-wiejskiej"
	case 8:
		return "dzielnica w m.st. Warszawa"
	case 9:
		return "delegatury miast: Kraków, Łódź, Poznań i Wrocław"
	}
	return ""
}

type RodzMiejscowosci int

func GetRodzMiejscowosci(v string) int {
	switch v {
	case "00":
		return 0
	case "01":
		return 1
	case "02":
		return 2
	case "03":
		return 3
	case "04":
		return 4
	case "05":
		return 5
	case "06":
		return 6
	case "07":
		return 7
	case "95":
		return 95
	case "96":
		return 96
	case "98":
		return 98
	case "99":
		return 99
	}
	return -1
}

func (r RodzMiejscowosci) Name() string {
	switch r {
	case 0:
		return "część miejscowości"
	case 1:
		return "wieś"
	case 2:
		return "kolonia"
	case 3:
		return "przysiółek"
	case 4:
		return "osada"
	case 5:
		return "osada leśna"
	case 6:
		return "osiedle"
	case 7:
		return "schronisko turystyczne"
	case 95:
		return "dzielnica m. st. Warszawy"
	case 96:
		return "miasto"
	case 98:
		return "delegatura"
	case 99:
		return "część miasta"
	}
	return ""
}
