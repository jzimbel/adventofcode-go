package color

import (
	"fmt"
	"github.com/mgutz/ansi"
)

// Functions that surround a value in ANSI escape codes to make it stand out when printed.
var (
	// R surrounds a value in ANSI escape codes to make it red when printed.
	R func(interface{}) string
	// G surrounds a value in ANSI escape codes to make it green when printed.
	G func(interface{}) string
	// Y surrounds a value in ANSI escape codes to make it yellow when printed.
	Y func(interface{}) string
	// B surrounds a value in ANSI escape codes to make it blue when printed.
	B func(interface{}) string
)

func morePermissiveColorFunc(f func(string) string) func(interface{}) string {
	return func(v interface{}) string {
		s := fmt.Sprintf("%v", v)
		return f(s)
	}
}

func init() {
	R = morePermissiveColorFunc(ansi.ColorFunc(ansi.Red))
	G = morePermissiveColorFunc(ansi.ColorFunc(ansi.Green))
	Y = morePermissiveColorFunc(ansi.ColorFunc(ansi.Yellow))
	B = morePermissiveColorFunc(ansi.ColorFunc(ansi.Blue))
}
