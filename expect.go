package expect

import (
	"os"
	"runtime"
	"testing"

	"github.com/kr/pretty"
	"github.com/stretchr/testify/assert"
)

type Expectation struct {
	file     string
	line     int
	expected string
}

func should_update_expect() bool {
	for _, env := range []string{"UPDATE_SNAPS", "UPDATE_SNAPSHOTS", "UPDATE_EXPECT"} {
		if _, ok := os.LookupEnv(env); ok {
			return true
		}
	}

	return false
}

func (exp Expectation) AssertEqual(t testing.TB, actual any) {
	t.Helper()
	var formatted string
	switch actual := actual.(type) {
	case string:
		formatted = actual
	default:
		formatted = pretty.Sprintf("%# v", actual)
	}

	if exp.expected == formatted {
		return
	}

	if !should_update_expect() {
		assert.Equal(t, exp.expected, formatted)
		return
	}

	exp.update(t, formatted)
}

func (exp Expectation) update(t testing.TB, actual string) {
	getrt().update(t, exp, actual)
}

func Expect(expected string) Expectation {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		panic("failed to get caller for `expect.Expect`")
	}

	return Expectation{file, line, expected}
}
