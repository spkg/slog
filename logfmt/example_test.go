package logfmt_test

import (
	"fmt"
	"time"

	"github.com/spkg/slog/logfmt"
)

func Example() {
	// create a buffer and ensure that it's
	// internal memory buffer is released when no
	// longer needed
	buf := logfmt.Buffer{}
	defer buf.Reset()

	buf.WriteTimestamp(time.Unix(1234567890, 987654321).UTC())
	buf.WriteKey("info")
	buf.WriteProperty("key1", 1)
	buf.WriteProperty("key2", "value 2")

	fmt.Println(buf.String())
	// Output: 2009-02-13T23:31:30.987654+0000 info key1=1 key2="value 2"
}
