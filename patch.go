package expect

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
)

type patches struct {
	text    []byte
	patches []patch
}

func (c *patches) apply(p patch) {
	c.patches = append(c.patches, p)
	slices.SortFunc(c.patches, func(i, j patch) bool {
		return i.loc.start < j.loc.start
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
	start := 0
	for i := 1; i < line; i++ {
		require.True(t, scanner.Scan())
		// + 1 for the newline (sorry windows)
		start += 1 + len(scanner.Text())
	}

	require.True(t, scanner.Scan())

	startColumn := 0
	delimiter := '`'
	for i, char := range scanner.Bytes() {
		if char == '`' {
			startColumn = i + 1
			break
		} else if char == '"' {
			startColumn = i + 1
			delimiter = '"'
			break
		}
	}

	if startColumn == 0 {
		require.FailNow(t, "no start marker found")
	}

	start += startColumn

	for j, char := range text[start:] {
		if char == delimiter {
			return location{start, start + j}, delimiter
		}
	}

	require.FailNow(t, "no end marker found")
	return location{}, delimiter
}
