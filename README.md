# Carbyne Stack Common

[![codecov](https://codecov.io/gh/carbynestack/common/branch/master/graph/badge.svg?token=XY3TLGJYEC)](https://codecov.io/gh/carbynestack/common)
[![Known Vulnerabilities](https://snyk.io/test/github/carbynestack/common/badge.svg)](https://snyk.io/test/github/carbynestack/common)
[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit&logoColor=white)](https://github.com/pre-commit/pre-commit)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](CODE_OF_CONDUCT.md)

This library provides a set of shared types and testing utilities for all
Carbyne Stack services, clients, and other components written in Go. The content
is split into two packages. **common** for public API signatures, as well as
other shared code and **commontest** for utilities which are needed to implement
tests in Carbyne Stack in a desirable manner, with result assertions and
generators of test inputs.

> **NOTE**: This library has been migrated from Java to Go. See the
> [migration report](MIGRATION_REPORT.md) for the API mapping and the changes in
> behaviour.

## Usage

The module requires Go 1.26 or later and has no dependencies besides the
standard library.

```shell
go get github.com/Samrat7745/common
```

A `Result` holds either a success value or a failure reason:

```go
import "github.com/Samrat7745/common"

// A success if the function returns no error, a failure with the reason if it
// does.
port := common.Of(func() (int, error) {
    return strconv.Atoi(input)
}, "not a number")

// Operations that keep the types of the result are methods.
port = port.
    Filter(func(p int) bool { return p > 0 && p < 65536 }, "out of range").
    Recover(func(reason string) int { return 8080 })

// Operations that change the success or the failure type are functions.
address := common.Map(port, func(p int) string {
    return fmt.Sprintf("localhost:%d", p)
})

if value, ok := address.Value(); ok {
    fmt.Println(value)
}
```

Failure reasons can be of any type. Reasons that describe themselves implement
`common.FailureReason`; `common.Unexpected` and `common.ErrorReason` adapt an
error to it.

The package `commontest` supports tests of code that returns results:

```go
import "github.com/Samrat7745/common/commontest"

func TestParsePort(t *testing.T) {
    commontest.AssertThat(t, ParsePort("8080")).HasValue(8080)

    for _, blank := range commontest.BlankStrings() {
        t.Run(blank.Name, func(t *testing.T) {
            commontest.AssertThat(t, ParsePort(blank.Value)).IsFailure()
        })
    }
}
```

## Development

```shell
go build ./...
go vet ./...
go test -race ./...
```

## License

Carbyne Stack *Common* is open-sourced under the Apache License 2.0. See the
[LICENSE](LICENSE) file for details.

### 3rd Party Licenses

For information on how license obligations for 3rd party OSS dependencies are
fulfilled see the [README](https://github.com/carbynestack/carbynestack) file of
the Carbyne Stack repository.

## Contributing

Please see the Carbyne Stack
[Contributor's Guide](https://github.com/carbynestack/carbynestack/blob/master/CONTRIBUTING.md)
.
