package expect

import (
	"testing"

	_ "embed"

	"github.com/stretchr/testify/require"
)

//go:embed testdata/example_test.go
var example string

func TestLocate(t *testing.T) {
	check := func(t *testing.T, line int, expected string) {
		t.Helper()
		location := locate(t, example, line)
		require.Equal(t, expected, example[location.start:location.end])
	}

	check(t, 10, `foo`)
	check(t, 11, ``)
	check(t, 12, `some
multiline
string`)
}
