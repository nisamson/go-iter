/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package iter

func makeIter[T any]() ChanIterator[T] {
    return make(chan T, 1)
}

// FromChannel creates a new iterator which iterates over the channel `c`
func FromChannel[T any](c chan T) Iterator[T] {
    return ChanIterator[T](c)
}

func (c ChanIterator[T]) Next() (T, bool) {
    o, ok := <-c
    return o, ok
}

func (c ChanIterator[T]) Advance(n int) int {
    for i := 0; i < n; i++ {
        _, ok := c.Next()
        if !ok {
            return i
        }
    }
    return n
}

func (c ChanIterator[T]) EstimatedRemaining() int {
    return len(c)
}

func (c ChanIterator[T]) Collect() []T {
    var out []T
    for item := range c {
        out = append(out, item)
    }
    return out
}
