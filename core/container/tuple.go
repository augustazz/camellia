package container

import "fmt"

type TwoTuple[T any] struct {
	First, Second T
}

type ThirdTuple[T any] struct {
	First, Second, Third T
}

func NewTwoTuple[T any](first, second T) *TwoTuple[T] {
	return &TwoTuple[T]{first, second}
}

func (t *TwoTuple[T]) String() string {
	return fmt.Sprintf("Tuple[%v, %v]", t.First, t.Second)
}

func NewThirdTuple[T any](first, second, third T) *ThirdTuple[T] {
	return &ThirdTuple[T]{first, second, third}
}

func (t *ThirdTuple[T]) String() string {
	return fmt.Sprintf("Tuple[%v, %v, %v]", t.First, t.Second, t.Third)
}
