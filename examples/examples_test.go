package examples

import (
	"testing"
	"github.com/facebookgo/ensure"
)

func TestRoundtrip(t *testing.T) {
	b := []byte{1, 2, 3, 4, 5}
	rt := cccRoundtrip(b)
	ensure.DeepEqual(t, b, rt)
}
