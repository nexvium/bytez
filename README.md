[![Go Report Card](https://goreportcard.com/badge/github.com/codexnull/bytez)](https://goreportcard.com/report/github.com/codexnull/bytez)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/codexnull/bytez)
[![LICENSE](https://img.shields.io/github/license/codexnull/bytez.svg?style=flat-square)](https://github.com//codexnull/bytez/blob/master/LICENSE)

# bytez

Package *bytez* provides functionality for working with large **byte** si**z**es and capacities in a
human-friendly way in the belief that code and configuration files are more readable and less
error-prone when they are specified with, for example, "4MiB" rather than "4194304".

The 'z' in the name is not an attempt to be cute but rather to create a package name that is
short yet unique.

## Installation

Download and install module with:

    go get github.com/codexnull/bytez

## Usage

To use, simply

    import github.com/codexnull/bytez

in your code.

Then, sizes can be specified using constants. For example, the code

```go
bufSize := 16*bytez.Kibibyte
buffer := make([]byte, bufSize)
fmt.Printf("Buffer Size: %v (%v)\n", len(buffer), bytez.AsStr(bufSize))
```

outputs

```text
Buffer Size: 16384 (16KiB)
```

The functions `bytez.AsStr()` and `bytez.AsInt()` can be used to convert to and from size
specifications and number of bytes. For example, if `bufSize` is 16384, `bytez.AsStr(bufSize)` will return `"16KiB"`. And if `cacheSize`
is "8MiB", `bytez.AsInt(cacheSize)` will return `8388608`.

The package also provides a `Size` type that supports marshaling to and from text using the
functions above when marshaling to and from JSON, YAML, etc. So the code

```go
type Conf struct {
    ...
    CacheSize bytez.Size `json:"cache_size"`
    ...
}
conf := Conf{CacheSize: bytez.Size(2 * bytez.Gibibyte)}
```

Will result in

```text
{
  ...
  "cache_size":"2GiB",
  ...
}
```

when `conf` is marshaled to JSON. Unmarshaling will result in the value `2 * bytez.Gibibyte` back.

See the [godoc](https://godoc.org/github.com/codexnull/bytez) for details.
