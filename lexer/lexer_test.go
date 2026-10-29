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
				{NumberToken, "1"},
				{DieToken, "d"},
				{NumberToken, "20"},
			},
		},
		{
			input: "2d4", expected: []Token{
				{NumberToken, "2"},
				{DieToken, "d"},
				{NumberToken, "4"},
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
				{DropHighestToken, "dh"},
				{NumberToken, "1"},
				{DieToken, "d"},
				{NumberToken, "20"},
			},
		},
		{
			input: "dl2d4", expected: []Token{
				{DropLowestToken, "dl"},
				{NumberToken, "2"},
				{DieToken, "d"},
				{NumberToken, "4"},
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
				{KeepHighestToken, "kh"},
				{NumberToken, "1"},
				{DieToken, "d"},
				{NumberToken, "20"},
			},
		},
		{
			input: "kl2d4", expected: []Token{
				{KeepLowestToken, "kl"},
				{NumberToken, "2"},
				{DieToken, "d"},
				{NumberToken, "4"},
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
				{NumberToken, "1"},
				{DieToken, "d"},
				{NumberToken, "20"},
				{SubtractionToken, "-"},
				{NumberToken, "4"},
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
				{NumberToken, "1"},
				{DieToken, "d"},
				{NumberToken, "20"},
				{AdditionToken, "+"},
				{NumberToken, "4"},
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
				{NumberToken, "1"},
				{DieToken, "d"},
				{FacesOpenToken, "{"},
				{FaceToken, "red"},
				{FacesSeparatorToken, ","},
				{FaceToken, "blue"},
				{FacesSeparatorToken, ","},
				{FaceToken, "green"},
				{FacesCloseToken, "}"},
			},
		},
		{
			input: "1d{2,7,194}", expected: []Token{
				{NumberToken, "1"},
				{DieToken, "d"},
				{FacesOpenToken, "{"},
				{FaceToken, "2"},
				{FacesSeparatorToken, ","},
				{FaceToken, "7"},
				{FacesSeparatorToken, ","},
				{FaceToken, "194"},
				{FacesCloseToken, "}"},
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
				{DropHighestToken, "dh"},
				{NumberToken, "3"},
				{DieToken, "d"},
				{NumberToken, "20"},
				{SubtractionToken, "-"},
				{NumberToken, "4"},
			},
		},
		{
			input: "kl2d69+4", expected: []Token{
				{KeepLowestToken, "kl"},
				{NumberToken, "2"},
				{DieToken, "d"},
				{NumberToken, "69"},
				{AdditionToken, "+"},
				{NumberToken, "4"},
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
				{NumberToken, "3"},
				{DieToken, "d"},
				{NumberToken, "20"},
				{AdditionToken, "+"},
				{NumberToken, "1"},
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
				{NumberToken, "3"},
				{DieToken, "d"},
				{NumberToken, "20"},
				{SubtractionToken, "-"},
				{NumberToken, "1"},
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
				{NumberToken, "3"},
				{DieToken, "d"},
				{NumberToken, "20"},
				{AdditionToken, "+"},
				{NumberToken, "1"},
				{DieToken, "d"},
				{NumberToken, "4"},
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
				{NumberToken, "3"},
				{DieToken, "d"},
				{NumberToken, "20"},
				{SubtractionToken, "-"},
				{NumberToken, "1"},
				{DieToken, "d"},
				{NumberToken, "4"},
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
				{NumberToken, "3"},
				{DieToken, "d"},
				{NumberToken, "20"},
				{SubtractionToken, "-"},
				{NumberToken, "1"},
				{DieToken, "d"},
				{NumberToken, "4"},
				{AdditionToken, "+"},
				{NumberToken, "1"},
				{DieToken, "d"},
				{NumberToken, "5"},
			},
		},
		{
			input: "3d20+1d4-1d5", expected: []Token{
				{NumberToken, "3"},
				{DieToken, "d"},
				{NumberToken, "20"},
				{AdditionToken, "+"},
				{NumberToken, "1"},
				{DieToken, "d"},
				{NumberToken, "4"},
				{SubtractionToken, "-"},
				{NumberToken, "1"},
				{DieToken, "d"},
				{NumberToken, "5"},
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
				{ExplodeToken, "!"},
				{NumberToken, "1"},
				{DieToken, "d"},
				{NumberToken, "20"},
			},
		},
		{
			input: "!{3}1d20", expected: []Token{
				{ExplodeToken, "!"},
				{ExplodeOpenToken, "{"},
				{ExplodeCountToken, "3"},
				{ExplodeCloseToken, "}"},
				{NumberToken, "1"},
				{DieToken, "d"},
				{NumberToken, "20"},
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
				{NumberToken, "1"},
				{DieToken, "d"},
				{NumberToken, "20"},
				{AdditionToken, "+"},
				{NumberToken, "5"},
				{AdditionToken, "+"},
				{NumberToken, "3"},
			},
		},
		{
			input: "1d20-5+2d4-1", expected: []Token{
				{NumberToken, "1"},
				{DieToken, "d"},
				{NumberToken, "20"},
				{SubtractionToken, "-"},
				{NumberToken, "5"},
				{AdditionToken, "+"},
				{NumberToken, "2"},
				{DieToken, "d"},
				{NumberToken, "4"},
				{SubtractionToken, "-"},
				{NumberToken, "1"},
			},
		},
		{
			input: "1d20+kh2d6+!{2}1d{a,b}", expected: []Token{
				{NumberToken, "1"},
				{DieToken, "d"},
				{NumberToken, "20"},
				{AdditionToken, "+"},
				{KeepHighestToken, "kh"},
				{NumberToken, "2"},
				{DieToken, "d"},
				{NumberToken, "6"},
				{AdditionToken, "+"},
				{ExplodeToken, "!"},
				{ExplodeOpenToken, "{"},
				{ExplodeCountToken, "2"},
				{ExplodeCloseToken, "}"},
				{NumberToken, "1"},
				{DieToken, "d"},
				{FacesOpenToken, "{"},
				{FaceToken, "a"},
				{FacesSeparatorToken, ","},
				{FaceToken, "b"},
				{FacesCloseToken, "}"},
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
				{NumberToken, "2"},
				{DieToken, "d"},
				{NumberToken, "5"},
				{SubtractionToken, "-"},
				{NumberToken, "4"},
			},
		},
		{
			input: " kh2d20 + 3 ", expected: []Token{
				{KeepHighestToken, "kh"},
				{NumberToken, "2"},
				{DieToken, "d"},
				{NumberToken, "20"},
				{AdditionToken, "+"},
				{NumberToken, "3"},
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
			{NumberToken, "1"}, {DieToken, "d"}, {FacesOpenToken, "{"}, {FaceToken, "a"},
		}},
		{input: "1d{}", prefix: []Token{
			{NumberToken, "1"}, {DieToken, "d"}, {FacesOpenToken, "{"},
		}},
		{input: "1d{a,}", prefix: []Token{
			{NumberToken, "1"}, {DieToken, "d"}, {FacesOpenToken, "{"}, {FaceToken, "a"}, {FacesSeparatorToken, ","},
		}},
		{input: "1d{red", prefix: []Token{
			{NumberToken, "1"}, {DieToken, "d"}, {FacesOpenToken, "{"}, {FaceToken, "red"},
		}},
		{input: "0d20", prefix: []Token{}},
		{input: "0d0", prefix: []Token{}},
		{input: "1d0", prefix: []Token{{NumberToken, "1"}, {DieToken, "d"}}},
		{input: "4d6xyz", prefix: []Token{{NumberToken, "4"}, {DieToken, "d"}, {NumberToken, "6"}}},
		{input: "hello", prefix: []Token{}},
		{input: "", prefix: []Token{}},
		{input: "d20", prefix: []Token{}},
		{input: "1d20+", prefix: []Token{
			{NumberToken, "1"}, {DieToken, "d"}, {NumberToken, "20"}, {AdditionToken, "+"},
		}},
		{input: "1d20+0", prefix: []Token{
			{NumberToken, "1"}, {DieToken, "d"}, {NumberToken, "20"}, {AdditionToken, "+"},
		}},
		{input: "1d20+5x", prefix: []Token{
			{NumberToken, "1"}, {DieToken, "d"}, {NumberToken, "20"}, {AdditionToken, "+"}, {NumberToken, "5"},
		}},
		{input: "1d20 5", prefix: []Token{{NumberToken, "1"}, {DieToken, "d"}, {NumberToken, "20"}}},
		{input: "!{}1d20", prefix: []Token{{ExplodeToken, "!"}, {ExplodeOpenToken, "{"}}},
		{input: "!{3 1d20", prefix: []Token{
			{ExplodeToken, "!"}, {ExplodeOpenToken, "{"}, {ExplodeCountToken, "3"},
		}},
		{input: "kh1", prefix: []Token{{KeepHighestToken, "kh"}, {NumberToken, "1"}}},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			logger := logger.With("case", tc.input)
			actual := Lex(tc.input, logger).Items()
			if len(actual) != len(tc.prefix)+1 {
				t.Fatalf("got %v", actual)
			}
			if last := actual[len(actual)-1]; last.Type != ErrorToken || last.Raw == "" {
				t.Fatalf("last token is not an error: %v", actual)
			}
			if !reflect.DeepEqual(actual[:len(tc.prefix)], tc.prefix) {
				t.Fatalf("got %v", actual)
			}
		})
	}
}
