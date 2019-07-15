
Native Go bindings for Clang's C API (libclang) for versions 5.0, 6.0, 7.0 and 8.0.

Refer to https://llvm.org for the latest clang versions and documentation.

## Alpha

This should still be considered alpha as the API deserves further cleanup.
The *error* go type should be used more extensively. Enum types typically have
an *Invalid* value which should be used to trigger an error return.

## Forked

Forked from [https://github.com/go-clang/v3.9](https://github.com/go-clang/v3.9) some years ago.

This first incorporated the clang-c headers from 5.0.0 and came with
some necessary changes as a result.
Some bug fixes and some API changes are also included.

The git log will show what has changed. Or the quick way for users of previous versions to incorporate this
would be to change their import lines and see what no longer compiles. (Also using git log on the `*_test.go`
and cmd examples shows what API changes were made to keep the tests passing.)

This fork was not meant to be completely backward compatible with the v3.9 repository.
As with that repository, this comes with a liberal license and users are free to take any parts they like.

## Usage

Usage hasn't changed drastically. Some return types have changed so that some Dispose calls could be done away with.
Some return error values have been changed to be more go idiomatic. In at least one case, a routine's return values
were swapped so the error result is second.

As before, an example on how to use the AST visitor of the Clang API can be found in [/cmd/go-clang-dump/main.go](/cmd/go-clang-dump/main.go)

## Generated Bindings

The v3.9 bindings were used as a base.
The gen tool from github.com/go-clang was also used on the 3.9 headers to get a sense for what had to be manually changed when the v3.9 repository was created.
The gen tool was used on the new headers which saved a considerable amount of time.
A diff between 5.0 and 3.9 headers was done to see how the libclang API had changed and what might need to be accounted for manually.

Support for clang 6.0, 7.0, and 8.0 was added by hand by looking at the differences to the clang-c headers from one version to the next.

Unit tests were not added to keep up with the changes so the new code is largely untested.

## Build and run self tests.

Once you have downloaded the repository:

```bash
  source env.sh
  cd clang
  go install
  go test

  cd ../cmd/go-clang-dump
  go build
  go test

  cd ../go-clang-compdb
  go build
  go test

  cd ../go-clang-includes
  go build
  go test

  cd ../go-clang-globals
  go build
  go test
```

## Older platforms tested.

| Platform | clang+llvm |
| --- | --- |
| 10.12 Darwin | clang+llvm 5.0.0 |
| 10.13 Darwin | clang+llvm 5.0.1 |
| 4.9 Debian | clang+llvm 5.0.1 built for Debian8 |
| 11 FreeBSD | clang+llvm 5.0.1 built for FreeBSD10 |

All builds done with go1.10. At the time, this was also found to build and test successfully with the new vgo command.

## Current platform tested.

| Platform | clang+llvm |
| --- | --- |
| 10.14.5 Darwin | clang+llvm 8.0.0 |

Built with with go1.12.7.

# License

This repository comes with the 3-Clause BSD License.
