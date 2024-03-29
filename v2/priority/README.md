# Priority discipline

## Purpose

Used to distributes data among handlers according to priority

Also may be used to equaling distribution of data with different processing times

## Principle of operation

* Prioritization:

  ![Principle of operation, prioritization](./doc/operation-principle-321.svg)

* Equaling:

  ![Principle of operation, equaling](./doc/operation-principle-222.svg)

## Comparison with unmanaged distribution

If different times are spent processing data of different priorities, then we will get different processing speeds in the case of using the priority discipline and without it.

For example, suppose that data from channel of priority 3 is processed in time **T**, data from channel of priority 2 is processed in time 5\***T**, and data from channel of priority 1 is processed in time 10\***T**, then we will get the following results:

* equaling by priority discipline:

  ![Equaling by priority discipline](./doc/different-processing-time-equaling.png)

* unmanaged distribution:

  ![Unmanaged distribution](./doc/different-processing-time-unmanagement.png)

It can be seen that with unmanaged distribution, the processing speed of data with priority 3 is limited by the slowest processed data (with priority 1 and 2), but at with equaling by priority discipline the processing speed of data with priority 3 is no limited by others priorities

## Usage

Example:

```go
package main

import (
    "fmt"
    "strconv"
    "sync"

    "github.com/akramarenkov/cqos/v2/priority"
    "github.com/akramarenkov/cqos/v2/priority/divider"
)

func main() {
    handlersQuantity := 100
    // Preferably input channels should be buffered
    inputCapacity := 10
    itemsQuantity := 100

    inputs := map[uint]chan string{
        3: make(chan string, inputCapacity),
        2: make(chan string, inputCapacity),
        1: make(chan string, inputCapacity),
    }

    // Map key is a value of priority
    inputsOpts := map[uint]<-chan string{
        3: inputs[3],
        2: inputs[2],
        1: inputs[1],
    }

    // Used only in this example for detect that all written data are processed
    measures := make(chan string)
    defer close(measures)

    // For equaling use divider.Fair divider, for prioritization use
    // divider.Rate divider or custom divider
    opts := priority.Opts[string]{
        Divider:          divider.Rate,
        HandlersQuantity: uint(handlersQuantity),
        Inputs:           inputsOpts,
    }

    discipline, err := priority.New(opts)
    if err != nil {
        panic(err)
    }

    wg := &sync.WaitGroup{}
    defer wg.Wait()

    // Run writers, that write data to input channels
    for priority, input := range inputs {
        wg.Add(1)

        go func(precedency uint, channel chan string) {
            defer wg.Done()
            defer close(channel)

            base := strconv.Itoa(int(precedency))

            for id := 0; id < itemsQuantity; id++ {
                item := base + ":" + strconv.Itoa(id)

                channel <- item
            }
        }(priority, input)
    }

    // Run handlers, that process data
    for handler := 0; handler < handlersQuantity; handler++ {
        wg.Add(1)

        go func() {
            defer wg.Done()

            for prioritized := range discipline.Output() {
                // Data processing
                measures <- prioritized.Item

                // Handler must indicate that current data has been processed and
                // handler is ready to receive new data
                discipline.Release(prioritized.Priority)
            }
        }()
    }

    received := 0

    // Wait for process all written data
    for range measures {
        received++

        if received == itemsQuantity*len(inputs) {
            break
        }
    }

    fmt.Println("Processed items quantity:", received)
    // Output: Processed items quantity: 300
}
```
