# CQOS

[![Go Reference](https://pkg.go.dev/badge/github.com/akramarenkov/cqos/v2.svg)](https://pkg.go.dev/github.com/akramarenkov/cqos/v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/akramarenkov/cqos/v2)](https://goreportcard.com/report/github.com/akramarenkov/cqos/v2)
[![codecov](https://codecov.io/gh/akramarenkov/cqos/branch/master/graph/badge.svg?token=2E4F42B30C)](https://codecov.io/gh/akramarenkov/cqos)

## Purpose

Library that allows you to control passage of data between Go channels

## Implemented disciplines

* **priority** - distributes data among handlers according to priority. See [README](./priority/README.md)

* **join** - accumulates elements from an input channel into a slice and write that slice to an output channel when the maximum slice size or timeout for its accumulation is reached. See [README](./join/README.md)

* **limit** - limits the speed of passing data elements from the input channel to the output channel. See [README](./limit/README.md)
