package expect

import (
	"testing"
)

func TestExpect(t *testing.T) {
	t.Run("expect ok", func(t *testing.T) {
		exp := Expect(`foo`)
		exp.AssertEqual(t, `foo`)
	})

	Expect(`bar`).AssertEqual(t, `bar`)
	Expect(`bar`).AssertEqual(t, `bar`)
}

func TestMultipleMultilineUpdates(t *testing.T) {
	// actually "running" this test is a bit awkward
	// I just replace the expectation with an empty string and update expects to check it works
	Expect(`foo
bar
baz`).AssertEqual(t, `foo
bar
baz`)

	Expect(`a
b
c`).AssertEqual(t, `a
b
c`)
}
