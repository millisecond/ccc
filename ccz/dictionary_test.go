package ccz

import (
	"testing"
	"github.com/facebookgo/ensure"
)

func TestAppend(t *testing.T) {
	ensure.DeepEqual(t, Append([]byte{1, 2}, []byte{3, 4}), []byte{1, 2, 3, 4})
}

func TestLatestVersion(t *testing.T) {
	_, err := latestVersion("garbage.sdlfnsdkjfhsdkj")
	ensure.NotNil(t, err)
}
