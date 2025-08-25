package explearn

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// you just have to print the ouput
// new sttyle, bold. backgound, and forrground color, paddig, width, sttyle.Render

func TestLipGloss() {

	style := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	fmt.Println(style.Render("hello world"))
}
