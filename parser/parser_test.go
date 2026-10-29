package parser

import (
	"reflect"
	"testing"
)

func Test_Parse(t *testing.T) {
	cases := []struct {
		input    string
		expected Node
	}{
		{"1d20", Roll{Count: 1, Sides: 20}},
		{"13d3", Roll{Count: 13, Sides: 3}},
		{"2d5 - 4", BinOp{Sub, Roll{Count: 2, Sides: 5}, Num{4}}},
		{"1d20+3", BinOp{Add, Roll{Count: 1, Sides: 20}, Num{3}}},
		{"5d{red,blue,green}", Roll{Count: 5, Sides: 3, Faces: []string{"red", "blue", "green"}}},
		{"2d{1,3,17}", Roll{Count: 2, Sides: 3, Faces: []string{"1", "3", "17"}}},
		{"1d{solo}", Roll{Count: 1, Sides: 1, Faces: []string{"solo"}}},
		{"kh2d20 + 3", BinOp{Add, Roll{Count: 2, Sides: 20, Mod: Modifier{Kind: KeepHighest}}, Num{3}}},
		{"kl2d20", Roll{Count: 2, Sides: 20, Mod: Modifier{Kind: KeepLowest}}},
		{"dh3d6", Roll{Count: 3, Sides: 6, Mod: Modifier{Kind: DropHighest}}},
		{"dl3d13 - 4", BinOp{Sub, Roll{Count: 3, Sides: 13, Mod: Modifier{Kind: DropLowest}}, Num{4}}},
		{"!1d6", Roll{Count: 1, Sides: 6, Mod: Modifier{Kind: Explode}}},
		{"!{3}2d6", Roll{Count: 2, Sides: 6, Mod: Modifier{Kind: Explode, ExplodeMax: 3}}},
		{"3d20+1d4", BinOp{Add, Roll{Count: 3, Sides: 20}, Roll{Count: 1, Sides: 4}}},
		{"1d8+kh2d6", BinOp{Add, Roll{Count: 1, Sides: 8}, Roll{Count: 2, Sides: 6, Mod: Modifier{Kind: KeepHighest}}}},
		// Chained terms are left-associative.
		{"10d6-3-2", BinOp{Sub, BinOp{Sub, Roll{Count: 10, Sides: 6}, Num{3}}, Num{2}}},
		{"3d20-1d4+1d5", BinOp{Add,
			BinOp{Sub, Roll{Count: 3, Sides: 20}, Roll{Count: 1, Sides: 4}},
			Roll{Count: 1, Sides: 5},
		}},
		{"1d20+5+3", BinOp{Add, BinOp{Add, Roll{Count: 1, Sides: 20}, Num{5}}, Num{3}}},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			actual, err := Parse(tc.input, nil)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Fatalf("got %#v", actual)
			}
		})
	}
}

func Test_ParseErrors(t *testing.T) {
	cases := []string{
		"",
		"hello",
		"0d20",
		"1d0",
		"4d6xyz",
		"1d{a b}",
		"1d20+",
		"5+1d20",                 // the grammar needs a leading roll
		"5",                      // a lone number is not a roll
		"99999999999999999999d6", // overflows int
		"kh1d{a,b}",              // keep/drop need numeric dice
		"!2d{a,b}",               // so does exploding
	}

	for _, input := range cases {
		t.Run(input, func(t *testing.T) {
			if n, err := Parse(input, nil); err == nil {
				t.Fatalf("expected error, got %#v", n)
			}
		})
	}
}
