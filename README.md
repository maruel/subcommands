subcommands golang library
==========================

This package permits a Go application to implement subcommands support
similar to what is supported by the 'go' tool.

The library is designed so that the test cases can run concurrently.
Using global flags variables is discouraged to keep your program testable
conccurently.

See module subcommands/subcommandstest for tools to help testing an application
using subcommands.
