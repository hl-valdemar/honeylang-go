package lexer

import (
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
			case '(':
				l.advance()
				tokens.append(Token{LeftParen, start, l.pos})
			case ')':
				l.advance()
				tokens.append(Token{RightParen, start, l.pos})
			case '[':
				l.advance()
				tokens.append(Token{LeftBracket, start, l.pos})
			case ']':
				l.advance()
				tokens.append(Token{RightBracket, start, l.pos})
			case '{':
				l.advance()
				tokens.append(Token{LeftCurly, start, l.pos})
			case '}':
				l.advance()
				tokens.append(Token{RightCurly, start, l.pos})
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
