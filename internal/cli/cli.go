package cli

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
)

func Run() {
	if len(os.Args) < 2 {
		if err := NewWatcher().Run(); err != nil {
			log.Fatal(err)
		}
		os.Exit(1)
	}
	switch os.Args[1] {
	case "generate":
		p := NewCompilerPool(nil)
		p.GenerateAssetsAndTempl()
	case "help", "--help":
		fmt.Print("boring is a tool that bundles sass, esbuild and templ\n\n")
		fmt.Print("Usage:\n")
		fmt.Print("\tboring\t\t\tHot code reloading of .scss, .css, .ts, .js, .templ and .go files\n")
		fmt.Print("\tboring generate\t\tGenerates, bundles, minifies .scss, .css, .ts, .js and .templ files\n")
		fmt.Print("\tboring help\t\tPrints the help screen\n\n")
		fmt.Print("For more informations visit https://github.com/Instantan/boring\n\n")
	default:
		fmt.Print("Command not supported.\nUse 'boring help' for a list of supported commands\n")
	}
}

func printMeasuredAction(start string, action func() error, end string) {
	printInternal("%v\n", start)
	t := time.Now()
	err := action()
	if err != nil {
		printError(err)
		return
	}
	printInternal("%v. Took %v\n", end, time.Since(t).String())
}

func printInternal(format string, a ...interface{}) {
	color.Cyan(format, a...)
}

func printError(err error) {
	color.Red(err.Error())
}
