/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package iter

type (
    // Iterator represents a one-time one-by-one view of elements through some internal state.
    Iterator[T any] interface {
        // Next returns the next element, or a zero element and false if there are no more elements
        Next() (T, bool)

        // Advance moves the iterator forward by n elements, returning how much it was actually able to progress.
        Advance(n int) int

        // EstimatedRemaining returns a lower bound on the number of elements remaining,
        // i.e. there are at least Remaining() elements remaining.
        // WARNING: This is a *lower* bound! There may be more elements, even if this returns zero!
        EstimatedRemaining() int

        // Collect returns a slice containing the elements of the iterator.
        Collect() []T
    }

    // ChanIterator is an iterator backed by a channel.
    ChanIterator[T any] chan T

    // Iterable represents an object which can be iterated over.
    Iterable[T any] interface {
        Iter() Iterator[T]
    }
)
