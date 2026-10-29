package lexer

import (
	"log/slog"
	"os"
	"reflect"
	"testing"
)

var logger = slog.New(
	slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		},
	),
)

func Test_Simple(t *testing.T) {
	logger := logger.With("test", "Test_Simple")
	cases := []struct {
		input    string
		expected []Token
	}{
		{
			input: "1d20", expected: []Token{
				{numberToken, "1"},
				{dieToken, "d"},
				{numberToken, "20"},
			},
		},
		{
			input: "2d4", expected: []Token{
				{numberToken, "2"},
				{dieToken, "d"},
				{numberToken, "4"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			logger := logger.With("case", tc.input)
			actual := Lex(tc.input, logger).Items()
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Fail()
			}
		})
	}
}

func Test_Dropping(t *testing.T) {
	logger := logger.With("test", "Test_Dropping")
	cases := []struct {
		input    string
		expected []Token
	}{
		{
			input: "dh1d20", expected: []Token{
				{dropHighestToken, "dh"},
				{numberToken, "1"},
				{dieToken, "d"},
				{numberToken, "20"},
			},
		},
		{
			input: "dl2d4", expected: []Token{
				{dropLowestToken, "dl"},
				{numberToken, "2"},
				{dieToken, "d"},
				{numberToken, "4"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			logger := logger.With("case", tc.input)
			actual := Lex(tc.input, logger).Items()
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Fail()
			}
		})
	}
}

func Test_Keeping(t *testing.T) {
	logger := logger.With("test", "Test_Keeping")
	cases := []struct {
		input    string
		expected []Token
	}{
		{
			input: "kh1d20", expected: []Token{
				{keepHighestToken, "kh"},
				{numberToken, "1"},
				{dieToken, "d"},
				{numberToken, "20"},
			},
		},
		{
			input: "kl2d4", expected: []Token{
				{keepLowestToken, "kl"},
				{numberToken, "2"},
				{dieToken, "d"},
				{numberToken, "4"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			logger := logger.With("case", tc.input)
			actual := Lex(tc.input, logger).Items()
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Fail()
			}
		})
	}
}

func Test_NumberSubtraction(t *testing.T) {
	logger := logger.With("test", "Test_NumberSubtraction")
	cases := []struct {
		input    string
		expected []Token
	}{
		{
			input: "1d20-4", expected: []Token{
				{numberToken, "1"},
				{dieToken, "d"},
				{numberToken, "20"},
				{subtractionToken, "-"},
				{numberToken, "4"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			logger := logger.With("case", tc.input)
			actual := Lex(tc.input, logger).Items()
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Fail()
			}
		})
	}
}

func Test_NumberAddition(t *testing.T) {
	logger := logger.With("test", "Test_NumberAddition")
	cases := []struct {
		input    string
		expected []Token
	}{
		{
			input: "1d20+4", expected: []Token{
				{numberToken, "1"},
				{dieToken, "d"},
				{numberToken, "20"},
				{additionToken, "+"},
				{numberToken, "4"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			logger := logger.With("case", tc.input)
			actual := Lex(tc.input, logger).Items()
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Fail()
			}
		})
	}
}

func Test_Faces(t *testing.T) {
	logger := logger.With("test", "Test_Faces")
	cases := []struct {
		input    string
		expected []Token
	}{
		{
			input: "1d{red,blue,green}", expected: []Token{
				{numberToken, "1"},
				{dieToken, "d"},
				{facesOpenToken, "{"},
				{faceToken, "red"},
				{facesSeparatorToken, ","},
				{faceToken, "blue"},
				{facesSeparatorToken, ","},
				{faceToken, "green"},
				{facesCloseToken, "}"},
			},
		},
		{
			input: "1d{2,7,194}", expected: []Token{
				{numberToken, "1"},
				{dieToken, "d"},
				{facesOpenToken, "{"},
				{faceToken, "2"},
				{facesSeparatorToken, ","},
				{faceToken, "7"},
				{facesSeparatorToken, ","},
				{faceToken, "194"},
				{facesCloseToken, "}"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			logger := logger.With("case", tc.input)
			actual := Lex(tc.input, logger).Items()
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Fail()
			}
		})
	}
}

func Test_Complex(t *testing.T) {
	logger := logger.With("test", "Test_Complex")
	cases := []struct {
		input    string
		expected []Token
	}{
		{
			input: "dh3d20-4", expected: []Token{
				{dropHighestToken, "dh"},
				{numberToken, "3"},
				{dieToken, "d"},
				{numberToken, "20"},
				{subtractionToken, "-"},
				{numberToken, "4"},
			},
		},
		{
			input: "kl2d69+4", expected: []Token{
				{keepLowestToken, "kl"},
				{numberToken, "2"},
				{dieToken, "d"},
				{numberToken, "69"},
				{additionToken, "+"},
				{numberToken, "4"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			logger := logger.With("case", tc.input)
			actual := Lex(tc.input, logger).Items()
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Fail()
			}
		})
	}
}

func Test_RollAdditionSimple(t *testing.T) {
	logger := logger.With("test", "Test_RollAdditionSimple")
	cases := []struct {
		input    string
		expected []Token
	}{
		{
			input: "3d20+1", expected: []Token{
				{numberToken, "3"},
				{dieToken, "d"},
				{numberToken, "20"},
				{additionToken, "+"},
				{numberToken, "1"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			logger := logger.With("case", tc.input)
			actual := Lex(tc.input, logger).Items()
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Fail()
			}
		})
	}
}

func Test_RollSubtractionSimple(t *testing.T) {
	logger := logger.With("test", "Test_RollAdditionSimple")
	cases := []struct {
		input    string
		expected []Token
	}{
		{
			input: "3d20-1", expected: []Token{
				{numberToken, "3"},
				{dieToken, "d"},
				{numberToken, "20"},
				{subtractionToken, "-"},
				{numberToken, "1"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			logger := logger.With("case", tc.input)
			actual := Lex(tc.input, logger).Items()
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Fail()
			}
		})
	}
}

func Test_RollAdditionComplex(t *testing.T) {
	logger := logger.With("test", "Test_RollAdditionComplex")
	cases := []struct {
		input    string
		expected []Token
	}{
		{
			input: "3d20+1d4", expected: []Token{
				{numberToken, "3"},
				{dieToken, "d"},
				{numberToken, "20"},
				{additionToken, "+"},
				{numberToken, "1"},
				{dieToken, "d"},
				{numberToken, "4"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			logger := logger.With("case", tc.input)
			actual := Lex(tc.input, logger).Items()
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Fail()
			}
		})
	}
}

func Test_RollSubtractionComplex(t *testing.T) {
	logger := logger.With("test", "Test_RollSubtractionComplex")
	cases := []struct {
		input    string
		expected []Token
	}{
		{
			input: "3d20-1d4", expected: []Token{
				{numberToken, "3"},
				{dieToken, "d"},
				{numberToken, "20"},
				{subtractionToken, "-"},
				{numberToken, "1"},
				{dieToken, "d"},
				{numberToken, "4"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			logger := logger.With("case", tc.input)
			actual := Lex(tc.input, logger).Items()
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Fail()
			}
		})
	}
}

func Test_RollComplexRolls(t *testing.T) {
	logger := logger.With("test", "Test_RollSubtractionComplex")
	cases := []struct {
		input    string
		expected []Token
	}{
		{
			input: "3d20-1d4+1d5", expected: []Token{
				{numberToken, "3"},
				{dieToken, "d"},
				{numberToken, "20"},
				{subtractionToken, "-"},
				{numberToken, "1"},
				{dieToken, "d"},
				{numberToken, "4"},
				{additionToken, "+"},
				{numberToken, "1"},
				{dieToken, "d"},
				{numberToken, "5"},
			},
		},
		{
			input: "3d20+1d4-1d5", expected: []Token{
				{numberToken, "3"},
				{dieToken, "d"},
				{numberToken, "20"},
				{additionToken, "+"},
				{numberToken, "1"},
				{dieToken, "d"},
				{numberToken, "4"},
				{subtractionToken, "-"},
				{numberToken, "1"},
				{dieToken, "d"},
				{numberToken, "5"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			logger := logger.With("case", tc.input)
			actual := Lex(tc.input, logger).Items()
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Fail()
			}
		})
	}
}

func Test_Explode(t *testing.T) {
	logger := logger.With("test", "Test_Simple")
	cases := []struct {
		input    string
		expected []Token
	}{
		{
			input: "!1d20", expected: []Token{
				{explodeToken, "!"},
				{numberToken, "1"},
				{dieToken, "d"},
				{numberToken, "20"},
			},
		},
		{
			input: "!{3}1d20", expected: []Token{
				{explodeToken, "!"},
				{explodeOpenToken, "{"},
				{explodeCountToken, "3"},
				{explodeCloseToken, "}"},
				{numberToken, "1"},
				{dieToken, "d"},
				{numberToken, "20"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			logger := logger.With("case", tc.input)
			actual := Lex(tc.input, logger).Items()
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Logf("%v\n", actual)
				t.Fail()
			}
		})
	}
}

func Test_Chained(t *testing.T) {
	logger := logger.With("test", "Test_Chained")
	cases := []struct {
		input    string
		expected []Token
	}{
		{
			input: "1d20+5+3", expected: []Token{
				{numberToken, "1"},
				{dieToken, "d"},
				{numberToken, "20"},
				{additionToken, "+"},
				{numberToken, "5"},
				{additionToken, "+"},
				{numberToken, "3"},
			},
		},
		{
			input: "1d20-5+2d4-1", expected: []Token{
				{numberToken, "1"},
				{dieToken, "d"},
				{numberToken, "20"},
				{subtractionToken, "-"},
				{numberToken, "5"},
				{additionToken, "+"},
				{numberToken, "2"},
				{dieToken, "d"},
				{numberToken, "4"},
				{subtractionToken, "-"},
				{numberToken, "1"},
			},
		},
		{
			input: "1d20+kh2d6+!{2}1d{a,b}", expected: []Token{
				{numberToken, "1"},
				{dieToken, "d"},
				{numberToken, "20"},
				{additionToken, "+"},
				{keepHighestToken, "kh"},
				{numberToken, "2"},
				{dieToken, "d"},
				{numberToken, "6"},
				{additionToken, "+"},
				{explodeToken, "!"},
				{explodeOpenToken, "{"},
				{explodeCountToken, "2"},
				{explodeCloseToken, "}"},
				{numberToken, "1"},
				{dieToken, "d"},
				{facesOpenToken, "{"},
				{faceToken, "a"},
				{facesSeparatorToken, ","},
				{faceToken, "b"},
				{facesCloseToken, "}"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			logger := logger.With("case", tc.input)
			actual := Lex(tc.input, logger).Items()
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Logf("%v\n", actual)
				t.Fail()
			}
		})
	}
}

func Test_Whitespace(t *testing.T) {
	logger := logger.With("test", "Test_Whitespace")
	cases := []struct {
		input    string
		expected []Token
	}{
		{
			input: "2d5 - 4", expected: []Token{
				{numberToken, "2"},
				{dieToken, "d"},
				{numberToken, "5"},
				{subtractionToken, "-"},
				{numberToken, "4"},
			},
		},
		{
			input: " kh2d20 + 3 ", expected: []Token{
				{keepHighestToken, "kh"},
				{numberToken, "2"},
				{dieToken, "d"},
				{numberToken, "20"},
				{additionToken, "+"},
				{numberToken, "3"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			logger := logger.With("case", tc.input)
			actual := Lex(tc.input, logger).Items()
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Logf("%v\n", actual)
				t.Fail()
			}
		})
	}
}

// Test_Errors checks that malformed input terminates and ends in an error
// token, preceded by whatever valid tokens were lexed before the problem.
func Test_Errors(t *testing.T) {
	logger := logger.With("test", "Test_Errors")
	cases := []struct {
		input  string
		prefix []Token
	}{
		{input: "1d{a b}", prefix: []Token{
			{numberToken, "1"}, {dieToken, "d"}, {facesOpenToken, "{"}, {faceToken, "a"},
		}},
		{input: "1d{}", prefix: []Token{
			{numberToken, "1"}, {dieToken, "d"}, {facesOpenToken, "{"},
		}},
		{input: "1d{a,}", prefix: []Token{
			{numberToken, "1"}, {dieToken, "d"}, {facesOpenToken, "{"}, {faceToken, "a"}, {facesSeparatorToken, ","},
		}},
		{input: "1d{red", prefix: []Token{
			{numberToken, "1"}, {dieToken, "d"}, {facesOpenToken, "{"}, {faceToken, "red"},
		}},
		{input: "0d20", prefix: []Token{}},
		{input: "0d0", prefix: []Token{}},
		{input: "1d0", prefix: []Token{{numberToken, "1"}, {dieToken, "d"}}},
		{input: "4d6xyz", prefix: []Token{{numberToken, "4"}, {dieToken, "d"}, {numberToken, "6"}}},
		{input: "hello", prefix: []Token{}},
		{input: "", prefix: []Token{}},
		{input: "d20", prefix: []Token{}},
		{input: "1d20+", prefix: []Token{
			{numberToken, "1"}, {dieToken, "d"}, {numberToken, "20"}, {additionToken, "+"},
		}},
		{input: "1d20+0", prefix: []Token{
			{numberToken, "1"}, {dieToken, "d"}, {numberToken, "20"}, {additionToken, "+"},
		}},
		{input: "1d20+5x", prefix: []Token{
			{numberToken, "1"}, {dieToken, "d"}, {numberToken, "20"}, {additionToken, "+"}, {numberToken, "5"},
		}},
		{input: "1d20 5", prefix: []Token{{numberToken, "1"}, {dieToken, "d"}, {numberToken, "20"}}},
		{input: "!{}1d20", prefix: []Token{{explodeToken, "!"}, {explodeOpenToken, "{"}}},
		{input: "!{3 1d20", prefix: []Token{
			{explodeToken, "!"}, {explodeOpenToken, "{"}, {explodeCountToken, "3"},
		}},
		{input: "kh1", prefix: []Token{{keepHighestToken, "kh"}, {numberToken, "1"}}},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			logger := logger.With("case", tc.input)
			actual := Lex(tc.input, logger).Items()
			if len(actual) != len(tc.prefix)+1 {
				t.Fatalf("got %v", actual)
			}
			if last := actual[len(actual)-1]; last.Type != errorToken || last.Raw == "" {
				t.Fatalf("last token is not an error: %v", actual)
			}
			if !reflect.DeepEqual(actual[:len(tc.prefix)], tc.prefix) {
				t.Fatalf("got %v", actual)
			}
		})
	}
}
