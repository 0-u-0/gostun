package test

import (
	"testing"

	"git.learning-tech.cn/realtimecat/gostun/libs"
)

func Test_rand_bytes(t *testing.T) {
	t.Logf("randbyte test : %x",libs.RandBytes(10))
}