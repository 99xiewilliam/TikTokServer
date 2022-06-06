package util

import (
	"fmt"
	"testing"
)

func TestSha(t *testing.T) {
	s := GenSaltPassword("1", "123456")
	fmt.Println(s)
	fmt.Println(len(s))
}
