package container

import "testing"

func TestNewThirdTuple(t *testing.T) {
	tt := NewThirdTuple(30, 50, 90)
	a := tt.Third
	println(tt.String())
	println()
	println(a)
}
