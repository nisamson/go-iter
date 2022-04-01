/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package iter

import (
    "golang.org/x/exp/constraints"
)

type mapIter[T, U any] struct {
    i      Iterator[T]
    mapper func(T) U
}

func (m mapIter[T, U]) Next() (u U, ok bool) {
    if item, ok := m.i.Next(); ok {
        return m.mapper(item), ok
    }
    return
}

func (m mapIter[T, U]) Advance(n int) int {
    return m.i.Advance(n)
}

func (m mapIter[T, U]) EstimatedRemaining() int {
    return m.i.EstimatedRemaining()
}

func (m mapIter[T, U]) Collect() (out []U) {
    for item, ok := m.i.Next(); ok; item, ok = m.i.Next() {
        out = append(out, m.mapper(item))
    }
    return out
}

type filterIter[T any] struct {
    i      Iterator[T]
    filter func(T) bool
}

func (f filterIter[T]) Next() (t T, ok bool) {
    for t, ok = f.i.Next(); ok; t, ok = f.i.Next() {
        if f.filter(t) {
            return
        }
    }
    return t, false
}

func (f filterIter[T]) Advance(n int) int {
    for i := 0; i < n; i++ {
        _, ok := f.Next()
        if !ok {
            return i
        }
    }
    return n
}

func (f filterIter[T]) EstimatedRemaining() int {
    return 0 // must be zero because we don't know if any elements will pass the predicate
}

func (f filterIter[T]) Collect() (out []T) {
    for t, ok := f.Next(); ok; t, ok = f.Next() {
        out = append(out, t)
    }
    return out
}

// Map applies a transformation to all elements and yields those transformed elements.
func Map[T, U any](i Iterator[T], mapper func(T) U) Iterator[U] {
    return mapIter[T, U]{
        i:      i,
        mapper: mapper,
    }
}

// Filter wraps an iterator, skipping elements for which filter returns false.
func Filter[T any](i Iterator[T], filter func(T) bool) Iterator[T] {
    return filterIter[T]{
        i:      i,
        filter: filter,
    }
}

// Fold accumulates a value by applying the step function to each element.
func Fold[Acc, T any](i Iterator[T], acc Acc, step func(Acc, T) Acc) Acc {
    for next, ok := i.Next(); ok; next, ok = i.Next() {
        acc = step(acc, next)
    }
    return acc
}

// Reduce accumulates a value by applying the step function to the accumulated value and the next element.
// This is a specialized version of Fold, and returns a boolean value which is false if the collection was empty.
func Reduce[T any](i Iterator[T], step func(T, T) T) (acc T, ok bool) {
    acc, ok = i.Next()
    if !ok {
        return
    }

    return Fold(i, acc, step), true
}

// Count returns the number of elements in the iterator. This consumes the iterator.
func Count[T any](i Iterator[T]) int {
    return Fold(i, 0, func(sz int, item T) int { sz++; return sz })
}

// Sum returns the sum of elements in the iterator. This consumes the iterator.
func Sum[T constraints.Integer | constraints.Float | constraints.Complex](i Iterator[T]) T {
    return Fold(i, 0, func(a T, b T) T { return a + b })
}

// Max returns the maximum value and true if the collection was not empty, and zero and false otherwise.
func Max[T constraints.Ordered](i Iterator[T]) (T, bool) {
    return Reduce(i, func(a, b T) T {
        if a > b {
            return a
        } else {
            return b
        }
    })
}

// Min returns the minimum value and true if the collection was not empty, and zero and false otherwise.
func Min[T constraints.Ordered](i Iterator[T]) (T, bool) {
    return Reduce(i, func(a, b T) T {
        if a < b {
            return a
        } else {
            return b
        }
    })
}

// MinBy returns the minimum value using the given comparison function and true if the collection was not empty,
// and zero and false otherwise.
func MinBy[T any](i Iterator[T], less func(a, b T) bool) (T, bool) {
    return Reduce(i, func(a, b T) T {
        if less(a, b) {
            return a
        } else {
            return b
        }
    })
}

// MaxBy returns the maximum value using the given comparison function and true if the collection was not empty,
// and zero and false otherwise.
func MaxBy[T any](i Iterator[T], less func(a, b T) bool) (T, bool) {
    return Reduce(i, func(a, b T) T {
        if less(a, b) {
            return b
        } else {
            return a
        }
    })
}

// Last returns the last value of this iterator, if any exists. If empty, Last returns zero and false.
func Last[T any](i Iterator[T]) (t T, ok bool) {
    var atLeastOne bool
    var lastOk T
    if i.EstimatedRemaining() > 1 {
        i.Advance(i.EstimatedRemaining() - 1)
    }
    for t, ok = i.Next(); ok; t, ok = i.Next() {
        lastOk = t
        atLeastOne = true
    }
    return lastOk, atLeastOne
}

// Empty creates an iterator which contains no values.
func Empty[T any]() Iterator[T] {
    return FromSlice[T]([]T{})
}
