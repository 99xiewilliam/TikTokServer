package middleware

import (
	"fmt"
	"testing"
)

func TestToken(t *testing.T) {
	token, err := GenToken(1004)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s\n", token)
	claims, ok := parseToken(token)
	if ok {
		fmt.Printf("%+v\n", claims)
	}
}
