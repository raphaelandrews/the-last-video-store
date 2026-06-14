package styles

import "github.com/charmbracelet/lipgloss"

var (
	Cyan      = lipgloss.Color("#00FFFF")
	Magenta   = lipgloss.Color("#FF00FF")
	Yellow    = lipgloss.Color("#FFFF00")
	NeonGreen = lipgloss.Color("#39FF14")
	NeonPink  = lipgloss.Color("#FF6EC7")

	Background = lipgloss.Color("#0A0A2E")
	Surface    = lipgloss.Color("#121240")
	BorderDim  = lipgloss.Color("#333366")

	Error   = lipgloss.Color("#FF4444")
	Success = lipgloss.Color("#44FF44")
	Warning = lipgloss.Color("#FFAA00")

	BronzeColor = lipgloss.Color("#CD7F32")
	SilverColor = lipgloss.Color("#C0C0C0")
	GoldColor   = lipgloss.Color("#FFD700")
	EmpColor    = lipgloss.Color("#FF00FF")
	SupColor    = lipgloss.Color("#FF8C00")
	MgrColor    = lipgloss.Color("#FFFF00")
	OwnColor    = lipgloss.Color("#00FFFF")
)

var (
	AppStyle = lipgloss.NewStyle().
			Background(Background).
			Foreground(lipgloss.Color("#CCCCCC"))

	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Cyan).
			BorderBackground(Background)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Cyan).
			Background(Background)

	HeadingStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Yellow).
			Background(Background)

	TextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CCCCCC")).
			Background(Background)

	DimTextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666688")).
			Background(Background)

	ErrorTextStyle = lipgloss.NewStyle().
			Foreground(Error).
			Background(Background).
			Bold(true)

	SuccessTextStyle = lipgloss.NewStyle().
				Foreground(Success).
				Background(Background).
				Bold(true)

	HighlightStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#1A1A4E")).
			Foreground(Cyan).
			Bold(true)

	SelectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#222266")).
			Foreground(NeonGreen)

	FooterStyle = lipgloss.NewStyle().
			Background(Surface).
			Foreground(lipgloss.Color("#8888AA")).
			Padding(0, 1)

	HeaderStyle = lipgloss.NewStyle().
			Background(Background).
			Foreground(Cyan)

	ModalOverlayStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#000000")).
				Foreground(lipgloss.Color("#CCCCCC"))

	ModalBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(Cyan).
			Padding(2, 4).
			Background(Surface)
)

var TierColors = map[string]lipgloss.Color{
	"Bronze":     BronzeColor,
	"Silver":     SilverColor,
	"Gold":       GoldColor,
	"Employee":   EmpColor,
	"Supervisor": SupColor,
	"Manager":    MgrColor,
	"Owner":      OwnColor,
}

func TierBadgeStyle(tierName string) lipgloss.Style {
	color, ok := TierColors[tierName]
	if !ok {
		color = Cyan
	}
	return lipgloss.NewStyle().
		Background(color).
		Foreground(lipgloss.Color("#000000")).
		Bold(true).
		Padding(0, 1)
}

func StatusStyle(status string) lipgloss.Style {
	switch status {
	case "active":
		return lipgloss.NewStyle().Foreground(NeonGreen).Bold(true)
	case "overdue":
		return lipgloss.NewStyle().Foreground(Error).Bold(true)
	case "returned":
		return lipgloss.NewStyle().Foreground(DimTextStyle.GetForeground())
	default:
		return TextStyle
	}
}
