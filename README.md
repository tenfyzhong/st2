# st2
[![test](https://github.com/tenfyzhong/st2/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/tenfyzhong/st2/actions/workflows/test.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/tenfyzhong/st2.svg)](https://pkg.go.dev/github.com/tenfyzhong/st2)

`st2` provide a package to parse json/protobuf/thrift/go/csv code and generage go/protobuf/thrift code.

## Cli
`st2` provide a terminal command line tool `st2`, which can be used to generate go/protobuf/thrift code from json/protobuf/thrift/go/csv code.

### Install
####  Use home brew
```bash
brew tap tenfyzhong/tap
brew install st2
```

#### Download from release
You can download the release of `st2` from the [GitHub releases](https://github.com/tenfyzhong/st2/releases).  

#### build from source
```bash
go install github.com/tenfyzhong/st2/cmd/st2@latest
```

### Usage
```
NAME:
   st2 - convert between json, protobuf, thrift, go struct

USAGE:
   st2 [global options] [arguments...]

VERSION:
   developing

AUTHOR:
   tenfyzhong <tenfy@tenfy.cn>

GLOBAL OPTIONS:
   common

   --root name, -r name  The root struct name (default: Root)

   input

   --input file, -i file  Input file, if not set, it will read from stdio
   --rc                   Read input from clipboard (default: false)
   --src type, -s type    The source data type, it will use the suffix of the input file if not set, available value: `[json,proto,thrift,go,csv]`

   output

   --dst type, -d type     The destination data type, it will use the suffix of the output file if not set, available value: `[go,proto,thrift]`
   --output file, -o file  Output file, if not set, it will write to stdout
   --wc                    Write output to clipboard (default: false)


COPYRIGHT:
   Copyright (c) 2022 tenfy
```
