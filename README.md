# Advent of Code solutions

This project holds my Go solutions + utilities for **[Advent of Code](https://adventofcode.com/)** puzzles over the years.

## Prereqs
To run solutions and tests, the following must be true:
1. The `$GOPATH` env var is set
2. This repository is cloned into `$GOPATH/src/github.com/jzimbel/adventofcode-go`
3. Your `$PATH` includes `$GOPATH/bin`

## Run it
```sh
$ go install github.com/jzimbel/adventofcode-go
$ adventofcode-go <year> <day>
```

If the input for the solution you're trying to run hasn't already been saved, the program will try to download it from the Advent of Code site first. If this is your first time downloading an input, you'll be asked to provide your unique session id. It's held in a cookie named `session` saved by the site--you can view it using your browser's dev tools or a number of cookie-viewing browser extensions.

## Test it
```sh
$ go test github.com/jzimbel/adventofcode-go
```

## Copy it
If you like the setup I've got here, feel free to use it yourself. Or let me know if it's garbage! This is my first time making a proper Go project from scratch, so I'm sure it's not perfect. ¯\\\_(ツ)\_/¯
