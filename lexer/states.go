package lexer

import (
	"strings"
)

func isDecimalDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isASCIIAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

type stateFunc func(*lexer) stateFunc

// skipSpace discards spaces and tabs between tokens.
func skipSpace(l *lexer) {
	l.acceptRun(" \t")
	l.ignore()
}

func lexModifier(l *lexer) stateFunc {
	logger := l.logger.WithGroup("lexModifier")
	logger.Debug("entering lexModifier")
	skipSpace(l)
	switch {
	case strings.HasPrefix(l.input[l.pos:], keepHighest):
		l.pos += len(keepHighest)
		l.emit(KeepHighestToken)
		return lexDieCount
	case strings.HasPrefix(l.input[l.pos:], keepLowest):
		l.pos += len(keepLowest)
		l.emit(KeepLowestToken)
		return lexDieCount
	case strings.HasPrefix(l.input[l.pos:], dropHighest):
		l.pos += len(dropHighest)
		l.emit(DropHighestToken)
		return lexDieCount
	case strings.HasPrefix(l.input[l.pos:], dropLowest):
		l.pos += len(dropLowest)
		l.emit(DropLowestToken)
		return lexDieCount
	case strings.HasPrefix(l.input[l.pos:], explode):
		l.pos += len(explode)
		l.emit(ExplodeToken)
		return lexExplode
	default:
		logger.Debug("no modifier found")
		return lexDieCount
	}
}

func lexExplode(l *lexer) stateFunc {
	if !strings.HasPrefix(l.input[l.pos:], explodeOpen) {
		return lexDieCount
	}
	l.pos += len(explodeOpen)
	l.emit(ExplodeOpenToken)
	if !l.accept(NON_ZERO_DIGITS) {
		return l.errorf("expected explode count at position %d", l.pos)
	}
	l.acceptRun(DIGITS)
	l.emit(ExplodeCountToken)
	if !l.accept(explodeClose) {
		return l.errorf("expected %q at position %d", explodeClose, l.pos)
	}
	l.emit(ExplodeCloseToken)
	return lexDieCount
}

func lexDieCount(l *lexer) stateFunc {
	logger := l.logger.WithGroup("lexDieCount")
	logger.With("lexer", l).Debug("entering lexDieCount")
	skipSpace(l)
	if !l.accept(NON_ZERO_DIGITS) {
		return l.errorf("expected die count at position %d", l.pos)
	}
	l.acceptRun(DIGITS)
	l.emit(NumberToken)
	return lexDie
}

func lexDie(l *lexer) stateFunc {
	if !strings.HasPrefix(l.input[l.pos:], die) {
		return l.errorf("expected %q at position %d", die, l.pos)
	}
	l.pos += len(die)
	l.emit(DieToken)
	return lexFaces
}

func lexFaces(l *lexer) stateFunc {
	switch {
	case strings.HasPrefix(l.input[l.pos:], facesOpen):
		l.pos += len(facesOpen)
		l.emit(FacesOpenToken)
		for {
			l.acceptRun(DIGITS + ASCII_ALPHA)
			if l.pos == l.start {
				return l.errorf("expected face at position %d", l.pos)
			}
			l.emit(FaceToken)
			switch {
			case l.accept(facesSeparator):
				l.emit(FacesSeparatorToken)
			case l.accept(facesClose):
				l.emit(FacesCloseToken)
				return lexAddSubtract
			default:
				return l.errorf("expected %q or %q at position %d", facesSeparator, facesClose, l.pos)
			}
		}
	case l.accept(NON_ZERO_DIGITS):
		l.acceptRun(DIGITS)
		l.emit(NumberToken)
		return lexAddSubtract
	}
	return l.errorf("expected face count or %q at position %d", facesOpen, l.pos)
}

func lexAddSubtract(l *lexer) stateFunc {
	skipSpace(l)
	switch {
	case l.peek() == EOF:
		return nil
	case strings.HasPrefix(l.input[l.pos:], addition):
		l.pos += len(addition)
		l.emit(AdditionToken)
	case strings.HasPrefix(l.input[l.pos:], subtraction):
		l.pos += len(subtraction)
		l.emit(SubtractionToken)
	default:
		return l.errorf("unexpected input %q at position %d", l.input[l.pos:], l.pos)
	}
	return lexOperand
}

// lexOperand lexes the right-hand side of an operator: either a plain number
// or another roll, which may carry its own modifier.
func lexOperand(l *lexer) stateFunc {
	skipSpace(l)
	for _, m := range []string{keepHighest, keepLowest, dropHighest, dropLowest, explode} {
		if strings.HasPrefix(l.input[l.pos:], m) {
			return lexModifier
		}
	}
	if !l.accept(NON_ZERO_DIGITS) {
		return l.errorf("expected number or roll at position %d", l.pos)
	}
	l.acceptRun(DIGITS)
	l.emit(NumberToken)
	if strings.HasPrefix(l.input[l.pos:], die) {
		return lexDie
	}
	return lexAddSubtract
}
