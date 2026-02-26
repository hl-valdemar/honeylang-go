package lexer

import "fmt"

type TokenKind uint

const (
	// complex tokens
	Identifier TokenKind = iota
	Number

	// single-char tokens
	Colon
	Comma
	LeftParen
	RightParen
	LeftBracket
	RightBracket
	LeftCurly
	RightCurly

	// double-char tokens
	DoubleColon

	// special
	EOF
)

var tokenKindNames = [...]string{
	Identifier:   "Identifier",
	Number:       "Number",
	Colon:        "Colon",
	Comma:        "Comma",
	LeftParen:    "LeftParen",
	RightParen:   "RightParen",
	LeftBracket:  "LeftBracket",
	RightBracket: "RightBracket",
	LeftCurly:    "LeftCurly",
	RightCurly:   "RightCurly",
	DoubleColon:  "DoubleColon",
	EOF:          "EOF",
}

func (k TokenKind) String() string {
	if int(k) < len(tokenKindNames) {
		return tokenKindNames[k]
	}
	return "Unknown"
}

type Token struct {
	kind       TokenKind
	start, end uint
}

type Tokens struct {
	Kinds        []TokenKind
	Starts, Ends []uint
}

func initTokens() Tokens {
	tokInitCap := 100
	return Tokens{
		make([]TokenKind, 0, tokInitCap),
		make([]uint, 0, tokInitCap),
		make([]uint, 0, tokInitCap),
	}
}

func (t *Tokens) append(tok Token) {
	t.assertHealth()
	t.Kinds = append(t.Kinds, tok.kind)
	t.Starts = append(t.Starts, tok.start)
	t.Ends = append(t.Ends, tok.end)
}

func (t *Tokens) Len() int {
	t.assertHealth()
	return len(t.Kinds)
}

func (t *Tokens) assertHealth() {
	if !(len(t.Kinds) == len(t.Starts) && len(t.Starts) == len(t.Ends)) {
		panic(fmt.Sprintf("Parallel arrays out of sync! [%T]", *t))
	}
}
