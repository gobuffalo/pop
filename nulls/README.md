# github.com/markbates/pop/nulls

This package should be used in place of the built-in null types in the `sql` package.

The real benefit of this packages comes in its implementation of `MarshalJSON` and `UnmarshalJSON` to properly encode/decode `null` values.

## Installation

``` bash
$ go get github.com/markbates/pop/nulls
```

## Supported Datatypes

* `string` (`nulls.NullString`) - Replaces `sql.NullString`
* `int64` (`nulls.NullInt64`) - Replaces `sql.NullInt64`
* `float64` (`nulls.NullFloat64`) - Replaces `sql.NullFloat64`
* `bool` (`nulls.NullBool`) - Replaces `sql.NullBool`
* `[]byte` (`nulls.NullByteSlice`)
* `float32` (`nulls.NullFloat32`)
* `int` (`nulls.NullInt`)
* `int32` (`nulls.NullInt32`)
* `uint32` (`nulls.NullUInt32`)
* `time.Time` (`nulls.NullTime`)
