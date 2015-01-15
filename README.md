parfind
=======
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
Just run `go build` in the project directory.

Sources
-------
Sources available on [GitHub](https://github.com/rakutentech/parfind).

License
-------
Copyright (c) 2014 Rakuten, Inc.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

Author
------
Carlo Alberto Ferraris (Rakuten, Inc.)
