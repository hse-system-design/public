package ordcol

type IterationOrder int

const (
	ByInsertion IterationOrder = 1
	ByKey IterationOrder = 2
)

type Iterator interface {
	HasNext() bool
	Next() (Item, error)
}

type Item interface {
	Key() int
	Value() int
}

type Collection interface {
	Add(item Item)

	IterateBy(order IterationOrder) Iterator
	At(key int) (Item, error)
}