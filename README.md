# subcommands golang library

This package permits a Go application to implement subcommands support
similar to what is supported by the 'go' tool.

The library is designed so that the test cases can run concurrently.
Using global flags variables is discouraged to keep your program testable
concurrently.

The intended command is found via heuristic search;

  - exact match
  - unique prefix, e.g. `lo` will run `longcommand` as long as there's no
    command with the same prefix.
  - case insensitivity; for those weird enough to use Upper Cased Commands.
  - [levenshtein distance](http://en.wikipedia.org/wiki/Levenshtein_distance);
    where `longcmmand` or `longcmomand` will properly trigger `longcommand`.

[![PkgGoDev](https://pkg.go.dev/badge/github.com/maruel/subcommands)](https://pkg.go.dev/github.com/maruel/subcommands)
[![Coverage Status](https://codecov.io/gh/maruel/subcommands/graph/badge.svg)](https://codecov.io/gh/maruel/subcommands)


## Examples

  - See [sample-simple](sample-simple) for a barebone sample skeleton usable
    as-is.
  - See [sample-complex](sample-complex) for a complex sample using advanced
    features.
  - See module
    [subcommands/subcommandstest](https://pkg.go.dev/github.com/maruel/subcommands/subcommandtest)
    for tools to help *testing* an application using subcommands. One of the
    main benefit is t.Parallel() just works, because subcommands help wrapping
    global variables.
