# 🌳 Go Bonzai™ Composite Command Tree

## TODO for me
Must be compatible with MacOS (M1) and Linux. Windows? Meh.

- datetime with multiple formats
- epoch to datetime (+ vice versa)
- yaml to json (+ vice versa)
- multiple yaml to json (+ vice versa)
- create uuid
- random hex string

[![GoDoc](https://godoc.org/github.com/coopernetes/bz?status.svg)](https://godoc.org/github.com/coopernetes/bz)
[![License](https://img.shields.io/badge/license-Apache2-brightgreen.svg)](LICENSE)

## Install

This command can be installed as a standalone program or composed into a
Bonzai command tree.

Standalone

```
go install github.com/coopernetes/bz/cmd/bz@latest
```

Composed

```go
package z

import (
	Z "github.com/coopernetes/bonzai/z"
	example "github.com/coopernetes/bz"
)

var Cmd = &Z.Cmd{
	Name:     `z`,
	Commands: []*Z.Cmd{help.Cmd, example.Cmd, example.BazCmd},
}
```

## Tab Completion

To activate bash completion just use the `complete -C` option from your
`.bashrc` or command line. There is no messy sourcing required. All the
completion is done by the program itself.

```
complete -C bz bz
```

If you don't have bash or tab completion check use the shortcut
commands instead.

## Embedded Documentation

All documentation (like manual pages) has been embedded into the source
code of the application. See the source or run the program with help to
access it.

