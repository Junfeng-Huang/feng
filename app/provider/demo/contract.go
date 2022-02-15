package demo

const Key = "feng:demo"

type Service interface {
	GetFoo() Foo
}

type Foo struct {
	Name string
}
