#!/usr/bin/env bash

set -euf -o pipefail

set -x

CGO_LDFLAGS=
case $(uname) in
    Darwin)
        DYLD_LIBRARY_PATH=
        ;;
    *)
        LD_LIBRARY_PATH=
        ;;
esac

. env.sh

cd clang
go install ./...
go test ./...

cd ../cmd/go-clang-dump
go build
go test

cd ../go-clang-compdb
go build
go test

cd ../go-clang-includes
go build
go test -cflags="$CGO_CPPFLAGS"

cd ../go-clang-globals
go build
go test -cflags="$CGO_CPPFLAGS"
