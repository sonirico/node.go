package repl

import (
	"bufio"
	"io"
	"node.go/lexer"
	"node.go/parser"
)

const PROMPT = "{o_o}-> "

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
				io.WriteString(out, "\n")
			}
			io.WriteString(out, "\n")
			continue
		}

		io.WriteString(out, "\n")
		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}
