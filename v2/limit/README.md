# Limit discipline

## Purpose

Limits the speed of passing data elements from the input channel to the output channel

The speed limit is set by the **Rate** structure, in which the **Quantity** field specifies the number of data elements that must pass over the time interval specified in the **Interval** field

As we know, the speed of 1000 data elements per second is, in fact, the same speed as 1 data element per millisecond specified in different units of measurement

However, the units of measurement affect the distribution of data elements written to the output channel over time and the performance of the discipline

If the speed is specified as 1000 data elements per second, first 1000 data elements will be written to the output channel, and then a pause will be made equal to 1 second minus the time spent writing 1000 data elements

If the speed is specified in the form of 1 data element per millisecond, first 1 data element will be written to the output channel, and then a pause will be made equal to 1 millisecond minus the time spent on writing 1 data element

However, the performance of the discipline if the speed is specified in the form of 1 data element per millisecond will be lower

Thus, when choosing units of measurement, you can balance between the uniform distribution of data elements over time and performance (the maximum achievable speed)

Based on measurements, specifying a time interval of less than 10 milliseconds greatly reduces performance

## Usage

Example:

```go
package main

import (
    "fmt"
    "time"

    "github.com/akramarenkov/cqos/v2/limit"
)

func main() {
    quantity := 10

    input := make(chan int)

    opts := limit.Opts[int]{
        Input: input,
        Limit: limit.Rate{
            Interval: time.Second,
            Quantity: 1,
        },
    }

    discipline, err := limit.New(opts)
    if err != nil {
        panic(err)
    }

    outSequence := make([]int, 0, quantity)

    startedAt := time.Now()

    go func() {
        defer close(input)

        for stage := 1; stage <= quantity; stage++ {
            input <- stage
        }
    }()

    for item := range discipline.Output() {
        outSequence = append(outSequence, item)
    }

    duration := time.Since(startedAt)
    expected := (time.Duration(quantity) * opts.Limit.Interval) / time.Duration(opts.Limit.Quantity)
    deviation := 0.01

    fmt.Println(duration <= time.Duration(float64(expected)*(1.0+deviation)))
    fmt.Println(duration >= time.Duration(float64(expected)*(1.0-deviation)))
    fmt.Println(outSequence)
    // Output:
    // true
    // true
    // [1 2 3 4 5 6 7 8 9 10]
}
```
