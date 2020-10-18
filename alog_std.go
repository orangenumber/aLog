// (c) 2020 Gon Y Yi. <https://gonyyi.com/copyright.txt>

package alog

import (
	"io"
	"os"
)

// =====================================================================================================================
// ALOG STD
// =====================================================================================================================
var std = New(os.Stdout, "", F_STD)

func Printf(format string, s ...interface{}) {
	std.Printf(format, s...)
}
func Printj(optionalPrefix string, a interface{}) {
	std.Printj(optionalPrefix, a)
}
func Print(s ...interface{}) {
	std.Print(s...)
}
func SetOutput(output io.Writer) {
	std.SetOutput(output)
}
func SetPrefix(prefix string) {
	std.SetPrefix(prefix)
}
func SetFlag(flag uint16) {
	std.SetFlag(flag)
}