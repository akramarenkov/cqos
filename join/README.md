# Join discipline

## Purpose

Accumulates elements from the input channel into a slice and writes it to the output channel when the size or timeout is reached

Works in two modes:

1. Making a copy of the slice before writing it to the output channel

2. Writes to the output channel of the accumulated slice without copying, in this case it is necessary to inform the discipline that the slice is no longer used by writing to the Released channel

## Usage

Example:

```go
package main

import (
    "fmt"
    "time"

    "github.com/akramarenkov/cqos/join"
)

func main() {
    quantity := 27

    input := make(chan int)

    opts := join.Opts[int]{
        Input:    input,
        JoinSize: 5,
        Timeout:  10 * time.Second,
    }

    discipline, err := join.New(opts)
    if err != nil {
        panic(err)
    }

    go func() {
        defer close(input)

        for stage := 1; stage <= quantity; stage++ {
            input <- stage
        }
    }()

    outSequence := make([]int, 0, quantity)

    for slice := range discipline.Output() {
        outSequence = append(outSequence, slice...)
    }

    fmt.Println(outSequence)
    // Output:[1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27]
}
```
