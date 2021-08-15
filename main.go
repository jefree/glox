package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jefree/glox/errors"
)

func main() {
	if len(os.Args) > 2 {
		printHelp()
		os.Exit(errors.Usage)
	}

	if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

func runFile(path string) {
	bytes, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Println(err)
		os.Exit(errors.File)
	}

	run(string(bytes))
}

func runPrompt() {
	fmt.Println("insert an interactive shell here")
}

func run(source string) {
	fmt.Println("execute:\n\n" + source)
}

func printHelp() {
	fmt.Println("Usage: glox [script]")
}
