package lexer

import (
	"fmt"
	"honey/source"
	"unicode"
)

type Lexer struct {
	src *source.Source
	pos uint
}

func Init(src *source.Source) Lexer {
	return Lexer{src, 0}
}

func Scan(src *source.Source) (Tokens, ScanErrors) {
	lexer := Init(src)
	return lexer.Scan()
}

func (l *Lexer) Scan() (Tokens, ScanErrors) {
	tokens := initTokens()
	errors := initErrors()

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
			tokens.append(t.kind, t.start, t.end)
		} else if unicode.IsDigit(*r) {
			t := l.scanNum()
			tokens.append(t.kind, t.start, t.end)
		} else {
			switch *r {
			case ',':
				l.advance()
				tokens.append(Comma, start, l.pos)
			case ':':
				l.advance()
				next := l.peek()
				if next != nil && *next == ':' {
					l.advance()
					tokens.append(DoubleColon, start, l.pos)
				} else {
					tokens.append(Colon, start, l.pos)
				}
			default:
				// unknown character encountered, report error!
				l.advance()
				errors.append(UnrecognizedCharacter, start, start+1)
			}
		}
	}

	tokens.append(EOF, l.pos, l.pos)
	return tokens, errors
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

	// TODO: handle potential keyword

	return Token{Identifier, start, l.pos}
}

func (l *Lexer) scanNum() (*Token, *ScanError) {
	start := l.pos
	has_decimal := false
	has_error := false

	// check for hex (0x/0X) or binary (0b/0B) prefix
	r := l.peek()
	if r != nil && *r == '0' {
		next := l.peekOffset(1)
		if next != nil && (*next == 'x' || *next == 'X') {
			l.consumeHex()
			return &Token{Number, start, l.pos}, nil
		} else if next != nil && (*next == 'b' || *next == 'B') {
			l.consumeBin()
			return &Token{Number, start, l.pos}, nil
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
				// TODO: report error
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
	// consume '0' and 'x'/'X'
	l.advanceBy(2)

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
}

func (l *Lexer) consumeBin() {
	// consume '0' and 'b'/'B'
	l.advanceBy(2)

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
	l.pos += 1
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

func (e *Tokens) append(kind TokenKind, start, end uint) {
	e.assertHealth()
	e.Kinds = append(e.Kinds, kind)
	e.Starts = append(e.Starts, start)
	e.Ends = append(e.Ends, end)
}

func (e *Tokens) Len() int {
	e.assertHealth()
	return len(e.Kinds)
}

func (e *Tokens) assertHealth() {
	if !(len(e.Kinds) == len(e.Starts) && len(e.Starts) == len(e.Ends)) {
		panic(fmt.Sprintf("Parallel arrays out of sync! [%T]", *e))
	}
}

// ERRORS

type ScanErrorKind = uint

const (
	UnrecognizedCharacter ScanErrorKind = iota
)

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

func (e *ScanErrors) append(kind ScanErrorKind, start, end uint) {
	e.assertHealth()
	e.Kinds = append(e.Kinds, kind)
	e.Starts = append(e.Starts, start)
	e.Ends = append(e.Ends, end)
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
