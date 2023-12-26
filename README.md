# CQOS

[![Go Reference](https://pkg.go.dev/badge/github.com/akramarenkov/cqos.svg)](https://pkg.go.dev/github.com/akramarenkov/cqos)
[![Go Report Card](https://goreportcard.com/badge/github.com/akramarenkov/cqos)](https://goreportcard.com/report/github.com/akramarenkov/cqos)
[![codecov](https://codecov.io/gh/akramarenkov/cqos/branch/master/graph/badge.svg?token=2E4F42B30C)](https://codecov.io/gh/akramarenkov/cqos)

## Purpose

Library that allows you to control passage of data between Go channels

## Implemented disciplines

* **priority** - distributes data among handlers according to priority. See [README](./priority/README.md)

* **join** - accumulates elements from the input channel into a slice and writes it to the output channel when the size or timeout is reached. See [README](./join/README.md)

[Improved version V2 available](./v2)
