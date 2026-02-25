package lexer

import (
	"fmt"
	"honey/source"
	"unicode"
)

type Lexer struct {
	src    *source.Source
	pos    uint
	errors ScanErrors
}

func Init(src *source.Source) Lexer {
	return Lexer{src, 0, initErrors()}
}

func Scan(src *source.Source) (Tokens, ScanErrors) {
	lexer := Init(src)
	return lexer.Scan()
}

func (l *Lexer) Scan() (Tokens, ScanErrors) {
	tokens := initTokens()

	for {
		r := l.peek()
		if r == nil {
			break
		}

		// skip whitespace
		if unicode.IsSpace(*r) {
			l.advance()
			continue
		}

		// skip comments
		if *r == '#' {
			l.advance()
			for {
				next := l.peek()
				if next == nil || *next == '\n' {
					break
				}
				l.advance()
			}
			continue
		}

		start := l.pos
		if unicode.IsLetter(*r) || *r == '_' {
			t := l.scanIdent()
			tokens.append(t)
		} else if unicode.IsDigit(*r) {
			t := l.scanNum()
			tokens.append(t)
		} else {
			switch *r {
			case ',':
				l.advance()
				tokens.append(Token{Comma, start, l.pos})
			case ':':
				l.advance()
				next := l.peek()
				if next != nil && *next == ':' {
					l.advance()
					tokens.append(Token{DoubleColon, start, l.pos})
				} else {
					tokens.append(Token{Colon, start, l.pos})
				}
			default:
				// unknown character encountered, report error!
				l.advance()
				l.errors.append(ScanError{UnrecognizedCharacter, start, start + 1})
			}
		}
	}

	tokens.append(Token{EOF, l.pos, l.pos})
	return tokens, l.errors
}

func (l *Lexer) scanIdent() Token {
	start := l.pos

	for {
		r := l.peek()
		if r == nil {
			break
		}

		if !unicode.IsLetter(*r) && *r != '_' {
			break
		}

		l.advance()
	}

	// TODO: handle potential keywords

	return Token{Identifier, start, l.pos}
}

func (l *Lexer) scanNum() Token {
	start := l.pos
	has_decimal := false
	has_error := false

	// check for hex (0x/0X) or binary (0b/0B) prefix
	r := l.peek()
	if r != nil && *r == '0' {
		next := l.peekOffset(1)
		if next != nil && (*next == 'x' || *next == 'X') {
			l.consumeHex()
			return Token{Number, start, l.pos}
		} else if next != nil && (*next == 'b' || *next == 'B') {
			l.consumeBin()
			return Token{Number, start, l.pos}
		}
	}

	// check for decimal
	for {
		c := l.peek()
		if c == nil {
			break
		}

		if unicode.IsDigit(*c) {
			l.advance()
		} else if *c == '.' && !has_decimal {
			// check if next char is also a digit
			next := l.peekOffset(1)
			if next != nil {
				if unicode.IsDigit(*next) {
					has_decimal = true
					l.advance()
				} else {
					break
				}
			} else {
				has_decimal = true
				l.advance()
			}
		} else if *c == '.' && has_decimal {
			// multiple decimal points
			if !has_error {
				l.errors.append(ScanError{MultipleDecimalPoints, l.pos, l.pos + 1})
				has_error = true
			}
			l.advance()
		} else {
			break
		}
	}

	return Token{Number, start, l.pos}
}

func (l *Lexer) consumeHex() {
	start := l.pos
	// consume '0' and 'x'/'X'
	l.advanceBy(2)

	digitStart := l.pos
	for {
		h := l.peek()
		if h == nil {
			break
		}

		if unicode.IsDigit(*h) || (*h >= 'a' && *h <= 'f') || (*h >= 'A' && *h <= 'F') {
			l.advance()
		} else {
			break
		}
	}

	if l.pos == digitStart {
		l.errors.append(ScanError{EmptyHexLiteral, start, l.pos})
	}
}

func (l *Lexer) consumeBin() {
	start := l.pos
	// consume '0' and 'b'/'B'
	l.advanceBy(2)

	digitStart := l.pos
	for {
		h := l.peek()
		if h == nil {
			break
		}

		if *h == '0' || *h == '1' {
			l.advance()
		} else {
			break
		}
	}

	if l.pos == digitStart {
		l.errors.append(ScanError{EmptyBinaryLiteral, start, l.pos})
	}
}

func (l *Lexer) peek() *rune {
	if l.pos < uint(len(l.src.Contents)) {
		return &l.src.Contents[l.pos]
	}
	return nil
}

func (l *Lexer) peekOffset(offset uint) *rune {
	if l.pos+offset < uint(len(l.src.Contents)) {
		return &l.src.Contents[l.pos+offset]
	}
	return nil
}

func (l *Lexer) advance() {
	l.pos++
}

func (l *Lexer) advanceBy(n uint) {
	l.pos += n
}

// TOKENS

type TokenKind uint

const (
	// complex tokens
	Identifier TokenKind = iota
	Number

	// single-char tokens
	Colon
	Comma

	// double-char tokens
	DoubleColon

	// special
	EOF
)

var tokenKindNames = [...]string{
	Identifier:  "Identifier",
	Number:      "Number",
	Colon:       "Colon",
	Comma:       "Comma",
	DoubleColon: "DoubleColon",
	EOF:         "EOF",
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

// ERRORS

type ScanErrorKind uint

const (
	UnrecognizedCharacter ScanErrorKind = iota
	MultipleDecimalPoints
	EmptyHexLiteral
	EmptyBinaryLiteral
)

var scanErrorKindNames = [...]string{
	UnrecognizedCharacter: "unrecognized character",
	MultipleDecimalPoints: "multiple decimal points in number literal",
	EmptyHexLiteral:       "empty hex literal (expected digits after '0x')",
	EmptyBinaryLiteral:    "empty binary literal (expected digits after '0b')",
}

func (k ScanErrorKind) String() string {
	if int(k) < len(scanErrorKindNames) {
		return scanErrorKindNames[k]
	}
	return "unknown error"
}

type ScanError struct {
	kind       ScanErrorKind
	start, end uint
}

type ScanErrors struct {
	Kinds        []ScanErrorKind
	Starts, Ends []uint
}

func initErrors() ScanErrors {
	errInitCap := 10
	return ScanErrors{
		make([]ScanErrorKind, 0, errInitCap),
		make([]uint, 0, errInitCap),
		make([]uint, 0, errInitCap),
	}
}

func (e *ScanErrors) append(err ScanError) {
	e.assertHealth()
	e.Kinds = append(e.Kinds, err.kind)
	e.Starts = append(e.Starts, err.start)
	e.Ends = append(e.Ends, err.end)
}

func (e *ScanErrors) Len() int {
	e.assertHealth()
	return len(e.Kinds)
}

func (e *ScanErrors) assertHealth() {
	if !(len(e.Kinds) == len(e.Starts) && len(e.Starts) == len(e.Ends)) {
		panic(fmt.Sprintf("Parallel arrays out of sync! [%T]", *e))
	}
}
