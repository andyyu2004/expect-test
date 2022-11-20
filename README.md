# Expect Test


```go
import "github.com/andyyu2004/expect-test"

func TestExample(t *testing.T) {
    expectation := expect.Expect(`expected result here, automatically update by setting 'UPDATE_SNAPSHOTS=1'`)
    actual := "..."
    expectation.AssertEqual(t, actual)
}
```

This works with either raw strings or double-quoted string. Raw strings have the benefit of supporting multiline but have the issue where it can't contain more backticks.
