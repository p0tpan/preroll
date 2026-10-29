package lexer

//go:generate stringer -type=tokenType
type tokenType int

const (
	EOF             = '\u0000'
	DIGITS          = "0123456789"
	NON_ZERO_DIGITS = "123456789"
	ASCII_ALPHA     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

const (
	keepHighest    = "kh"
	keepLowest     = "kl"
	dropHighest    = "dh"
	dropLowest     = "dl"
	explode        = "!"
	die            = "d"
	facesOpen      = "{"
	facesClose     = "}"
	explodeOpen    = "{"
	explodeClose   = "}"
	addition       = "+"
	subtraction    = "-"
	facesSeparator = ","
)

const (
	errorToken tokenType = iota
	eofToken
	keepHighestToken
	keepLowestToken
	dropHighestToken
	dropLowestToken
	explodeToken
	dieToken
	numberToken
	additionToken
	subtractionToken
	facesOpenToken
	facesCloseToken
	facesSeparatorToken
	faceToken
	explodeOpenToken
	explodeCloseToken
	explodeCountToken
)

type Token struct {
	Type tokenType
	Raw  string
}
