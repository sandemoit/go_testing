# TEKNALOGI

[![Go Reference](https://pkg.go.dev/badge/golang.org/x/example.svg)](https://pkg.go.dev/golang.org/x/example)

This repository contains a collection of Go programs and libraries that
demonstrate the language, standard libraries, and tools.

## Clone the project

```
$ git clone https://gitlab.pusri.co.id/teknologi/be-teknologi.git
$ cd example
```

## [hello](hello/) and [hello/reverse](hello/reverse/)

```
$ cd hello
$ go build
$ ./hello -help
```

A trivial "Hello, world" program that uses a library package.

The [hello](hello/) command covers:

- The basic form of an executable command
- Importing packages (from the standard library and the local repository)
- Printing strings ([fmt](//golang.org/pkg/fmt/))
- Command-line flags ([flag](//golang.org/pkg/flag/))
- Logging ([log](//golang.org/pkg/log/))

The [reverse](hello/reverse/) reverse covers:

- The basic form of a library
- Conversion between string and []rune
- Table-driven unit tests ([testing](//golang.org/pkg/testing/))
