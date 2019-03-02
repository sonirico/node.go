package repl

import (
	"bufio"
	"io"
	"node.go/evaluator"
	"node.go/lexer"
	"node.go/object"
	"node.go/parser"
)

const (
	PROMPT      = "(o_o) > "
	PROMPT_OOPS = "(T_T) >"
)

func printPrompt(out io.Writer, lastStatusCode byte) {
	if lastStatusCode < 1 {
		io.WriteString(out, PROMPT)
	} else {
		io.WriteString(out, PROMPT_OOPS)
	}
}

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	var lastStatus byte = 0
	var environment = object.NewEnvironment()

	for {
		printPrompt(out, lastStatus)
		_ = scanner.Scan()

		line := scanner.Text()
		lex := lexer.New(line)
		par := parser.New(lex)
		program := par.ParseProgram()

		if len(par.Errors()) > 0 {
			lastStatus = 1
			for _, errorMessage := range par.Errors() {
				io.WriteString(out, errorMessage)
				io.WriteString(out, "\n")
			}
			io.WriteString(out, "\n")
		} else {
			evaluatedObject := evaluator.Eval(program, environment)
			if evaluatedObject.Type() == object.ERROR {
				lastStatus = 1
			} else {
				lastStatus = 0
			}
			io.WriteString(out, "\n")
			io.WriteString(out, evaluatedObject.Inspect())
			io.WriteString(out, "\n")
		}
	}
}
