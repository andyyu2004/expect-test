package expect

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type patches struct {
	text    []byte
	patches []patch
}

func (c *patches) apply(p patch) {
	c.patches = append(c.patches, p)
	slices.SortFunc(c.patches, func(i, j patch) int {
		return i.loc.start - j.loc.start
	})

	// need to shift the location of patch accordingly due to the previous patches
	var deleted, inserted int
	for _, patch := range c.patches {
		if patch.loc.start >= p.loc.start {
			break
		}
		deleted += patch.loc.len()
		inserted += len(patch.with)
	}

	c.text = replace(c.text, p.loc.shifted(inserted-deleted), p.with)
}

func newPatches(text []byte) patches {
	return patches{text: text}
}

type patch struct {
	loc  location
	with string
}

func replace(bytes []byte, loc location, with string) []byte {
	if loc.len() == len(with) {
		copy(bytes[loc.start:loc.end], with)
		return bytes
	} else if loc.len() > len(with) {
		copy(bytes[loc.start:], with)
		copy(bytes[loc.start+len(with):], bytes[loc.end:])
		return bytes[:len(bytes)-loc.len()+len(with)]
	} else {
		// loc.len() < len(with)
		bytes = append(bytes, make([]byte, len(with)-loc.len())...)
		copy(bytes[loc.end+len(with)-loc.len():], bytes[loc.end:])
		copy(bytes[loc.start:], with)
		return bytes
	}
}

func locate(t testing.TB, text string, line int) (location, rune) {
	scanner := bufio.NewScanner(strings.NewReader(text))
	startIdx := 0
	for i := 1; i < line; i++ {
		require.True(t, scanner.Scan())
		// + 1 for the newline (sorry windows)
		startIdx += 1 + len(scanner.Text())
	}

	require.True(t, scanner.Scan())

	col := strings.Index(scanner.Text(), "Expect")
	if col == -1 {
		panic("no `Expect` call found on expected line")
	}

	// we're looking for the start of the string literal
	// this slice will be of the pattern `Expect(t, <what we want>) ... the rest of the code`
	sliceToParse := text[startIdx+col:]
	endIdx := 0
	var expr ast.Expr
	for expr == nil {
		idx := strings.Index(sliceToParse[endIdx:], ")")
		if idx == -1 {
			panic("no closing paren found, this code must be syntactically correct otherwise this wouldn't be running")
		}
		endIdx += idx + 1

		expr, _ = parser.ParseExpr(sliceToParse[:endIdx])
	}

	call := expr.(*ast.CallExpr)
	if len(call.Args) != 1 {
		panic(fmt.Sprintf("expected 1 args to Expect, got %d", len(call.Args)))
	}

	expectedLit := call.Args[0].(*ast.BasicLit)

	end := startIdx + col + endIdx - 2
	start := end - len(expectedLit.Value) + 2
	delimiter := expectedLit.Value[0]
	if text[start-1] != delimiter || text[end] != delimiter || text[start:end] != expectedLit.Value[1:len(expectedLit.Value)-1] {
		panic("expected to find string literal")
	}

	return location{start, end}, rune(delimiter)
}
