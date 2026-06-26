package styles

import "github.com/charmbracelet/lipgloss"

var (
	BG0      = lipgloss.Color("#1D2021")
	BG1      = lipgloss.Color("#282828")
	BG3      = lipgloss.Color("#3C3836")
	BG4      = lipgloss.Color("#3C3836")
	BG5      = lipgloss.Color("#504945")
	BGRed    = lipgloss.Color("#3C1F1E")
	BGGreen  = lipgloss.Color("#32361A")
	BGVisual = lipgloss.Color("#473C29")

	FG0    = lipgloss.Color("#D4BE98")
	FG1    = lipgloss.Color("#DDC7A1")
	Red    = lipgloss.Color("#EA6962")
	Green  = lipgloss.Color("#A9B665")
	Blue   = lipgloss.Color("#7DAEA3")
	Yellow = lipgloss.Color("#D8A657")
	Orange = lipgloss.Color("#E78A4E")
	Purple = lipgloss.Color("#D3869B")
	Aqua   = lipgloss.Color("#89B482")

	Grey0 = lipgloss.Color("#7C6F64")
	Grey1 = lipgloss.Color("#928374")
	Grey2 = lipgloss.Color("#A89984")
)

var (
	AppStyle = lipgloss.NewStyle()

	TitleStyle = lipgloss.NewStyle().
			Bold(true).Foreground(Green)

	HeadingStyle = lipgloss.NewStyle().
			Bold(true).Foreground(Green)

	TextStyle = lipgloss.NewStyle().
			Foreground(FG0)

	DimTextStyle = lipgloss.NewStyle().
			Foreground(Grey1)

	ErrorTextStyle = lipgloss.NewStyle().
			Foreground(Red).Bold(true)

	SuccessTextStyle = lipgloss.NewStyle().
				Foreground(Green).Bold(true)

	HighlightStyle = lipgloss.NewStyle().
			Foreground(BG0).Background(Yellow).Bold(true)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(Yellow).Bold(true)

	FooterStyle = lipgloss.NewStyle().
			Foreground(Grey1)

	HeaderStyle = lipgloss.NewStyle().
			Foreground(Green)

	CardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BG4).
			Padding(1, 2)

	ModalStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(Yellow).
			Background(BG1).
			Padding(2, 4)
)

var TierColors = map[string]lipgloss.Color{
	"Wood":       lipgloss.Color("#8B7355"),
	"Bronze":     lipgloss.Color("#CD7F32"),
	"Silver":     lipgloss.Color("#BDC3C7"),
	"Gold":       lipgloss.Color("#D8A657"),
	"Employee":   lipgloss.Color("#D3869B"),
	"Supervisor": lipgloss.Color("#E78A4E"),
	"Manager":    lipgloss.Color("#A9B665"),
	"Owner":      lipgloss.Color("#7DAEA3"),
}

func TierBadgeStyle(name string) lipgloss.Style {
	c, ok := TierColors[name]
	if !ok {
		c = Blue
	}
	return lipgloss.NewStyle().
		Foreground(c).
		Bold(true)
}

func StatusStyle(status string) lipgloss.Style {
	switch status {
	case "active":
		return lipgloss.NewStyle().Foreground(Green).Bold(true)
	case "overdue":
		return lipgloss.NewStyle().Foreground(Red).Bold(true)
	case "returned":
		return lipgloss.NewStyle().Foreground(Grey1)
	default:
		return TextStyle
	}
}

func FormatBadge(format string) string {
	switch format {
	case "VHS":
		return lipgloss.NewStyle().Foreground(Orange).Render("📼 VHS")
	case "DVD":
		return lipgloss.NewStyle().Foreground(Aqua).Render("📀 DVD")
	case "Blu-ray":
		return lipgloss.NewStyle().Foreground(Blue).Render("💿 BD")
	default:
		return format
	}
}

var (
	Cyan       = Aqua
	Magenta    = Purple
	NeonGreen  = Green
	Error      = Red
	Success    = Green
	Background = lipgloss.Color("")
	Surface    = lipgloss.Color("")
	BorderDim  = BG5
	Warning    = Yellow
	WarningAmb = Orange
	ErrorRed   = Red
	SuccessGrn = Green
	FreshGreen = Green
	SkyBlue    = Aqua
	Coral      = Orange
	Gold       = Yellow
	GlassBlue  = BG4
	BgWhite    = lipgloss.Color("")
	BgBlue     = lipgloss.Color("")
	BgLight    = lipgloss.Color("")
	TextDark   = FG0
	TextMedium = Grey1
	TextLight  = Grey0
	SoftWhite  = FG1
)

var (
	BronzeColor = TierColors["Bronze"]
	SilverColor = TierColors["Silver"]
	GoldColor   = Yellow
	EmpColor    = Purple
	SupColor    = Orange
	MgrColor    = Green
	OwnColor    = Aqua
)
