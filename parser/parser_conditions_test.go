package parser

import (
	"testing"
)

func TestElseBlockStatement(t *testing.T) {
	// Ref: https://github.com/sonirico/node.go/issues/5

	// Ensures that a conditional expression within other block statements
	// are correctly parsed
	payload := `
		if (true) {
			if (true) {
			} else {
			}
		}
	`
	ParseTesting(t, payload)
}
