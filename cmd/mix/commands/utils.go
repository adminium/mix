package commands

import (
	"github.com/gozelle/color"
	"os"
)

func Warning(format string, a ...any) {
	color.Yellow(format, a...)
}

func Info(format string, a ...any) {
	color.Green(format, a...)
}

func Fatal(err error) {
	color.Red(err.Error())
	os.Exit(1)
}
