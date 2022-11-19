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
	Expect(``).AssertEqual(t, `foo
bar
baz`)

	Expect(``).AssertEqual(t, `foo
bar
baz`)
}
