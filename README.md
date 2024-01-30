# ltsvcut [![PkgGoDev](https://pkg.go.dev/badge/github.com/hnakamur/go-ltsvcut)](https://pkg.go.dev/github.com/hnakamur/go-ltsvcut)

ltsvcut provides features to cut labels, values from an escaped LTSV (Labeled Tab Separated Values) line, and to unescape values.

Supported escapes are: \t -> tab, \n -> newline, \\ -> backslash.

Note: LTSV http://ltsv.org/ specification does not define escapes.
Escapes are our own extensions.
