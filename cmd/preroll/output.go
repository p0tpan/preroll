package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/therealfakemoot/preroll/evaluator"
)

// writers maps each -format value to its renderer.
var writers = map[string]func(io.Writer, table) error{
	"csv":  writeCSV,
	"json": writeJSON,
	"md":   writeMarkdown,
}

// dieText renders one die: its face or value, with dropped dice in
// parentheses.
func dieText(d evaluator.Die) string {
	s := d.Face
	if s == "" {
		s = strconv.Itoa(d.Value)
	}
	if d.Dropped {
		return "(" + s + ")"
	}
	return s
}

func diceText(r evaluator.Result, sep string) string {
	var parts []string
	for _, t := range r.Terms {
		for _, d := range t.Dice {
			parts = append(parts, dieText(d))
		}
	}
	return strings.Join(parts, sep)
}

// totalText is blank when the roll has no numeric total (custom faces).
func totalText(r evaluator.Result) string {
	if !r.Numeric {
		return ""
	}
	return strconv.Itoa(r.Total)
}

// writeCSV emits a "# roll=… seed=… n=…" comment line so the table can be
// regenerated, followed by a regular CSV body.
func writeCSV(w io.Writer, t table) error {
	if _, err := fmt.Fprintf(w, "# roll=%s seed=%d n=%d\n", t.Roll, t.Seed, len(t.Results)); err != nil {
		return err
	}
	cw := csv.NewWriter(w)
	cw.Write([]string{"index", "total", "dice"})
	for i, r := range t.Results {
		cw.Write([]string{strconv.Itoa(i + 1), totalText(r), diceText(r, "|")})
	}
	cw.Flush()
	return cw.Error()
}

func writeMarkdown(w io.Writer, t table) error {
	var b strings.Builder
	fmt.Fprintf(&b, "# preroll\n\n- Roll: `%s`\n- Seed: %d\n- Count: %d\n\n", t.Roll, t.Seed, len(t.Results))
	b.WriteString("| # | Total | Dice |\n|---|---|---|\n")
	for i, r := range t.Results {
		fmt.Fprintf(&b, "| %d | %s | %s |\n", i+1, totalText(r), diceText(r, ", "))
	}
	_, err := io.WriteString(w, b.String())
	return err
}

type jsonDie struct {
	Value    int    `json:"value"`
	Face     string `json:"face,omitempty"`
	Dropped  bool   `json:"dropped,omitempty"`
	Exploded bool   `json:"exploded,omitempty"`
}

type jsonRoll struct {
	Total *int      `json:"total"`
	Dice  []jsonDie `json:"dice"`
}

type jsonTable struct {
	Roll  string     `json:"roll"`
	Seed  int64      `json:"seed"`
	Count int        `json:"count"`
	Rolls []jsonRoll `json:"rolls"`
}

// writeJSON emits metadata at the top level next to the rolls. A roll with no
// numeric total has a null "total".
func writeJSON(w io.Writer, t table) error {
	out := jsonTable{Roll: t.Roll, Seed: t.Seed, Count: len(t.Results), Rolls: make([]jsonRoll, 0, len(t.Results))}
	for _, r := range t.Results {
		jr := jsonRoll{Dice: []jsonDie{}}
		if r.Numeric {
			total := r.Total
			jr.Total = &total
		}
		for _, term := range r.Terms {
			for _, d := range term.Dice {
				jr.Dice = append(jr.Dice, jsonDie{d.Value, d.Face, d.Dropped, d.Exploded})
			}
		}
		out.Rolls = append(out.Rolls, jr)
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
