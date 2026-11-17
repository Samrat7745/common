# Java to Go Migration Report

<!-- markdownlint-disable MD013 -->

This report describes the migration of the Carbyne Stack Common library from
Java (version `0.2-SNAPSHOT`) to Go on the branch `go-migration`: what was
changed, how the API maps, and where the behaviour differs.

## Summary

- The library is now a Go module, `github.com/Samrat7745/common`, with two
  packages: `common` in the repository root replaces the Maven module
  `common-types`, and `commontest` replaces `common-testing`. Both use the
  standard library only.
- All Java sources, the Maven build and the Java-specific third-party license
  files were removed.
- The Java test suite was run before removal and passed (161 tests). Every Java
  test has a Go counterpart or is listed under [Tests](#tests) with the reason
  why it has none. All Go tests pass with 100% statement coverage.
- `Result` keeps all its operations and their semantics. The differences are
  listed in [Behaviour changes](#behaviour-changes); the ones most likely to
  affect users are that errors are returned instead of thrown (and panics are
  not turned into failures), that operations changing a type parameter are
  functions instead of methods, and that the exception-to-reason mapping is
  keyed by `reflect.Type`.

## What changed in the repository

### Removed

| Path                                                                 | Reason                                                                                                             |
| -------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------ |
| `common-types/src/main/java/**` (15 types)                           | Replaced by the package `common`                                                                                   |
| `common-types/src/main/java17/**` (3 types, `module-info.java`)      | Java 17 variant of `Result`, `Success`, `Failure`; see [Java 8 and Java 17 variants](#java-8-and-java-17-variants) |
| `common-types/src/test/java/**` (10 test classes)                    | Replaced by Go tests                                                                                               |
| `common-testing/src/main/java/**` (5 types)                          | Replaced by the package `commontest`                                                                               |
| `common-testing/src/test/java/**` (6 test classes)                   | Replaced by Go tests                                                                                               |
| `pom.xml`, `common-types/pom.xml`, `common-testing/pom.xml`          | Replaced by `go.mod`                                                                                               |
| `.github/templates/sbom.ftl`                                         | Template of the Maven license plugin                                                                               |
| `common-testing/3RD-PARTY-LICENSES/<java dependency>/**` (6 folders) | The Go module has no third-party dependencies                                                                      |

### Added

| Path                                | Content                                                                                           |
| ----------------------------------- | ------------------------------------------------------------------------------------------------- |
| `go.mod`                            | Module definition, Go 1.26 minimum, no dependencies                                               |
| `doc.go`                            | Package documentation                                                                             |
| `result.go`                         | `Result`, its constructors, methods and functions                                                 |
| `failure_reason.go`                 | `Describable`, `FailureReason`, its optional interfaces and their accessor functions, `ToFailure` |
| `unexpected.go`                     | `Unexpected`                                                                                      |
| `error_reason.go`                   | `FailureError`, `ErrorReason`                                                                     |
| `commontest/result_assert.go`       | `AssertThat`, `ResultAssert`                                                                      |
| `commontest/blank_strings.go`       | `BlankStrings`, `NamedString`                                                                     |
| `commontest/nil_combinations.go`    | `NilCombinations`                                                                                 |
| `commontest/doc.go`                 | Package documentation                                                                             |
| `*_test.go`, `commontest/*_test.go` | Tests                                                                                             |
| `MIGRATION_REPORT.md`               | This report                                                                                       |

### Moved

`common-testing/3RD-PARTY-LICENSES/sbom.xml` and `module.name` moved to
`3RD-PARTY-LICENSES/` in the repository root, as the folder `common-testing` no
longer exists. `sbom.xml` is now an empty dependency list and `module.name`
contains `common`.

### Layout

The root package is named `common`, so that the import path is the module path
and the types read as `common.Result`. The testing utilities are a separate
package because they import `testing`, which must not end up in production
binaries; the name `commontest` follows the standard library convention
(`httptest`, `iotest`). The three Java packages of `common-types` (`common`,
`common.result`, `common.function`) became one Go package: the first two depend
on each other (`CsFailureReason.toFailure` returns a `Failure`), which Go does
not allow across packages, and the third has no Go counterpart.

### Maintenance files

| File                                         | Change                                                                                                                                                                                                                                                                                                                                                                                                       |
| -------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `.github/workflows/build-test-deploy.yaml`   | Same name, workflow name (`Artifacts`), triggers and jobs (`verification`, `build`, `test`). The JDK 17 setup and the Maven calls became a Go 1.26 / 1.27 matrix running `gofmt`, `go build`, `go vet` and `go test -race`; coverage of the Go 1.27 run is uploaded to Codecov. The `verification` job (pre-commit) is unchanged.                                                                            |
| `.github/workflows/publish-and-release.yaml` | Same name, trigger and jobs (`prepare`, `draft-release`, `publish`). The version is no longer read from the POM: it is `v<BASE_VERSION>-<run id>-<run number>-<short sha>` with `BASE_VERSION: "0.3.0"` set in the workflow. The Maven deploy to GitHub Packages was dropped because a Go module is consumed from the repository. Collecting license obligation resources no longer fails if there are none. |
| `.github/scripts/generate-changelog.sh`      | The Maven dependency snippets in the release notes became a `go get <module>@<version>` snippet.                                                                                                                                                                                                                                                                                                             |
| `.pre-commit-config.yaml`                    | The exclusion of the third-party license folder follows its new location. There was no Java formatter hook to replace; see [Open points](#open-points).                                                                                                                                                                                                                                                      |
| `.gitignore`                                 | Go entries added; existing entries kept.                                                                                                                                                                                                                                                                                                                                                                     |
| `README.md`                                  | Describes the Go library, with usage and development instructions. Badges, license and contributing sections kept.                                                                                                                                                                                                                                                                                           |
| `3RD-PARTY-LICENSES/sbom.xml`                | Now an empty dependency list (see [Moved](#moved)).                                                                                                                                                                                                                                                                                                                                                          |

Unchanged: `LICENSE`, `NOTICE.md`, `CODEOWNERS`, `CODE_OF_CONDUCT.md`,
`CONTRIBUTING.md`, `.codacy.yaml`, `.github/workflows/auto-labeler.yml`,
`.github/workflows/stale-issue-cleanup.yml`.

The repository has no commitlint configuration and no release-please
configuration, so there was nothing to adapt for either.

## API mapping

### Package `common` (was `common-types`)

| Java                                                                               | Go                                                                                                                                                     |
| ---------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `Result<S, F>` (interface), `Success<S, F>`, `Failure<S, F>` (classes)             | `Result[S, F]` (one struct type)                                                                                                                       |
| `new Success<>(value)`                                                             | `Success[S, F](value)`                                                                                                                                 |
| `new Failure<>(reason)`                                                            | `Failure[S, F](reason)`                                                                                                                                |
| `Success.value()`                                                                  | `r.Value()` returning `(S, bool)`                                                                                                                      |
| `Failure.reason()`                                                                 | `r.Reason()` returning `(F, bool)`                                                                                                                     |
| `Result.of(AnyThrowingSupplier<S>, F reason)`                                      | `Of(supplier func() (S, error), reason F)`                                                                                                             |
| `Result.of(ThrowingSupplier<S, E>)`                                                | `OfReason(supplier func() (S, error))` returning `Result[S, FailureReason]`                                                                            |
| `Result.of(AnyThrowingSupplier<S>, Map<Class<? extends Throwable>, F>, F missing)` | `OfMapped(supplier func() (S, error), reasons map[reflect.Type]F, missing F)`                                                                          |
| `isSuccess()`, `isFailure()`                                                       | `r.IsSuccess()`, `r.IsFailure()`                                                                                                                       |
| `map(function)`                                                                    | `Map(r, function)`                                                                                                                                     |
| `mapFailure(function)`                                                             | `MapFailure(r, function)`                                                                                                                              |
| `tryMap(function, reason)`                                                         | `TryMap(r, function func(S) (N, error), reason)`                                                                                                       |
| `peek(consumer)`                                                                   | `r.Peek(consumer)`                                                                                                                                     |
| `tryPeek(consumer, reason)`                                                        | `r.TryPeek(consumer func(S) error, reason)`                                                                                                            |
| `recover(function)`                                                                | `r.Recover(function)`                                                                                                                                  |
| `flatMap(function)`                                                                | `FlatMap(r, function)`                                                                                                                                 |
| `unsafeFlatten()`                                                                  | `UnsafeFlatten[V, R](r)`                                                                                                                               |
| `fold(successFunction, failureFunction)`                                           | `Fold(r, onSuccess, onFailure)`                                                                                                                        |
| `filter(predicate, reason)`                                                        | `r.Filter(predicate, reason)`                                                                                                                          |
| `or(supplier)`                                                                     | `r.Or(supplier)`                                                                                                                                       |
| `swap()`                                                                           | `r.Swap()`                                                                                                                                             |
| `toOptional()`                                                                     | `r.Value()` returning `(S, bool)`                                                                                                                      |
| `stream()`                                                                         | `r.All()` returning `iter.Seq[S]`                                                                                                                      |
| `toString()`                                                                       | `r.String()`                                                                                                                                           |
| `equals()`, `hashCode()`                                                           | `==` (if `S` and `F` are comparable)                                                                                                                   |
| `Describable`                                                                      | `Describable`                                                                                                                                          |
| `CsFailureReason`                                                                  | `FailureReason` (requires `Synopsis()` only)                                                                                                           |
| `CsFailureReason.description()` (default: synopsis)                                | `Description(reason)`; a reason overrides it by implementing `Describable`                                                                             |
| `CsFailureReason.stackTrace()` (default: empty)                                    | `StackTrace(reason)` returning `([]runtime.Frame, bool)`; a reason provides one by implementing `StackTracer`                                          |
| `CsFailureReason.reportIssue()` (default: false)                                   | `ReportIssue(reason)`; a reason overrides it by implementing `IssueReporter`                                                                           |
| `CsFailureReason.toFailure()`                                                      | `ToFailure[S](reason)`                                                                                                                                 |
| `CsFailureThrowable`                                                               | `FailureError` (interface: `error` and `FailureReason`); `ErrorReason` / `NewErrorReason(err)` provides the Java default methods for an existing error |
| `Unexpected`, `new Unexpected(throwable)`                                          | `*Unexpected`, `NewUnexpected(err)`                                                                                                                    |
| `function.ThrowingSupplier<T, E>`, `AnyThrowingSupplier<T>`                        | No equivalent; `func() (T, error)`                                                                                                                     |
| `function.ThrowingConsumer<T, E>`, `AnyThrowingConsumer<T>`                        | No equivalent; `func(T) error`                                                                                                                         |
| `function.ThrowingFunction<T, R, E>`, `AnyThrowingFunction<T, R>`                  | No equivalent; `func(T) (R, error)`                                                                                                                    |
| `@Generated`                                                                       | Dropped: marked code for JaCoCo to ignore                                                                                                              |
| `@Stub`                                                                            | Dropped: marker annotation without behaviour                                                                                                           |
| `module-info.java`                                                                 | Dropped: Go has no module descriptor besides `go.mod`                                                                                                  |

New in Go, without a Java counterpart:

- `Unexpected.Err()` returns the error an `Unexpected` was created for. The Java
  class kept its throwable private.
- `ErrorReason.Error()` and `ErrorReason.Unwrap()`, so that the adapter works
  with `errors.Is` and `errors.As`.
- `StackTracer` and `IssueReporter`, the optional interfaces behind `StackTrace`
  and `ReportIssue`.

### Package `commontest` (was `common-testing`)

| Java                                                                                       | Go                                                                                         |
| ------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------ |
| `ResultAssert.assertThat(result)`                                                          | `AssertThat(t, result)` with `t` a `testing.TB`                                            |
| `ResultAssert<S, F>`                                                                       | `ResultAssert[S, F]`                                                                       |
| `isSuccess()`, `isFailure()`, `hasValue(v)`, `hasReason(r)`                                | `IsSuccess()`, `IsFailure()`, `HasValue(v)`, `HasReason(r)`; chainable as in Java          |
| Assertions inherited from AssertJ's `AbstractAssert` (`isNotNull`, `isEqualTo`, `as`, ...) | Dropped: there is no AssertJ; a `Result` cannot be nil                                     |
| `@EmptyOrBlankStringSource`, `EmptyOrBlankStringProvider`                                  | `BlankStrings()` returning `[]NamedString{Name, Value}`, to be used in a loop with `t.Run` |
| `@NullableParamSource("FIELD")`, `NullableParamGenerator`                                  | `NilCombinations(baseline...)` returning `[][]any`                                         |

JUnit annotations and argument providers have no Go counterpart. The two sources
were ported as functions that return the same data, as Go tests iterate over
their inputs themselves.

### Differences in API shape

- **Type-changing operations are functions, not methods.** Go methods cannot
  have their own type parameters, so `Map`, `MapFailure`, `TryMap`, `FlatMap`,
  `Fold` and `UnsafeFlatten` are package-level functions that take the result as
  first argument. Operations that keep both types are methods. A chain that
  mixes both has to be split into statements.
- **One type instead of three.** `Result` is a struct value, not an interface
  with two implementations. Code that used `instanceof Success` uses
  `IsSuccess()` or the boolean of `Value()`. A `Result` cannot be nil. Its zero
  value is a failure with the zero value of `F` as reason.
- **Errors are returned, not thrown.** Where Java took a function that may
  throw, Go takes a function that returns an `error`.
- **Optional values are `(value, bool)` pairs** and streams are iterators
  (`iter.Seq`).
- **Default methods became functions.** A Go interface cannot provide default
  implementations. `FailureReason` requires only `Synopsis()`; the other
  properties are read with `Description`, `StackTrace` and `ReportIssue`, which
  apply the Java defaults unless the reason implements the corresponding
  optional interface.
- **Names.** The prefix `Cs` was dropped, as the package name already qualifies
  the types (`common.FailureReason`).

## Behaviour changes

### Java 8 and Java 17 variants

`common-types` was a multi-release JAR with two implementations of `Result`,
`Success` and `Failure`. Comparing `src/main/java` and `src/main/java17` showed
them to be the same code except for language features (sealed interface,
records, pattern matching for `instanceof`) and one visible difference:

|              | Java 8 variant                            | Java 17 variant                                                     | Go                                        |
| ------------ | ----------------------------------------- | ------------------------------------------------------------------- | ----------------------------------------- |
| `toString()` | `Success{value=12}`, `Failure{reason=21}` | `Success[value=12]`, `Failure[reason=21]` (generated by the record) | `Success[value=12]`, `Failure[reason=21]` |

Go uses the Java 17 format, which is also the one the Java tests of
`ResultAssert` expected.

### Preserved on purpose

- **Semantics of every operation**, including which function is called for a
  success and which for a failure, and that the function of the other case is
  not called.
- **Nil functions are rejected in both cases.** Java threw a
  `NullPointerException` for a `null` function even if the operation would not
  have called it (for example `map` on a `Failure`). Go panics with a message
  starting with `common:` in the same situations. Without the explicit check, a
  nil function would only have failed when called.
- **Resolution order of the mapped reasons.** `OfMapped` looks up the type of
  the cause first, then the type of the error, then a key the error's type is
  assignable to, and uses the `missing` reason last.
- **Text of `Unexpected`.** The synopsis is
  `An unknown exception has occurred.`, the description is the message of the
  error, and `ReportIssue` is true.
- **Defaults of a failure reason.** The description defaults to the synopsis,
  the stack trace to none, and `ReportIssue` to false.
- **Failure messages of `ResultAssert`.**
  `Expecting a success result but was: Failure[reason=21]`,
  `Expecting result success value to equal 21 but was: 12` and their failure
  counterparts are unchanged, including the single quotes around strings and the
  escaping of line breaks.
- **A failed assertion ends the test.** AssertJ threw an `AssertionError`; Go
  calls `Fatalf`.
- **The list of blank strings** is the same 26 entries with the same names.
- **Order of the nil combinations.** `NilCombinations` returns the combinations
  in the order `NullableParamGenerator` produced them, without the combination
  that replaces nothing.

### Changed

| Area                                                      | Java                                                                                                                                                           | Go                                                                                                                                                                                     | Impact                                                                                                         |
| --------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------- |
| Failing functions                                         | `Result.of`, `tryMap` and `tryPeek` caught every `Throwable`, including unchecked exceptions such as `NullPointerException`                                    | Only a returned `error` makes a failure; a panic is not recovered and propagates                                                                                                       | A function that panics crashes the caller instead of yielding a failure.                                       |
| `null` reason                                             | `Result.of(supplier, null)` threw a `NullPointerException`                                                                                                     | A nil reason is accepted                                                                                                                                                               | Go cannot test a value of a type parameter for nil without reflection.                                         |
| `null` mapping                                            | `Result.of(supplier, null, missing)` threw a `NullPointerException`                                                                                            | A nil map is an empty map                                                                                                                                                              | Every error maps to `missing`.                                                                                 |
| `null` results                                            | `flatMap` and `or` threw a `NullPointerException` if the function returned `null`                                                                              | Not applicable: a `Result` is a value and cannot be nil                                                                                                                                | None.                                                                                                          |
| Reason from the error (`of(ThrowingSupplier)`)            | The thrown exception itself had to implement `CsFailureReason`; the failure type was the exception type `E`; any other exception caused a `ClassCastException` | `OfReason` uses the first error in the chain (`errors.As`) that implements `FailureReason`; the failure type is the interface `FailureReason`; any other error becomes an `Unexpected` | Wrapped reasons are found. Callers that need the concrete type use a type assertion or `MapFailure`.           |
| Keys of the mapped reasons                                | `Class<? extends Throwable>`                                                                                                                                   | `reflect.Type`, obtained with `reflect.TypeFor[T]()`                                                                                                                                   | Call sites change.                                                                                             |
| Subclass matching                                         | A key matched every subclass of its exception class (`isAssignableFrom`)                                                                                       | Go has no subclasses: step 3 matches keys that are interface types implemented by the error, for example `reflect.TypeFor[error]()` for every error                                    | A hierarchy of exception classes has to be expressed by an interface.                                          |
| Cause of an error                                         | `Throwable.getCause()`                                                                                                                                         | `errors.Unwrap(err)`                                                                                                                                                                   | An error that wraps several errors (`errors.Join`, several `%w`) has no cause for step 1.                      |
| Sentinel errors                                           | Not applicable                                                                                                                                                 | Errors made by `errors.New` or `fmt.Errorf` all have the same type and cannot be told apart by `OfMapped`                                                                              | See [Open points](#open-points).                                                                               |
| Several assignable keys                                   | The first match in the iteration order of the map, which is unspecified for a `HashMap`                                                                        | The match whose package path and type name sort first                                                                                                                                  | Deterministic.                                                                                                 |
| `toOptional()` of a `null` success value                  | Empty `Optional`                                                                                                                                               | `Value()` returns `nil, true`                                                                                                                                                          | A nil success value is reported as present.                                                                    |
| `unsafeFlatten()` type parameters                         | Inferred from the assignment; a mismatch surfaced later as a `ClassCastException` where the value was used                                                     | Given explicitly; a mismatch panics in `UnsafeFlatten`                                                                                                                                 | Fails earlier.                                                                                                 |
| Stack traces                                              | `StackTraceElement[]` of the place the exception was created                                                                                                   | `[]runtime.Frame` of the call to `NewUnexpected` or `NewErrorReason`, as Go errors carry no stack                                                                                      | The trace shows where the error was converted, not where it originated.                                        |
| `CsFailureThrowable` on a class that is not a `Throwable` | Synopsis `Throwable <toString> has been caught.`, empty stack trace                                                                                            | Not applicable: `FailureError` includes `error`                                                                                                                                        | The fallback text is gone.                                                                                     |
| Equality                                                  | `equals` and `hashCode` for every `Success` and `Failure`                                                                                                      | `==`, which compiles only if `S` and `F` are comparable                                                                                                                                | Results of slices, maps or functions cannot be compared with `==` or used as map keys.                         |
| Equality in assertions                                    | `Objects.equals`                                                                                                                                               | `reflect.DeepEqual`                                                                                                                                                                    | Pointers to equal values are equal; slices and maps are compared by content.                                   |
| Line breaks in assertion messages                         | `System.lineSeparator()` was escaped                                                                                                                           | `\n` is escaped                                                                                                                                                                        | Same on Linux and macOS.                                                                                       |
| Order of blank strings                                    | Iteration order of a `HashMap` (unspecified)                                                                                                                   | Fixed, in the order of the Java source                                                                                                                                                 | Deterministic subtest order.                                                                                   |
| White space definition                                    | All 26 blank strings are blank: every character is white space for `Character.isWhitespace`                                                                    | `FILE`, `GROUP`, `RECORD` and `UNIT SEPARATOR` (U+001C to U+001F) are not white space for `unicode.IsSpace` and `strings.TrimSpace`                                                    | Code under test that trims with `strings.TrimSpace` does not treat these four as blank. A test documents this. |
| Baseline of the nil combinations                          | Read by reflection from a static `Arguments` field named in the annotation; a missing field raised `Failed to load test params from static variable.`          | Passed as arguments                                                                                                                                                                    | No reflection, no such error. Values come back as `any` and are converted with a type assertion.               |
| Distribution                                              | Maven artifacts `common-types` and `common-testing` on GitHub Packages                                                                                         | One Go module fetched from the repository                                                                                                                                              | See [Open points](#open-points) on versions.                                                                   |

## Tests

### Java baseline

Before the Java sources were removed, `mvn verify` was run on OpenJDK 17.0.21
with Maven 3.9.9. The numbers are taken from the Surefire output in the build
log; the build ended with `BUILD SUCCESS`.

| Module           | Tests run | Failures | Errors | Skipped |
| ---------------- | --------- | -------- | ------ | ------- |
| `common-types`   | 37        | 0        | 0      | 0       |
| `common-testing` | 124       | 0        | 0      | 0       |
| Total            | 161       | 0        | 0      | 0       |

The tests of `common-testing` ran on JDK 17 against the multi-release JAR and so
against the Java 17 variant, which shows in the record-style `toString` text
they expected. The tests of `common-types` ran against the module's class
folder, where the Java 8 variant is at the top level and the Java 17 variant
under `META-INF/versions/17`; they most likely exercised the Java 8 variant, but
none of them distinguishes the two.

### Go

| Java test class (module)                                  | Tests | Go counterpart                                                                            |
| --------------------------------------------------------- | ----- | ----------------------------------------------------------------------------------------- |
| `UnexpectedTest` (types)                                  | 5     | `unexpected_test.go`                                                                      |
| `CsFailureReasonTest` (types)                             | 4     | `failure_reason_test.go`                                                                  |
| `CsFailureThrowableTest` (types)                          | 4     | `error_reason_test.go`                                                                    |
| `ResultTest` (types)                                      | 9     | `result_test.go` (`TestResult_Chain*`, `TestUnsafeFlatten`, `TestResult_Swap`, `TestOf*`) |
| `ThrowingSupplierTest`, `AnyThrowingSupplierTest` (types) | 5     | `function_test.go` (`TestSupplierReturningError`)                                         |
| `ThrowingConsumerTest`, `AnyThrowingConsumerTest` (types) | 5     | `function_test.go` (`TestConsumerReturningError`)                                         |
| `ThrowingFunctionTest`, `AnyThrowingFunctionTest` (types) | 5     | `function_test.go` (`TestFunctionReturningError`)                                         |
| `SuccessTest` (testing)                                   | 35    | `success_test.go`                                                                         |
| `FailureTest` (testing)                                   | 35    | `failure_test.go`                                                                         |
| `ResultTest` (testing)                                    | 17    | `result_test.go` (`TestOf*`, `TestOfReason*`, `TestOfMapped*`)                            |
| `CsFailureReasonTest` (testing)                           | 1     | `failure_reason_test.go` (`TestFailureReason_ToFailure`)                                  |
| `BlankStringTest` (testing)                               | 26    | `commontest/blank_strings_test.go`                                                        |
| `ResultAssertTest` (testing)                              | 10    | `commontest/result_assert_test.go`                                                        |

Java tests without a direct counterpart, because the behaviour does not exist in
Go:

- `CsFailureThrowableTest`: the synopsis of an implementation that is not a
  `Throwable`.
- `SuccessTest` and `FailureTest`: the `NullPointerException` for a function
  that returns a `null` result in `flatMap` and `or`.
- `ResultTest` (testing): the `NullPointerException` for a `null` reason and a
  `null` mapping. The Go tests assert the new behaviour instead (nil accepted).
- The tests of the `function` package checked that a lambda can be assigned to
  the functional interfaces and called. The Go tests check the three function
  shapes where the `Result` API takes them.

The parameterized tests using `@NullableParamSource` were ported with
`NilCombinations`. `ResultAssert` is tested with a `testing.TB` that records the
calls to `Fatalf`.

Tests were added for what the Java suite did not cover: `OfReason` with wrapped
and with other errors, the precedence and the tie-break of `OfMapped`, `String`,
`==`, the zero value, the type mismatch of `UnsafeFlatten`, `All`, the optional
interfaces of a failure reason, `ErrorReason`, `NilCombinations`, and the
difference between the blank strings and `strings.TrimSpace`.

Results:

| Command                               | Toolchain | Result                                                                                                    |
| ------------------------------------- | --------- | --------------------------------------------------------------------------------------------------------- |
| `go test -race -count=1 -cover ./...` | Go 1.27.4 | passed; `common` 100.0%, `commontest` 100.0% of statements                                                |
| `go test -count=1 -json ./...`        | Go 1.27.4 | 139 test functions passed (`common` 120, `commontest` 19), 187 including subtests; none failed or skipped |
| `go test -count=1 ./...`              | Go 1.26.8 | passed                                                                                                    |
| `go build ./...`, `go vet ./...`      | Go 1.27.4 | clean                                                                                                     |
| `go vet ./...`                        | Go 1.26.8 | clean                                                                                                     |
| `gofmt -l .`                          | Go 1.27.4 | no output                                                                                                 |

Also checked locally:

- The code examples of `README.md` and of the package documentation compile and
  run in a scratch module.
- `mdformat` 0.7.22 with `mdformat-gfm` (`--check --wrap 80`) and
  `markdownlint-cli` 0.44.0, the versions pinned in `.pre-commit-config.yaml`,
  pass on `README.md` and on this report. They were called directly, not through
  pre-commit.
- The shell step collecting license obligation resources exits successfully with
  the empty `3RD-PARTY-LICENSES` folder, and `generate-changelog.sh` produces
  the release notes with the `go get` snippet.

Not verified: the GitHub workflows were not executed, as that needs a push.

## Open points

- **Module path.** The module path is `github.com/Samrat7745/common` as
  requested. The license headers, the README badges and `CODEOWNERS` still refer
  to `carbynestack/common`. Moving the module to that organisation changes the
  import path in `go.mod`, in the imports of `commontest` and the tests, and in
  the README.
- **Versions.** The POM carried the version (`0.2-SNAPSHOT`); `go.mod` has none.
  The release workflow now takes it from `BASE_VERSION`, which was set to
  `0.3.0` because of the breaking change. This value is a proposal. The tags
  have the form `v0.3.0-<run id>-<run number>-<short sha>`, which Go accepts as
  pre-release versions; there is still no process for a final `vX.Y.Z` tag, as
  there was none for Maven releases.
- **Workflow runners and actions.** The workflows keep `ubuntu-20.04` and the
  action versions of the original files where a step was not rewritten. They
  were not run; whether that runner image and those action versions are still
  available on GitHub was not checked.
- **Go formatting in pre-commit.** The pre-commit configuration had no Java
  formatter hook, so no Go formatter hook was added: the `verification` job that
  runs pre-commit sets up Python only, and a Go hook would need a Go toolchain
  there. Formatting is checked by `gofmt` in the `build` job instead.
- **Panics.** Java turned every exception thrown by a supplier, function or
  consumer into a failure. Go only does so for returned errors. If panics of
  callbacks should become failures as well, `Of`, `TryMap` and `TryPeek` would
  need a `recover`.
- **Sentinel errors in `OfMapped`.** The mapping by type is the faithful port,
  but idiomatic Go often identifies errors by value (`errors.Is`). A mapping
  that supports sentinel errors would be an addition to the API.
- **Blank strings.** `BlankStrings` keeps the four separator characters that are
  white space for Java but not for Go. Whether Go services should treat them as
  blank is a decision for the services; the list can be reduced if not.
- **`NOTICE.md`** was left as is; add the contributors of the migration as
  appropriate.
- **Codecov** needs no configuration change, but the first Go coverage upload
  will not be comparable with the Java history.
