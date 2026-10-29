package lexer

//go:generate stringer -type=TokenType
type TokenType int

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
	ErrorToken TokenType = iota
	EofToken
	KeepHighestToken
	KeepLowestToken
	DropHighestToken
	DropLowestToken
	ExplodeToken
	DieToken
	NumberToken
	AdditionToken
	SubtractionToken
	FacesOpenToken
	FacesCloseToken
	FacesSeparatorToken
	FaceToken
	ExplodeOpenToken
	ExplodeCloseToken
	ExplodeCountToken
)

type Token struct {
	Type TokenType
	Raw  string
}
