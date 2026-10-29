package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/therealfakemoot/preroll/evaluator"
	"github.com/therealfakemoot/preroll/parser"
)

// table is a generated set of rolls plus the metadata needed to reproduce it.
type table struct {
	Roll    string
	Seed    int64
	Results []evaluator.Result
}

// gen implements "preroll gen". It returns the process exit code: 0 on
// success, 1 for bad expressions or write failures, 2 for bad usage.
func gen(args []string, stdout, stderr io.Writer) int {
	var (
		expr   string
		n      int
		seed   int64
		format string
	)

	fs := flag.NewFlagSet("gen", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.StringVar(&expr, "roll", "1d20+3", "dice roll")
	fs.IntVar(&n, "n", 100, "number of rolls to generate")
	fs.Int64Var(&seed, "seed", 0, "RNG seed; random (and printed in the output) if omitted")
	fs.StringVar(&format, "format", "csv", "output format: csv|json|md")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	seeded := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == "seed" {
			seeded = true
		}
	})
	if !seeded {
		// Kept below 2^53 so the seed survives a round trip through JSON
		// consumers that parse numbers as float64.
		seed = time.Now().UnixNano() & (1<<53 - 1)
	}

	write, ok := writers[format]
	if !ok {
		fmt.Fprintf(stderr, "preroll gen: unknown format %q (want csv, json or md)\n", format)
		return 2
	}
	if n < 1 {
		fmt.Fprintln(stderr, "preroll gen: -n must be at least 1")
		return 2
	}

	root, err := parser.Parse(expr, nil)
	if err != nil {
		fmt.Fprintf(stderr, "preroll gen: %v\n", err)
		return 1
	}

	// One RNG for the whole table keeps output a pure function of
	// (expression, seed, n).
	rng := rand.New(rand.NewSource(seed))
	t := table{Roll: expr, Seed: seed, Results: make([]evaluator.Result, 0, n)}
	for i := 0; i < n; i++ {
		res, err := evaluator.Eval(root, rng)
		if err != nil {
			fmt.Fprintf(stderr, "preroll gen: %v\n", err)
			return 1
		}
		t.Results = append(t.Results, res)
	}

	if err := write(stdout, t); err != nil {
		fmt.Fprintf(stderr, "preroll gen: %v\n", err)
		return 1
	}
	return 0
}
