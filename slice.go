/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package iter

import "golang.org/x/exp/slices"

// SliceIterator enables iteration over slices.
type SliceIterator[T any] struct {
    base []T
    idx  int
}

func (s *SliceIterator[T]) Next() (out T, ok bool) {
    if s.EstimatedRemaining() > 0 {
        out = s.base[s.idx]
        s.idx++
        return out, true
    }

    return
}

func (s *SliceIterator[T]) Advance(n int) int {
    adv := min(n, s.EstimatedRemaining())
    s.idx += adv
    return adv
}

func (s *SliceIterator[T]) EstimatedRemaining() int {
    return len(s.base) - s.idx
}

func (s *SliceIterator[T]) Collect() []T {
    return slices.Clone(s.base[s.idx:])
}

func (s *SliceIterator[T]) Reset() {
    s.idx = 0
}

// FromSlice creates a new SliceIterator from a slice.
func FromSlice[T any](s []T) Iterator[T] {
    return &SliceIterator[T]{
        base: s,
        idx:  0,
    }
}
