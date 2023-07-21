package acquire

import (
	"acquire/internal/util"
	"fmt"
)

// Collection
// a collection is like a pile of stock certificates, a bunch of tiles, or a stack of paper money.
// it contains a fixed amount of identical items, which can be taken from other collections.
type Collection[T comparable] struct {
	Items []T
}

func newCollection[T comparable]() *Collection[T] {
	return &Collection[T]{
		Items: make([]T, 0, 0),
	}
}

func (c *Collection[T]) clone() *Collection[T] {
	return &Collection[T]{
		// must be cloned (but can't know if we need to clone deeper automatically)
		Items: util.Clone(c.Items),
	}
}

func (c *Collection[T]) add(n int, ctor func() T) {
	for i := 0; i < n; i++ {
		c.Items = append(c.Items, ctor())
	}
}

func (c *Collection[T]) remove(val T) {
	index, ok := c.indexOf(val)
	if !ok {
		return
	}

	c.Items = append(c.Items[:index], c.Items[index+1:]...)
}

// removeN
// removes n elements from the end of the collection
func (c *Collection[T]) removeN(amount int) {
	l := len(c.Items)
	clampedAmount := util.Min(l, amount)
	c.Items = c.Items[:l-clampedAmount]
}

func (c *Collection[T]) indexOf(val T) (int, bool) {
	return util.IndexOf(c.Items, val)
}

func (c *Collection[T]) take(other *Collection[T]) error {
	if len(other.Items) < 1 {
		return fmt.Errorf("cannot take from an empty collection")
	}

	item := other.Items[0]

	// remove item from other
	other.Items = append(other.Items[:0], other.Items[0+1:]...)

	// add it to this one
	c.Items = append(c.Items, item)

	return nil
}
