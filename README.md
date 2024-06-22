# Gabagool™️ file format

I lost a full nights of sleep for this and so here it is.

`.gabagool` is a universal file format, right now it can store:
- Text - (like any other text file)
- Images - (in a format that is compatible with the go canvas)
- Binary data - (any data in raw bytes)

## Usage

simply `go get` the gabagool package:
```shell
go get github.com/kociumba/gabagool-file/gabagool
```
the `gabagool` subpackage in this repo is the library while

the `main` package will be a `.gabagool` interpreter for use with desktop environments.

## Opening a `.gabagool` file

This repo also houses the `.gabagool` language features which can be installed with:
```shell
go install -ldflags "-s -w -H windowsgui" 'github.com/kociumba/gabagool-file'
```
this installs a program capable of decoding and interpreting `.gabagool` files

This program should be added to your `$PATH` and set as the default for `.gabagool` files.
 
