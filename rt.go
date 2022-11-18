package expect

import (
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

var rt once[globalrt]

func getrt() *globalrt {
	return rt.Get(func() globalrt {
		return globalrt{filerts: make(map[string]*filert)}
	})
}

type globalrt struct {
	sync.RWMutex
	filerts map[string]*filert
}

func (rt *globalrt) update(t testing.TB, exp Expectation, actual string) {
	rt.Lock()
	defer rt.Unlock()
	ft, ok := rt.filerts[exp.file]
	if !ok {
		ft = newfilert(t, exp)
		rt.filerts[exp.file] = ft
	}
	ft.update(t, actual)
}

type filert struct {
	exp      Expectation
	original string
	patches  patches
}

func newfilert(t testing.TB, exp Expectation) *filert {
	content, err := os.ReadFile(exp.file)
	require.NoError(t, err)
	return &filert{exp, string(content), newPatches(content)}
}

func (rt *filert) update(t testing.TB, actual string) {
	loc := locate(t, rt.original, rt.exp.line)
	rt.patches.apply(patch{loc, actual})
	require.NoError(t, os.WriteFile(rt.exp.file, rt.patches.text, 0))
}

type location struct {
	start int
	end   int
}

func (loc location) len() int {
	return loc.end - loc.start
}

func (loc location) shifted(k int) location {
	return location{loc.start + k, loc.end + k}
}

type once[T any] struct {
	once  sync.Once
	value T
}

func (l *once[T]) Get(f func() T) *T {
	l.once.Do(func() {
		l.value = f()
	})
	return &l.value
}
