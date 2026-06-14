package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/thelastvideostore/tui"
)

func main() {
	apiURL := flag.String("api-url", "http://localhost:8080", "API server base URL")
	debug := flag.Bool("debug", false, "enable debug mode")
	flag.Parse()

	if *debug {
		f, err := tea.LogToFile("tui-debug.log", "")
		if err != nil {
			fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	model := tui.NewModel(*apiURL)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
