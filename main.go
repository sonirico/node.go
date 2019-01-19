package main

import (
	"node.go/repl"
	"os"
)

func main () {
	repl.Start(os.Stdin, os.Stdout)
}