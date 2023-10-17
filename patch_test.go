package expect

import (
	"testing"

	_ "embed"

	"github.com/stretchr/testify/require"
)

//go:embed testdata/example_test.go
var example string

func TestLocate(t *testing.T) {
	check := func(t *testing.T, line int, expected string, expectedDelimeter rune) {
		t.Helper()
		location, delimeter := locate(t, example, line)
		require.Equal(t, expectedDelimeter, delimeter)
		require.Equal(t, expected, example[location.start:location.end])
	}

	check(t, 10, `foo`, '`')
	check(t, 11, ``, '`')
	check(t, 12, "double quoted string that can have ` in it", '"')
	check(t, 13, "escaped double quoted \\\"hi\\\" string", '"')
	check(t, 14, `backslash \ ignored in raw\`, '`')
	check(t, 15, `some
multiline
string`, '`')
	check(t, 24, `expected`, '`')
	check(t, 25, `expected2`, '"')
}

func TestPatch(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		p := newPatches([]byte(`foobar`))
		p.apply(patch{location{1, 3}, `baz`})
		require.Equal(t, `fbazbar`, string(p.text))
	})

	t.Run("sequenced", func(t *testing.T) {
		p := newPatches([]byte(`foo
bar`))

		p.apply(patch{location{0, 3}, `foo2`})
		require.Equal(t, `foo2
bar`, string(p.text))

		p.apply(patch{location{4, 7}, `bar2`})
		require.Equal(t, `foo2
bar2`, string(p.text))
	})

	t.Run("multiline", func(t *testing.T) {
		p := newPatches([]byte(``))

		p.apply(patch{location{0, 0}, `foo
bar`})
		require.Equal(t, `foo
bar`, string(p.text))
	})
}

func TestReplace(t *testing.T) {
	check := func(t *testing.T, text string, loc location, with string, expected string) {
		t.Helper()
		bytes := []byte(text)
		require.Equal(t, expected, string(replace(bytes, loc, with)))
	}

	check(t, `foobar`, location{0, 3}, `baz`, `bazbar`)

	check(t, `foobar`, location{0, 4}, `baz`, `bazar`)

	check(t, `foobar`, location{1, 3}, `baz`, `fbazbar`)
	check(t, `foobar`, location{0, 0}, `baz`, `bazfoobar`)
}
