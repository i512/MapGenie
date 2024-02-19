package log

import "fmt"

const (
	Red       = "\033[0;31m"
	LightRed  = "\033[1;31m"
	Yellow    = "\033]1;33m"
	LightGray = "\033]0;38m"
	NoColor   = "\033[0m"

	FColor = LightRed
	EColor = Red
	WColor = Yellow
	IColor = ""
	DColor = LightGray
)

func Color(color string, s string) string {
	if color == "" {
		return s
	}

	return fmt.Sprintf("%s%s%s", color, s, NoColor)
}
