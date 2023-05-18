# Priority discipline

## Purpose

Used to distributes data among handlers according to priority

Also may be used to equaling distribution of data with different processing times

## Principle of operation

Prioritization:

![Principle of operation, prioritization](./doc/operation-principle-321.svg)

Equaling:

![Principle of operation, equaling](./doc/operation-principle-222.svg)

## Comparison with unmanaged distribution

If different times are spent processing data of different priorities, then we will get different processing speeds in the case of using the priority discipline and without it:

Equaling by priority discipline:

![Equaling by priority discipline](./doc/different-processing-time-equaling.svg)

Unmanaged distribution:
![Unmanaged distribution](./doc/different-processing-time-unmanagement.svg)
