// Package evaluator rolls a parsed dice expression using an injected RNG and
// reports every individual die, not just the total.
package evaluator

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/therealfakemoot/preroll/parser"
)

// DefaultExplodeCap bounds how many extra rolls a single die may trigger
// under a bare "!", so a 1-sided die cannot explode forever.
const DefaultExplodeCap = 100

// Source is the RNG the evaluator draws from. *rand.Rand satisfies it.
type Source interface {
	// Intn returns a value in [0, n).
	Intn(n int) int
}

// Die is a single rolled die.
type Die struct {
	// Value is the numeric value of the die. For a custom face that is not
	// a number it is 0 and does not count toward any total.
	Value int
	// Face is the face string for custom-faced dice, empty otherwise.
	Face string
	// Dropped is set when a keep/drop modifier excluded this die from the sum.
	Dropped bool
	// Exploded is set on dice rolled because an earlier die exploded.
	Exploded bool
}

// Term is one operand of the expression.
type Term struct {
	// Sign is +1 or -1, from the operator that preceded the term.
	Sign int
	// Dice is every die rolled for this term; empty for a constant.
	Dice []Die
	// Constant is the value of a literal number term.
	Constant int
	// Subtotal is the unsigned sum of the term's kept dice (or Constant).
	Subtotal int

	// nonNumeric is set when a custom face could not be read as a number.
	nonNumeric bool
}

// Result is the outcome of evaluating one expression once.
type Result struct {
	// Total is the sum of all signed term subtotals. It is 0 for a lone
	// roll whose custom faces are not numbers.
	Total int
	Terms []Term
	// Numeric is false when Total is meaningless because the expression is a
	// lone roll of non-numeric custom faces.
	Numeric bool
	// Faces holds the face strings of a lone custom-faced roll such as
	// 5d{red,blue,green}, in roll order. It is nil otherwise.
	Faces []string
}

// Eval rolls the expression once.
// Eval rolls the expression once, drawing from src in left-to-right order.
// The same node and RNG state always produce the same Result.
func Eval(n parser.Node, src Source) (Result, error) {
	var res Result
	if err := collect(n, 1, src, &res); err != nil {
		return Result{}, err
	}
	res.Numeric = true
	for _, t := range res.Terms {
		if t.nonNumeric {
			res.Numeric = false
		}
		if t.nonNumeric && len(res.Terms) > 1 {
			return Result{}, errors.New("non-numeric faces cannot be combined with other terms")
		}
		res.Total += t.Sign * t.Subtotal
	}

	if len(res.Terms) == 1 && res.Terms[0].isCustom() {
		for _, d := range res.Terms[0].Dice {
			res.Faces = append(res.Faces, d.Face)
		}
	}
	return res, nil
}

func (t Term) isCustom() bool {
	return len(t.Dice) > 0 && t.Dice[0].Face != ""
}

// collect flattens the left-leaning tree into ordered terms, rolling in
// source order so output is deterministic for a given RNG state.
func collect(n parser.Node, sign int, src Source, res *Result) error {
	switch n := n.(type) {
	case parser.Num:
		res.Terms = append(res.Terms, Term{Sign: sign, Constant: n.Value, Subtotal: n.Value})
	case parser.Roll:
		t, err := rollTerm(n, sign, src)
		if err != nil {
			return err
		}
		res.Terms = append(res.Terms, t)
	case parser.BinOp:
		if err := collect(n.Left, sign, src, res); err != nil {
			return err
		}
		rs := sign
		if n.Op == parser.Sub {
			rs = -sign
		}
		return collect(n.Right, rs, src, res)
	default:
		return fmt.Errorf("unknown node %T", n)
	}
	return nil
}

// rollTerm rolls one Roll node: it draws every die (and explosions), applies
// any keep/drop modifier, then sums the dice that were kept.
func rollTerm(r parser.Roll, sign int, src Source) (Term, error) {
	if r.Count < 1 || r.Sides < 1 {
		return Term{}, errors.New("roll needs at least one die and one side")
	}

	values := make([]int, r.Sides)
	numeric := true
	if r.Faces != nil {
		for i, f := range r.Faces {
			v, err := strconv.Atoi(f)
			if err != nil {
				numeric = false
			}
			values[i] = v
		}
	} else {
		for i := range values {
			values[i] = i + 1
		}
	}

	t := Term{Sign: sign}
	roll := func(exploded bool) (idx int) {
		idx = src.Intn(r.Sides)
		d := Die{Value: values[idx], Exploded: exploded}
		if r.Faces != nil {
			d.Face = r.Faces[idx]
		}
		t.Dice = append(t.Dice, d)
		return idx
	}

	for i := 0; i < r.Count; i++ {
		idx := roll(false)
		if r.Mod.Kind != parser.Explode {
			continue
		}
		limit := r.Mod.ExplodeMax
		if limit == 0 {
			limit = DefaultExplodeCap
		}
		// A die explodes whenever it shows its highest face.
		for extra := 0; extra < limit && idx == r.Sides-1; extra++ {
			idx = roll(true)
		}
	}

	applyKeepDrop(t.Dice, r.Mod.Kind)

	t.nonNumeric = !numeric
	for _, d := range t.Dice {
		if !d.Dropped {
			t.Subtotal += d.Value
		}
	}
	return t, nil
}

// applyKeepDrop marks dice as dropped. Keep and drop modifiers act on exactly
// one die: the highest or lowest. Ties go to the earliest die.
func applyKeepDrop(dice []Die, kind parser.ModKind) {
	var pick func(a, b int) bool // reports whether value a beats value b
	switch kind {
	case parser.KeepHighest, parser.DropHighest:
		pick = func(a, b int) bool { return a > b }
	case parser.KeepLowest, parser.DropLowest:
		pick = func(a, b int) bool { return a < b }
	default:
		return
	}

	best := 0
	for i := range dice {
		if pick(dice[i].Value, dice[best].Value) {
			best = i
		}
	}

	switch kind {
	case parser.KeepHighest, parser.KeepLowest:
		for i := range dice {
			dice[i].Dropped = i != best
		}
	default:
		dice[best].Dropped = true
	}
}
