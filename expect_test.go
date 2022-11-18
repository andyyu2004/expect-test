package expect

import (
	"os"
	"testing"
)

func TestExpect(t *testing.T) {
	t.Run("expect ok", func(t *testing.T) {
		exp := Expect(`foo`)
		exp.Expect(t, `foo`)
	})

	t.Run("expect fail", func(t *testing.T) {
		exp := Expect(``)
		os.Setenv(env, "1")
		exp.Expect(t, `bar`)
		panic("")
	})
}
