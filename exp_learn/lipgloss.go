package explearn

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// you just have to print the ouput
// new sttyle, bold. backgound, and forrground color, paddig, width, sttyle.Render

func TestLipGloss() {

	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#b4b7d6")).
		Width(100).
		AlignHorizontal(lipgloss.Center)
	 
	fmt.Println(style.Render("hello world"))
}
