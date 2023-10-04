package example

import (
	"testing"

	"github.com/andyyu2004/expect-test"
)

func TestExample(t *testing.T) {
	expect.Expect(`foo`)
	expect.Expect(``)
	expect.Expect("double quoted string that can have ` in it")
	expect.Expect("escaped double quoted \"hi\" string")
	expect.Expect(`backslash \ ignored in raw\`)
	expect.Expect(`some
multiline
string`)

	check := func(input string, expectation expect.Expectation) {
		expectation.AssertEqual(t, input)
	}

	check(`test that the locate still works when there 
is an unrelated backtick on the same line`, expect.Expect(`expected`))
}
