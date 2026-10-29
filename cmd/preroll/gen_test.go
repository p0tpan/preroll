package main

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func runGen(t *testing.T, args ...string) (string, int) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := run(append([]string{"gen"}, args...), &stdout, &stderr)
	if code != 0 {
		t.Logf("stderr: %s", stderr.String())
	}
	return stdout.String(), code
}

func Test_GenExact(t *testing.T) {
	cases := []struct {
		format   string
		expected string
	}{
		{"csv", "# roll=kh2d20+3 seed=42 n=3\n" +
			"index,total,dice\n" +
			"1,11,(6)|8\n" +
			"2,14,(9)|11\n" +
			"3,9,(4)|6\n"},
		{"md", "# preroll\n\n" +
			"- Roll: `kh2d20+3`\n- Seed: 42\n- Count: 3\n\n" +
			"| # | Total | Dice |\n|---|---|---|\n" +
			"| 1 | 11 | (6), 8 |\n" +
			"| 2 | 14 | (9), 11 |\n" +
			"| 3 | 9 | (4), 6 |\n"},
	}

	for _, tc := range cases {
		t.Run(tc.format, func(t *testing.T) {
			out, code := runGen(t, "-roll", "kh2d20+3", "-n", "3", "-seed", "42", "-format", tc.format)
			if code != 0 || out != tc.expected {
				t.Fatalf("code %d, got:\n%s", code, out)
			}
		})
	}
}

func Test_GenJSON(t *testing.T) {
	out, code := runGen(t, "-roll", "kh2d20+3", "-n", "3", "-seed", "42", "-format", "json")
	if code != 0 {
		t.Fatal(code)
	}
	var got struct {
		Roll  string `json:"roll"`
		Seed  int64  `json:"seed"`
		Count int    `json:"count"`
		Rolls []struct {
			Total int `json:"total"`
			Dice  []struct {
				Value   int  `json:"value"`
				Dropped bool `json:"dropped"`
			} `json:"dice"`
		} `json:"rolls"`
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatal(err)
	}
	if got.Roll != "kh2d20+3" || got.Seed != 42 || got.Count != 3 || len(got.Rolls) != 3 {
		t.Fatalf("bad metadata: %+v", got)
	}
	first := got.Rolls[0]
	if first.Total != 11 || len(first.Dice) != 2 || first.Dice[0].Value != 6 || !first.Dice[0].Dropped || first.Dice[1].Value != 8 {
		t.Fatalf("bad first roll: %+v", first)
	}
}

func Test_GenCustomFaces(t *testing.T) {
	out, code := runGen(t, "-roll", "3d{red,blue,green}", "-n", "2", "-seed", "42")
	if code != 0 {
		t.Fatal(code)
	}
	// Non-numeric faces have no total.
	if !strings.Contains(out, "\n1,,green|green|green\n") {
		t.Fatalf("got:\n%s", out)
	}
}

func Test_GenDeterministic(t *testing.T) {
	for _, format := range []string{"csv", "json", "md"} {
		a, _ := runGen(t, "-roll", "!3d6-1d4+2", "-n", "50", "-seed", "7", "-format", format)
		b, _ := runGen(t, "-roll", "!3d6-1d4+2", "-n", "50", "-seed", "7", "-format", format)
		c, _ := runGen(t, "-roll", "!3d6-1d4+2", "-n", "50", "-seed", "8", "-format", format)
		if a != b || a == c {
			t.Fatalf("%s: seed does not determine output", format)
		}
	}
}

// Test_GenSeedZero makes sure an explicit seed of 0 is honoured rather than
// treated as "unset".
func Test_GenSeedZero(t *testing.T) {
	out, _ := runGen(t, "-n", "1", "-seed", "0")
	if !strings.HasPrefix(out, "# roll=1d20+3 seed=0 n=1\n") {
		t.Fatalf("got:\n%s", out)
	}
}

func Test_GenRandomSeedReproduces(t *testing.T) {
	first, code := runGen(t, "-roll", "2d6", "-n", "20")
	if code != 0 {
		t.Fatal(code)
	}
	m := regexp.MustCompile(`^# roll=2d6 seed=(\d+) n=20\n`).FindStringSubmatch(first)
	if m == nil {
		t.Fatalf("no seed header:\n%s", first)
	}
	if _, err := strconv.ParseInt(m[1], 10, 64); err != nil {
		t.Fatal(err)
	}
	again, _ := runGen(t, "-roll", "2d6", "-n", "20", "-seed", m[1])
	if first != again {
		t.Fatalf("seed %s did not reproduce the table", m[1])
	}
}

func Test_GenErrors(t *testing.T) {
	cases := [][]string{
		{"-roll", "hello"},
		{"-roll", "1d{a b}"},
		{"-roll", "1d{red,blue}+3"},
		{"-n", "0"},
		{"-format", "pdf"},
		{"-bogus"},
	}
	for _, args := range cases {
		if out, code := runGen(t, args...); code == 0 || out != "" {
			t.Fatalf("%v: code %d, output %q", args, code, out)
		}
	}

	var stdout, stderr bytes.Buffer
	if code := run(nil, &stdout, &stderr); code != 2 {
		t.Fatal(code)
	}
	if code := run([]string{"nope"}, &stdout, &stderr); code != 2 {
		t.Fatal(code)
	}
}
