package example

import (
	"testing"

	"github.com/andyyu2004/expect-test"
)

func TestExample(t *testing.T) {
	expect.Expect(t, `foo`)
	expect.Expect(t, ``)
	expect.Expect(t, "double quoted string that can have ` in it")
	expect.Expect(t, "escaped double quoted \"hi\" string")
	expect.Expect(t, `backslash \ ignored in raw\`)
	expect.Expect(t, `some
multiline
string`)
}
