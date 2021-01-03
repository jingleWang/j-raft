package transport

import (
	"bytes"
	"fmt"
	"math"
	"testing"
)

func TestNewPack(t *testing.T) {
	pack := NewPack(math.MaxUint64, []byte("hello world"), RESP)
	fmt.Println(pack.IsRequest())

	buf := pack.Encode()
	fmt.Println(buf)

	pack1 := Decode(bytes.NewReader(buf))
	fmt.Println(pack1.IsRequest())
}
