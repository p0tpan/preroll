package main

import (
	"fmt"
	"io"
	"os"
)

const usage = `usage: preroll <command> [flags]

commands:
  gen    generate a table of pre-calculated rolls
`

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, usage)
		return 2
	}
	switch args[0] {
	case "gen":
		return gen(args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "preroll: unknown command %q\n\n%s", args[0], usage)
		return 2
	}
}
