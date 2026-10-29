package lexer

import (
	"fmt"
	"log/slog"
	"strings"
	"unicode/utf8"
)

type lexer struct {
	input  string
	start  int
	pos    int
	width  int
	tokens []Token
	states []stateFunc
	logger *slog.Logger
}

func (l *lexer) String() string {
	return fmt.Sprintf("{input:'%s' start:%d pos:%d width:%d}", l.input, l.start, l.pos, l.width)
}

func Lex(input string, logger *slog.Logger) *lexer {
	l := &lexer{
		input:  input,
		pos:    0,
		tokens: make([]Token, 0),
		states: make([]stateFunc, 0),
		logger: logger,
	}

	l.run()
	return l
}

// next returns the next rune in the input.
func (l *lexer) next() (rune rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return EOF
	}
	rune, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return rune
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

// acceptFn consumes a run of runes from the valid set.
func (l *lexer) acceptFn(matchers ...func(rune) bool) {
	// for strings.IndexRune(valid, l.next()) >= 0 {
	// }
	l.backup()
}

func (l *lexer) emit(t TokenType) {
	token := Token{t, l.input[l.start:l.pos]}
	l.tokens = append(l.tokens, token)
	l.start = l.pos
}

func (l *lexer) run() {
	logger := l.logger.WithGroup("run")
	logger.Debug("entering run")
	for state := lexModifier; state != nil; {
		state = state(l)
	}
	logger.Debug("terminating lexer")
}

func (l *lexer) Items() []Token {
	return l.tokens
}

func (l *lexer) errorf(format string, args ...any) stateFunc {
	l.tokens = append(l.tokens, Token{
		Type: ErrorToken,
		Raw:  fmt.Sprintf(format, args...),
	},
	)
	return nil
}
