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
	AppStyle = lipgloss.NewStyle().
			Foreground(FG0).Background(BG0)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).Foreground(Green).Background(BG0)

	HeadingStyle = lipgloss.NewStyle().
			Bold(true).Foreground(Green).Background(BG0)

	TextStyle = lipgloss.NewStyle().
			Foreground(FG0).Background(BG0)

	DimTextStyle = lipgloss.NewStyle().
			Foreground(Grey1).Background(BG0)

	ErrorTextStyle = lipgloss.NewStyle().
			Foreground(Red).Bold(true).Background(BG0)

	SuccessTextStyle = lipgloss.NewStyle().
				Foreground(Green).Bold(true).Background(BG0)

	HighlightStyle = lipgloss.NewStyle().
			Foreground(BG0).Background(Green).Bold(true)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(FG0).Background(BGVisual)

	FooterStyle = lipgloss.NewStyle().
			Background(BG1).Foreground(Grey1).Padding(0, 1)

	HeaderStyle = lipgloss.NewStyle().
			Foreground(Green).Background(BG0)

	CardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BG4).
			Background(BG1).
			Padding(1, 2)

	ModalStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(Yellow).
			Background(BG1).
			Padding(2, 4)
)

var TierColors = map[string]lipgloss.Color{
	"Bronze":     lipgloss.Color("#CD7F32"),
	"Silver":     lipgloss.Color("#BDC3C7"),
	"Gold":       Yellow,
	"Employee":   Purple,
	"Supervisor": Orange,
	"Manager":    Green,
	"Owner":      Blue,
}

func TierBadgeStyle(name string) lipgloss.Style {
	c, ok := TierColors[name]
	if !ok {
		c = Blue
	}
	return lipgloss.NewStyle().
		Background(c).
		Foreground(BG0).
		Bold(true).
		Padding(0, 1)
}

func StatusStyle(status string) lipgloss.Style {
	switch status {
	case "active":
		return lipgloss.NewStyle().Foreground(Green).Bold(true).Background(BG0)
	case "overdue":
		return lipgloss.NewStyle().Foreground(Red).Bold(true).Background(BG0)
	case "returned":
		return lipgloss.NewStyle().Foreground(Grey1).Background(BG0)
	default:
		return TextStyle
	}
}

func FormatBadge(format string) string {
	switch format {
	case "VHS":
		return lipgloss.NewStyle().Foreground(Orange).Background(BG0).Render("📼 VHS")
	case "DVD":
		return lipgloss.NewStyle().Foreground(Blue).Background(BG0).Render("📀 DVD")
	case "Blu-ray":
		return lipgloss.NewStyle().Foreground(Aqua).Background(BG0).Render("💿 BD")
	default:
		return format
	}
}

var (
	Cyan       = Blue
	Magenta    = Orange
	NeonGreen  = Green
	Error      = Red
	Success    = Green
	Background = BG0
	Surface    = BG1
	BorderDim  = BG4
	Warning    = Yellow
	WarningAmb = Yellow
	ErrorRed   = Red
	SuccessGrn = Green
	FreshGreen = Green
	SkyBlue    = Green
	Coral      = Orange
	Gold       = Yellow
	GlassBlue  = BG4
	BgWhite    = BG1
	BgBlue     = BG0
	BgLight    = BG0
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
	OwnColor    = Blue
)
