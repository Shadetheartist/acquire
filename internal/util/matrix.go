package util

import (
	"errors"
	"fmt"
)

type Matrix[T comparable] struct {
	sizeX int
	sizeY int
	area  int
	data  []T
}

func NewMatrix[T comparable](sizeX int, sizeY int) *Matrix[T] {
	matrix := Matrix[T]{}

	matrix.sizeX = sizeX
	matrix.sizeY = sizeY
	matrix.area = sizeX * sizeY
	matrix.data = make([]T, matrix.area)

	return &matrix
}

func (m *Matrix[T]) Copy() *Matrix[T] {
	newMatrix := NewMatrix[T](m.sizeX, m.sizeY)

	copy(newMatrix.data, m.data)

	return newMatrix
}

func (m *Matrix[T]) Index(x int, y int) int {
	return (y * m.sizeX) + x
}

func (m *Matrix[T]) IndexPt(pt Point[int]) int {
	return (pt.Y * m.sizeX) + pt.X
}

func (m *Matrix[T]) Get(x int, y int) T {
	i := m.Index(x, y)

	return m.data[i]
}

func (m *Matrix[T]) GetPt(pt Point[int]) (T, error) {

	if !m.IsInBounds(pt.X, pt.Y) {
		var noop T
		return noop, errors.New("point not within bounds")
	}

	i := m.Index(pt.X, pt.Y)

	return m.data[i], nil
}

func (m *Matrix[T]) GetI(index int) T {
	return m.data[index]
}

func (m *Matrix[T]) GetNeighbors(pt Point[int]) []T {
	neighbors := make([]T, len(Directions))

	for d, npt := range pt.OrthogonalNeighbours() {
		nb, err := m.GetPt(npt)
		if err == nil {
			neighbors[d] = nb
		}
	}

	return neighbors
}

func (m *Matrix[T]) Set(x int, y int, d T) {
	i := m.Index(x, y)
	m.data[i] = d
}

func (m *Matrix[T]) SetPt(pt Point[int], d T) error {
	if !m.IsInBounds(pt.X, pt.Y) {
		return errors.New("point not within bounds")
	}

	i := m.Index(pt.X, pt.Y)

	m.data[i] = d

	return nil
}

func (m *Matrix[T]) SetI(i int, d T) {
	m.data[i] = d
}

func (m *Matrix[T]) Equal(m2 *Matrix[T]) bool {
	for y := 0; y < m.sizeY; y++ {
		for x := 0; x < m.sizeX; x++ {
			index := m.Index(x, y)
			if m.data[index] != m2.data[index] {
				return false
			}
		}
	}

	return true
}

func (m *Matrix[T]) Print() {
	fmt.Println()

	for y := 0; y < m.sizeY; y++ {
		for x := 0; x < m.sizeX; x++ {
			d := m.Get(x, y)
			fmt.Print(d, "\t")
		}
		fmt.Println()
	}

	fmt.Println()
}

type MatrixIterator[T comparable] func(rt T, x int, y int, idx int)

func (m *Matrix[T]) Iterate(iter MatrixIterator[T]) {
	for y := 0; y < m.sizeY; y++ {
		for x := 0; x < m.sizeX; x++ {
			d := m.Get(x, y)
			idx := m.Index(x, y)
			iter(d, x, y, idx)
		}
	}
}

type MatrixIteratorFind[T comparable] func(rt T, x int, y int, idx int) bool

func (m *Matrix[T]) Find(iter MatrixIteratorFind[T]) T {
	for y := 0; y < m.sizeY; y++ {
		for x := 0; x < m.sizeX; x++ {
			d := m.Get(x, y)
			idx := m.Index(x, y)
			if iter(d, x, y, idx) {
				return d
			}
		}
	}

	var z T
	return z
}

func (m *Matrix[T]) Area() int {
	return m.area
}

func (m *Matrix[T]) Len() int {
	return len(m.data)
}

func (m *Matrix[T]) IsInBounds(x int, y int) bool {
	return x >= 0 && x < m.sizeX && y >= 0 && y < m.sizeY
}

func (m *Matrix[T]) IsIndexInBounds(i int) bool {
	return i >= 0 && i < len(m.data)
}
