package container

type TwoTuple struct {
	First, Second interface{}
}

func NewTwoTuple(first, second interface{}) *TwoTuple {
	return &TwoTuple {first, second}
}
