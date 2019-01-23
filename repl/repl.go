package repl

import (
	"bufio"
	"fmt"
	"io"
	"node.go/lexer"
	"node.go/parser"
)

const PROMPT = "\n/> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		io.WriteString(out, PROMPT)
		_ = scanner.Scan()

		line := scanner.Text()
		lex := lexer.New(line)
		par := parser.New(lex)
		program := par.ParseProgram()

		if len(par.Errors()) > 0 {
			for _, errorMessage := range par.Errors() {
				io.WriteString(out, errorMessage)
			}
			continue
		}

		io.WriteString(out, fmt.Sprintf("Program has %d statement nodes", len(program.Statements)))
		io.WriteString(out, "\n")
		io.WriteString(out, fmt.Sprintf("Program back to string: %s", program.String()))
	}
}
