package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

var hadError = false

func main() {
	if len(os.Args) > 2 {
		printHelp()
		os.Exit(UsageExitCode)
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
		os.Exit(FileExitCode)
	}

	run(string(bytes))
}

func runPrompt() {
	fmt.Println("insert an interactive shell here")
}

func run(source string) {
	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()

	parser := NewParser(tokens)
	expr := parser.Parse()

	fmt.Println(expr)

	if hadError {
		return
	}

	fmt.Println(AstPrinter{}.Print(expr))
}

func report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error %s: %s\n", line, where, message)
	hadError = true
}

func fail(token Token, message string) {
	if token.Kind == EOF {
		report(token.Line, "at end", message)
	} else {
		report(token.Line, "at '"+token.Lexeme+"'", message)
	}
}

func printHelp() {
	fmt.Println("Usage: glox [script]")
}
