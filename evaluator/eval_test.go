package evaluator

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/therealfakemoot/preroll/parser"
)

// script returns a fixed sequence of Intn results, panicking if exhausted.
type script []int

func (s *script) Intn(n int) int {
	v := (*s)[0]
	*s = (*s)[1:]
	if v >= n {
		panic("scripted value out of range")
	}
	return v
}

func mustParse(t *testing.T, input string) parser.Node {
	t.Helper()
	n, err := parser.Parse(input, nil)
	if err != nil {
		t.Fatal(err)
	}
	return n
}

func values(dice []Die) (vals []int, dropped []bool) {
	for _, d := range dice {
		vals = append(vals, d.Value)
		dropped = append(dropped, d.Dropped)
	}
	return vals, dropped
}

// Test_Scripted exercises modifier logic with fully controlled rolls. A
// scripted index i is the face with value i+1 on a numeric die.
func Test_Scripted(t *testing.T) {
	cases := []struct {
		input   string
		rolls   script
		total   int
		values  []int
		dropped []bool
	}{
		{"2d6", script{2, 4}, 8, []int{3, 5}, []bool{false, false}},
		{"1d20+3", script{9}, 13, []int{10}, []bool{false}},
		{"1d20-3", script{9}, 7, []int{10}, []bool{false}},
		{"kh2d20", script{4, 14}, 15, []int{5, 15}, []bool{true, false}},
		{"kl2d20", script{4, 14}, 5, []int{5, 15}, []bool{false, true}},
		{"dl4d6", script{0, 5, 2, 3}, 13, []int{1, 6, 3, 4}, []bool{true, false, false, false}},
		{"dh4d6", script{0, 5, 2, 3}, 8, []int{1, 6, 3, 4}, []bool{false, true, false, false}},
		// Ties drop/keep the earliest die.
		{"kh2d6", script{3, 3}, 4, []int{4, 4}, []bool{false, true}},
		// A die showing its max rolls again; chains continue.
		{"!2d4", script{3, 3, 1, 0}, 4 + 4 + 2 + 1, []int{4, 4, 2, 1}, []bool{false, false, false, false}},
		// !{N} limits extra rolls per die.
		{"!{1}1d2", script{1, 1}, 4, []int{2, 2}, []bool{false, false}},
		{"1d20+5+3", script{0}, 9, []int{1}, []bool{false}},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			res, err := Eval(mustParse(t, tc.input), &tc.rolls)
			if err != nil {
				t.Fatal(err)
			}
			vals, dropped := values(res.Terms[0].Dice)
			if res.Total != tc.total || !reflect.DeepEqual(vals, tc.values) || !reflect.DeepEqual(dropped, tc.dropped) {
				t.Fatalf("got total=%d values=%v dropped=%v", res.Total, vals, dropped)
			}
			if len(tc.rolls) != 0 {
				t.Fatalf("unused scripted rolls: %v", []int(tc.rolls))
			}
		})
	}
}

func Test_ExplodeMarksExtraDice(t *testing.T) {
	rolls := script{1, 0}
	res, err := Eval(mustParse(t, "!1d2"), &rolls)
	if err != nil {
		t.Fatal(err)
	}
	want := []Die{{Value: 2}, {Value: 1, Exploded: true}}
	if !reflect.DeepEqual(res.Terms[0].Dice, want) {
		t.Fatalf("got %+v", res.Terms[0].Dice)
	}
}

func Test_ExplodeDefaultCap(t *testing.T) {
	// A one-sided die always shows its max; the cap must stop it.
	res, err := Eval(mustParse(t, "!1d1"), rand.New(rand.NewSource(1)))
	if err != nil {
		t.Fatal(err)
	}
	if got := len(res.Terms[0].Dice); got != DefaultExplodeCap+1 {
		t.Fatalf("got %d dice", got)
	}
}

func Test_Terms(t *testing.T) {
	rolls := script{0, 1}
	res, err := Eval(mustParse(t, "1d4-1d4+7"), &rolls)
	if err != nil {
		t.Fatal(err)
	}
	signs := []int{}
	for _, term := range res.Terms {
		signs = append(signs, term.Sign)
	}
	if res.Total != 1-2+7 || !reflect.DeepEqual(signs, []int{1, -1, 1}) {
		t.Fatalf("got %+v", res)
	}
}

func Test_CustomFaces(t *testing.T) {
	res, err := Eval(mustParse(t, "3d{red,blue,green}"), &script{2, 0, 1})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"green", "red", "blue"}; !reflect.DeepEqual(res.Faces, want) {
		t.Fatalf("got %v", res.Faces)
	}

	// Numeric custom faces are summed.
	res, err = Eval(mustParse(t, "2d{1,3,17}+1"), &script{2, 1})
	if err != nil {
		t.Fatal(err)
	}
	if res.Total != 17+3+1 || res.Faces != nil {
		t.Fatalf("got %+v", res)
	}
}

func Test_NonNumericFacesCannotBeCombined(t *testing.T) {
	if _, err := Eval(mustParse(t, "1d{red,blue}+3"), &script{0}); err == nil {
		t.Fatal("expected error")
	}
}

// Test_Seeded pins exact output for a seeded math/rand source. math/rand's
// seeded sequence is stable across Go versions.
func Test_Seeded(t *testing.T) {
	cases := []struct {
		input   string
		total   int
		values  []int
		dropped []bool
	}{
		{"3d6", 15, []int{6, 6, 3}, []bool{false, false, false}},
		{"kh2d20+3", 11, []int{6, 8}, []bool{true, false}},
		{"dl4d6", 15, []int{6, 6, 3, 1}, []bool{false, false, false, true}},
		{"!2d4", 7, []int{2, 4, 1}, []bool{false, false, false}},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			res, err := Eval(mustParse(t, tc.input), rand.New(rand.NewSource(42)))
			if err != nil {
				t.Fatal(err)
			}
			vals, dropped := values(res.Terms[0].Dice)
			if res.Total != tc.total || !reflect.DeepEqual(vals, tc.values) || !reflect.DeepEqual(dropped, tc.dropped) {
				t.Fatalf("got total=%d values=%v dropped=%v", res.Total, vals, dropped)
			}
		})
	}

	res, err := Eval(mustParse(t, "5d{red,blue,green}"), rand.New(rand.NewSource(42)))
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"green", "green", "green", "red", "blue"}; !reflect.DeepEqual(res.Faces, want) {
		t.Fatalf("got %v", res.Faces)
	}
}
