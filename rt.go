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
	patch    patch
}

func newfilert(t testing.TB, exp Expectation) *filert {
	content, err := os.ReadFile(exp.file)
	require.NoError(t, err)
	return &filert{exp, string(content), patch{}}
}

func (rt *filert) update(t testing.TB, actual string) {
	locate(t, rt.original, rt.exp.line)
}

type location struct {
	start int
	end   int
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
