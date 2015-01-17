parfind
=======

[![GitHub release](http://img.shields.io/github/release/rakutentech/parfind.svg?style=flat-square)][release]
[![Travis](https://img.shields.io/travis/rakutentech/parfind.svg?style=flat-square)][travis]
[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godocs]

[release]: https://github.com/rakutentech/parfind/releases
[travis]: https://travis-ci.org/rakutentech/parfind
[godocs]: http://godoc.org/github.com/rakutentech/parfind

A parallel, simplified version of `find(1)` for use on high-latency,
highly-parallel file systems.

Usage
-----
`parfind -root=<directory>` will recursively list all files and directories in
`<directory>`. The output order is undefined. For each file/directory a line
will be printed to stdout:

    <type> <mtime> <size> <name>

where `<type>` is the type (e.g. `f` for regular files, `d` for directories, `l`
for symbolic links, etc.), `<mtime>` is the UNIX timestamp of the file/directory
mtime, `<size>` is the size in bytes, `<name>` is the quoted and escaped
absolute path. Fields are space separated and lines are terminated by a newline.

If the `-print0` option is used fields and lines are terminated by a `NUL` byte
and the `<name>` field is not quoted or escaped: this is useful for passing the
output to command line tools that support this convention, e.g. `xargs -0`:

    <type><NUL><mtime><NUL><size><NUL><name><NUL>

Build
-----

Just run below in the project directory:

```bash
$ go build
```

Contribution
----

1. Fork ([https://github.com/rakutentech/parfind/fork](https://github.com/rakutentech/parfind/fork))
1. Create a feature branch
1. Add features and its tests as appropriate
1. Commit your changes 
1. Rebase your local changes against the master branch
1. Run test suite with the `go test` command and confirm that it passes
1. Run `gofmt -s`
1. Create new Pull Request


Author
------
Carlo Alberto Ferraris (Rakuten, Inc.)
