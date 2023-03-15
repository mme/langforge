package tui

import (
	"fmt"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

func DisplayBanner() {
	EmptyLine()
	EmptyLine()
	pterm.DefaultBigText.WithLetters(putils.LettersFromString("LangForge")).Render()
	EmptyLine()
}

func EmptyLine() {
	fmt.Println()
}

func Bold(text string) string {
	return pterm.Bold.Sprintf(text)
}
