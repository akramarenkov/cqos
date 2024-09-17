# Join discipline

## Purpose

Accumulates elements from an input channel into a slice and write that slice to an output channel when the maximum slice size or timeout for its accumulation is reached

Works in two modes:

1. Making a copy of the slice before writing it to the output channel

2. Writes to the output channel of the accumulated slice without copying, in this case it is necessary to inform the discipline that the slice is no longer used by call the Release() method

## Usage

Example:

```go
package main

import (
    "fmt"
    "time"

    "github.com/akramarenkov/cqos/v2/join"
)

func main() {
    data := []int{
        1, 2, 3, 4, 5, 6, 7, 8,
        9, 10, 11, 12, 13, 14, 15, 16,
        17, 18, 19, 20, 21, 22, 23, 24,
        25, 26, 27,
    }

    // Preferably input channel should be buffered for performance reasons.
    // Optimal capacity is in the range of one to three JoinSize
    input := make(chan int, 10)

    opts := join.Opts[int]{
        Input:    input,
        JoinSize: 10,
        Timeout:  time.Second,
    }

    discipline, err := join.New(opts)
    if err != nil {
        panic(err)
    }

    go func() {
        defer close(input)

        for _, item := range data {
            input <- item
        }
    }()

    for join := range discipline.Output() {
        fmt.Println(join)
    }

    // Output:
    // [1 2 3 4 5 6 7 8 9 10]
    // [11 12 13 14 15 16 17 18 19 20]
    // [21 22 23 24 25 26 27]
}
```
