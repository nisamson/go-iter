package iter

import (
    "testing"

    "golang.org/x/exp/rand"
    "golang.org/x/exp/slices"
)

var sampleInts = []int{
    0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
}
var isEven = func(i int) bool {
    return i%2 == 0
}
var double = func(i int) int {
    return i * 2
}

func shuffledInts() []int {
    myInts := slices.Clone(sampleInts)
    rand.Shuffle(len(myInts), func(i, j int) { tmp := myInts[i]; myInts[i] = myInts[j]; myInts[j] = tmp })
    return myInts
}

func TestCount(t *testing.T) {
    ints := shuffledInts()
    iter := FromSlice(ints)
    count := Count[int](iter)
    if count != len(ints) {
        t.Fatalf("expected %d, got %d", len(ints), count)
    }
}

func TestSum(t *testing.T) {
    iter := FromSlice(shuffledInts())
    sum := Sum(iter)
    if sum != 45 {
        t.Fatalf("expected 45, got %d", sum)
    }
}

func TestFilter(t *testing.T) {

    iter := FromSlice(shuffledInts())
    filtered := Filter(iter, isEven)
    sum := Sum(filtered)
    if sum != 20 {
        t.Fatalf("expected 20, %d", sum)
    }
}

func TestMap(t *testing.T) {

    iter := FromSlice(shuffledInts())
    mapped := Map(iter, double)
    sum := Sum(mapped)
    if sum != 90 {
        t.Fatalf("expected 90, got %d", sum)
    }
}

func TestLast(t *testing.T) {
    ints := shuffledInts()
    last := ints[len(ints)-1]
    iter := FromSlice(ints)
    testLast, ok := Last(iter)

    if !ok {
        t.Fail()
        t.Errorf("expected a result")
    }

    if testLast != last {
        t.Fail()
        t.Errorf("expected %d, got %d", last, testLast)
    }

    c := make(chan int, len(ints))
    for _, i := range ints {
        c <- i
    }
    close(c)

    iter = FromChannel(c)
    testLast, ok = Last(iter)

    if !ok {
        t.Fail()
        t.Errorf("expected a result")
    }

    if testLast != last {
        t.Fail()
        t.Errorf("expected %d, got %d", last, testLast)
    }

    _, ok = Last(iter)
    if ok {
        t.Fail()
        t.Errorf("expected false from empty iterator")
    }
}

func TestFilterIter_Collect(t *testing.T) {
    ints := shuffledInts()
    iter := Filter(FromSlice(ints), isEven)
    evenInts := iter.Collect()
    if len(evenInts) != len(ints)/2 {
        t.Fatalf("Expected evenInts to be exactly half the size of ints, got %d", len(evenInts))
    }
}

func TestMapIter_Collect(t *testing.T) {
    ints := shuffledInts()
    doubledInts := Map(FromSlice(ints), double).Collect()
    var sum int
    for _, i := range doubledInts {
        sum += i
    }
    if sum != 90 {
        t.Fatalf("expected %d, got %d", 90, sum)
    }
}

func TestComprehensive(t *testing.T) {
    data := []string{
        "hello",
        "world",
        "this",
        "world",
        "is",
        "ted",
        "shiny buttons",
    }

    iter := FromSlice(data)
    max, ok := Max(iter)
    if !ok || max != "world" {
        t.Errorf("Expected 'world', got '%s' or nothing", max)
    }

    iter = FromSlice(data)
    min, ok := Min(iter)
    if !ok || min != "hello" {
        t.Errorf("Expected 'hello', got '%s' or nothing", min)
    }

    iter = FromSlice(data)
    max, ok = MaxBy(iter, func(a, b string) bool {
        return len(a) < len(b)
    })

    if !ok || max != "shiny buttons" {
        t.Errorf("Expected 'shiny buttons', got '%s' or nothing", max)
    }

    iter = FromSlice(data)
    min, ok = MinBy(iter, func(a, b string) bool {
        return len(a) < len(b)
    })

    if !ok || min != "is" {
        t.Errorf("Expected 'is', got '%s' or nothing", min)
    }

    empty := Empty[int]()
    _, ok = Max(empty)

    if ok {
        t.Errorf("expected empty")
    }

    _, ok = Min(empty)

    if ok {
        t.Errorf("expected empty")
    }

    emptyStr := Empty[string]()
    _, ok = MaxBy(emptyStr, func(a, b string) bool {
        return len(a) < len(b)
    })

    if ok {
        t.Errorf("expected empty")
    }

    _, ok = MinBy(emptyStr, func(a, b string) bool {
        return len(a) < len(b)
    })

    if ok {
        t.Errorf("expected empty")
    }

    iter = FromSlice(data)
    adv := iter.Advance(20)

    if adv != len(data) {
        t.Errorf("expected to advance %d, got %d", len(data), adv)
    }

    if Count(iter) > 0 {
        t.Error("expected empty")
    }

    iter.(*SliceIterator[string]).Reset()
    if iter.EstimatedRemaining() != len(data) {
        t.Errorf("expected %d remaining, got %d", len(data), iter.EstimatedRemaining())
    }

    asInts := Map(iter, func(s string) int { return len(s) })
    evensOnly := Filter(asInts, isEven)
    collected := evensOnly.Collect()
    reiter := FromSlice(collected)
    evenLenSum := Sum(reiter)

    if evenLenSum != 6 {
        t.Errorf("expected 6, got %d", evenLenSum)
    }
}
