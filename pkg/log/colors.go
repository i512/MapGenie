package log

import "fmt"

const (
	Red      = "\033[31;1m"
	LightRed = "\033[31m"
	Yellow   = "\033[33m"
	Cyan     = "\033[36m"
	NoColor  = "\033[0m"

	FColor = LightRed
	EColor = Red
	WColor = Yellow
	IColor = ""
	DColor = Cyan
)

func Color(color string, s string) string {
	if color == "" {
		return s
	}

	return fmt.Sprintf("%s%s%s", color, s, NoColor)
}
