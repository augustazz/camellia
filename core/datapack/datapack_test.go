package datapack

import (
	"fmt"
	"testing"
)

func TestUnPack(t *testing.T) {

	b := []rune("qwertyuiop")
	f(b[:5])
	f(b[5:])


}

func f(b []rune) {
	for _, bb := range b {
		fmt.Print(string(bb), " ")
	}
	fmt.Println("\n=============")
}
