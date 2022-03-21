# assert

`assert` is an opinionated test assertion libary. It is not intended to fulfill every use case and is not intended to be flexible/extendable.

## Installing

Because some of the assertions use type parameters, Go 1.18+ is required. This library can be installed with

```sh
go get -u github.com/mattmeyers/assert
```

## Usage

The assertions found within this library can replace any simple assertion logic normally found in tests. All assertions are marked as `t.Helper`s so stack traces will point to the appropriate line in the test. Additionally, all assertions are not fatal by default. To convert an assertion to a fatal assertion, the `*testing.T` struct can be wrapped with `Fatal()`.

For example, consider the following function that returns a stringified JSON array.

```go
// In foo.go
package foo

func getJSON() string {
	return "[1,2,3,4,5]"
}
```

Then this test will assert that the JSON contains a certain value.

```go
// In foo_test.go
package foo

import (
	"encoding/json"
	"testing"

	"github.com/mattmeyers/assert"
)

func TestGetJSONContains3(t *testing.T) {
	str := getJSON()

	var values []int
	err := json.Unmarshal([]byte(str), &values)
    // Will fatally fail test if JSON is malformed
	assert.NoError(assert.Fatal(t), err) 

	assert.SliceContains(t, values, 3)
    // Will fail test
	// assert.SliceContains(t, values, 9)}

```