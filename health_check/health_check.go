package main

import (
	"fmt"
	"healthcheck/internal/cpu"
	"healthcheck/internal/memory"
	tui "healthcheck/internal/pkg"
)

var style = tui.DefaultStyles()

func main() {
	cpu.CPUHealth()
	memory.Memory()
	fmt.Println(style.Info.Render("(+) Everything seems normal. Thank you for using my application ദ്ദി(˵ •̀ ᴗ - ˵ ) ✧"))
}