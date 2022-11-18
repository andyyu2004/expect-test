package expect

import (
	"testing"
)

func TestExpect(t *testing.T) {
	t.Run("expect ok", func(t *testing.T) {
		exp := Expect(`foo`)
		exp.Expect(t, `foo`)
	})
}
