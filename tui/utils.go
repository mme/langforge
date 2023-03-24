package tui

import (
	"fmt"

	"github.com/pterm/pterm"
)

func DisplayBanner() {
	EmptyLine()
	EmptyLine()
	lines := []string{
		" _                      _____                    ",
		"| |    __ _ _ __   __ _|  ___|__  _ __ __ _  ___ ",
		"| |   / _` | '_ \\ / _` | |_ / _ \\| '__/ _` |/ _ \\",
		"| |__| (_| | | | | (_| |  _| (_) | | | (_| |  __/",
		"|_____\\__,_|_| |_|\\__, |_|  \\___/|_|  \\__, |\\___|",
		"                  |___/               |___/      ",
	}

	for _, line := range lines {
		fmt.Println(line)
	}
	EmptyLine()
}

func EmptyLine() {
	fmt.Println()
}

func Bold(text string, args ...any) string {
	return pterm.Bold.Sprintf(text, args...)
}
