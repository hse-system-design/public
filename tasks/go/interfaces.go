package ordcol

import "fmt"

var ErrDuplicateKey = fmt.Errorf("duplicate key")
var ErrEmptyIterator = fmt.Errorf("empty iterator")

type IterationOrder int

const (
	ByInsertion IterationOrder = 1
	ByKey IterationOrder = 2
)

type Iterator interface {
	HasNext() bool
	// Next returns ErrEmptyIterator if there is no more elements (HasNext() is false)
	Next() (Item, error)
}

type Item interface {
	Key() int
	Value() int
}

type Collection interface {
	// Add returns ErrDuplicateKey if there already is an element with the same Item.Key in collection
	Add(item Item) error

	IterateBy(order IterationOrder) Iterator
	At(key int) (Item, bool)
}