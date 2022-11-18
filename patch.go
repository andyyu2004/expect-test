package expect

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type patch struct{}

func locate(t testing.TB, text string, line int) location {
	scanner := bufio.NewScanner(strings.NewReader(text))
	start := 0
	for i := 1; i < line; i++ {
		require.True(t, scanner.Scan())
		// + 1 for the newline (sorry windows)
		start += 1 + len(scanner.Text())
	}

	require.True(t, scanner.Scan())

	startColumn := 0
	for i, char := range scanner.Text() {
		if char == '`' {
			startColumn = i + 1
			break
		}
	}

	if startColumn == 0 {
		require.FailNow(t, "no start marker found")
	}

	start += startColumn

	for j, char := range text[start:] {
		if char == '`' {
			return location{start, start + j}
		}
	}

	require.FailNow(t, "no end marker found")
	return location{}
}
