package parser

import (
	"fmt"
	"io"
	"log/slog"
	"strconv"

	"github.com/therealfakemoot/preroll/lexer"
)

type parser struct {
	tokens []lexer.Token
	pos    int
}

// Parse lexes and parses a dice expression. A nil logger discards lexer
// logging.
func Parse(input string, logger *slog.Logger) (Node, error) {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}
	p := &parser{tokens: lexer.Lex(input, logger).Items()}
	return p.parse()
}

func (p *parser) parse() (Node, error) {
	// The lexer stops at the first error and reports it as the last token.
	if n := len(p.tokens); n > 0 && p.tokens[n-1].Type == lexer.ErrorToken {
		return nil, fmt.Errorf("lex error: %s", p.tokens[n-1].Raw)
	}

	left, err := p.roll()
	if err != nil {
		return nil, err
	}
	for p.more() {
		var op Op
		switch p.peek().Type {
		case lexer.AdditionToken:
			op = Add
		case lexer.SubtractionToken:
			op = Sub
		default:
			return nil, p.errorf("expected + or -, got %q", p.peek().Raw)
		}
		p.pos++

		right, err := p.operand()
		if err != nil {
			return nil, err
		}
		left = BinOp{Op: op, Left: left, Right: right}
	}
	return left, nil
}

// operand parses the right-hand side of an operator: a bare number or a roll.
func (p *parser) operand() (Node, error) {
	if !p.more() {
		return nil, p.errorf("expected number or roll after operator")
	}
	// A number followed by a die is a roll's count, not a constant.
	if p.peek().Type == lexer.NumberToken &&
		(p.pos+1 >= len(p.tokens) || p.tokens[p.pos+1].Type != lexer.DieToken) {
		v, err := p.number()
		if err != nil {
			return nil, err
		}
		return Num{Value: v}, nil
	}
	return p.roll()
}

func (p *parser) roll() (Node, error) {
	var r Roll

	if p.more() {
		switch p.peek().Type {
		case lexer.KeepHighestToken:
			r.Mod.Kind = KeepHighest
		case lexer.KeepLowestToken:
			r.Mod.Kind = KeepLowest
		case lexer.DropHighestToken:
			r.Mod.Kind = DropHighest
		case lexer.DropLowestToken:
			r.Mod.Kind = DropLowest
		case lexer.ExplodeToken:
			r.Mod.Kind = Explode
		}
		if r.Mod.Kind != ModNone {
			p.pos++
		}
	}
	if r.Mod.Kind == Explode && p.more() && p.peek().Type == lexer.ExplodeOpenToken {
		p.pos++
		if !p.more() || p.peek().Type != lexer.ExplodeCountToken {
			return nil, p.errorf("expected explode count")
		}
		max, err := p.number()
		if err != nil {
			return nil, err
		}
		r.Mod.ExplodeMax = max
		if !p.more() || p.peek().Type != lexer.ExplodeCloseToken {
			return nil, p.errorf("expected closing }")
		}
		p.pos++
	}

	var err error
	if r.Count, err = p.expect(lexer.NumberToken, "die count"); err != nil {
		return nil, err
	}
	if !p.more() || p.peek().Type != lexer.DieToken {
		return nil, p.errorf("expected d")
	}
	p.pos++

	if p.more() && p.peek().Type == lexer.FacesOpenToken {
		p.pos++
		for {
			if !p.more() || p.peek().Type != lexer.FaceToken {
				return nil, p.errorf("expected face")
			}
			r.Faces = append(r.Faces, p.peek().Raw)
			p.pos++
			if !p.more() {
				return nil, p.errorf("unterminated faces list")
			}
			t := p.peek().Type
			p.pos++
			if t == lexer.FacesCloseToken {
				break
			}
			if t != lexer.FacesSeparatorToken {
				return nil, p.errorf("expected , or } in faces list")
			}
		}
		r.Sides = len(r.Faces)
		if r.Mod.Kind != ModNone {
			return nil, p.errorf("keep, drop and explode modifiers need numeric dice, not custom faces")
		}
	} else if r.Sides, err = p.expect(lexer.NumberToken, "die sides"); err != nil {
		return nil, err
	}
	return r, nil
}

func (p *parser) expect(t lexer.TokenType, what string) (int, error) {
	if !p.more() || p.peek().Type != t {
		return 0, p.errorf("expected %s", what)
	}
	return p.number()
}

// number consumes the current token as an integer.
func (p *parser) number() (int, error) {
	raw := p.peek().Raw
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, p.errorf("invalid number %q", raw)
	}
	p.pos++
	return v, nil
}

func (p *parser) more() bool {
	return p.pos < len(p.tokens)
}

func (p *parser) peek() lexer.Token {
	return p.tokens[p.pos]
}

func (p *parser) errorf(format string, args ...any) error {
	return fmt.Errorf("parse error at token %d: %s", p.pos, fmt.Sprintf(format, args...))
}
