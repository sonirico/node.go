package repl

import (
	"bufio"
	"fmt"
	"io"
	"node.go/lexer"
	"node.go/token"
)

const PROMPT = "\n/> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		io.WriteString(out, PROMPT)
		_ = scanner.Scan()

		line := scanner.Text()
		lex := lexer.New(line)
		mustQuit := false

		for !mustQuit {
			t := lex.NextToken()
			io.WriteString(out, fmt.Sprintln(t))
			mustQuit = t.Type == token.EOF
		}
	}
}
