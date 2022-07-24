package src

import (
	"fmt"
	"os"

	"github.com/gookit/color"
)

// something similar to printf but for checker logs
func Logf(format string, level int, a ...interface{}) {
	text := fmt.Sprintf(format, a...)

	switch level {
	case 0:
		fmt.Printf("%s %s\n", color.FgCyan.Render("[info]"), color.FgBlue.Render(text))
	case 1:
		fmt.Printf("%s %s\n", color.FgGreen.Render("[live]"), color.FgLightGreen.Render(text))

	case 2:
		fmt.Printf("%s %s\n", color.FgYellow.Render("[die]"), color.FgLightYellow.Render(text))
	case 3:
		fmt.Printf("%s %s\n", color.FgRed.Render("[error]"), color.FgDarkGray.Render(text))
		os.Exit(1)
	}
}
