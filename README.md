# bytez

Package *bytez* provides functionality for working with large **byte** si**z**es and capacities in a
human-friendly way in the belief that code and configuration files are more readable and less
error-prone when they are specified with, for example, "4MiB" rather than "4194304".

The 'z' in the name is not an attempt to be cute but rather to create a package name that is
short yet unique.

## Installation

Download and install module with

    go get github.com/codexnull/bytez

## Usage

To use, simply

    import github.com/codexnull/bytez

Then the functions `bytez.AsStr()` and `bytez.AsInt()` can be used to convert to and from size
specifications and number of bytes.

For example, if `bufSize` is 16384, `bytez.AsStr(bufSize)` will return `"16KiB"`. And if `cacheSize`
is "8MiB", `bytez.AsInt(cacheSize)` will return `8388608`.

The package also provides a `Size` type that supports marshaling to and from text using the
functions above when marshaling to and from JSON, YAML, etc.

See the [godoc](https://godoc.org/github.com/codexnull/bytez) for more details.
