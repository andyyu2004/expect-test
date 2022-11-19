package expect

import (
	"os"
	"runtime"
	"testing"

	"github.com/kr/pretty"
	"github.com/stretchr/testify/require"
)

type Expectation struct {
	file     string
	line     int
	expected string
}

const env = "UPDATE_SNAPSHOTS"

func should_update_expect() bool {
	_, ok := os.LookupEnv(env)
	return ok
}

func (exp Expectation) Expect(t testing.TB, actual string) {
	if exp.expected == actual {
		return
	}

	if !should_update_expect() {
		require.Equal(t, exp.expected, actual)
	}

	exp.update(t, actual)
}

func (exp Expectation) update(t testing.TB, actual string) {
	file, err := os.Open(exp.file)
	require.NoError(t, err)
	defer file.Close()
	rt := getrt()
	rt.update(t, exp, actual)

}

func Expect(expected any) Expectation {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		panic("failed to get caller for `expect.Expect`")
	}

	var formatted string
	switch expected := expected.(type) {
	case string:
		formatted = expected
	default:
		formatted = pretty.Sprintf("%# v", expected)
	}
	return Expectation{file, line, formatted}
}
