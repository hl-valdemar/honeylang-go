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

		// skip whitespace (except for newlines)
		if unicode.IsSpace(*r) && *r != '\n' {
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
			case '\n':
				l.advance()
				tokens.append(TokenDesc{NewLine, start, l.pos})
			case ',':
				l.advance()
				tokens.append(TokenDesc{Comma, start, l.pos})
			case '(':
				l.advance()
				tokens.append(TokenDesc{LeftParen, start, l.pos})
			case ')':
				l.advance()
				tokens.append(TokenDesc{RightParen, start, l.pos})
			case '[':
				l.advance()
				tokens.append(TokenDesc{LeftBracket, start, l.pos})
			case ']':
				l.advance()
				tokens.append(TokenDesc{RightBracket, start, l.pos})
			case '{':
				l.advance()
				tokens.append(TokenDesc{LeftCurly, start, l.pos})
			case '}':
				l.advance()
				tokens.append(TokenDesc{RightCurly, start, l.pos})
			case '=':
				l.advance()
				tokens.append(TokenDesc{Equal, start, l.pos})
			case ':':
				l.advance()
				next := l.peek()
				if next != nil && *next == ':' {
					l.advance()
					tokens.append(TokenDesc{DoubleColon, start, l.pos})
				} else {
					tokens.append(TokenDesc{Colon, start, l.pos})
				}
			default:
				// unknown character encountered, report error!
				l.advance()
				l.errors.append(ScanErrorDesc{UnrecognizedCharacter, start, start + 1})
			}
		}
	}

	tokens.append(TokenDesc{EOF, l.pos, l.pos})
	return tokens, l.errors
}

func (l *Lexer) scanIdent() TokenDesc {
	start := l.pos

	for {
		r := l.peek()
		if r == nil {
			break
		}

		if !unicode.IsLetter(*r) && !unicode.IsDigit(*r) && *r != '_' {
			break
		}

		l.advance()
	}

	// handle potential keywords
	ident := l.src.Contents[start:l.pos]
	if kind, ok := identKeyword[string(ident)]; ok {
		return TokenDesc{kind, start, l.pos}
	}

	return TokenDesc{Identifier, start, l.pos}
}

func (l *Lexer) scanNum() TokenDesc {
	start := l.pos

	// check for hex (0x/0X) or binary (0b/0B) prefix
	r := l.peek()
	if r != nil && *r == '0' {
		next := l.peekOffset(1)
		if next != nil && (*next == 'x' || *next == 'X') {
			l.consumeHex()
			return TokenDesc{Number, start, l.pos}
		} else if next != nil && (*next == 'b' || *next == 'B') {
			l.consumeBin()
			return TokenDesc{Number, start, l.pos}
		}
	}

	// check for decimal
	l.consumeDec()
	return TokenDesc{Number, start, l.pos}
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
		l.errors.append(ScanErrorDesc{EmptyHexLiteral, start, l.pos})
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
		l.errors.append(ScanErrorDesc{EmptyBinaryLiteral, start, l.pos})
	}
}

func (l *Lexer) consumeDec() {
	has_decimal := false
	has_error := false

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
			if next != nil && unicode.IsDigit(*next) {
				has_decimal = true
				l.advance()
			} else {
				break
			}
		} else if *c == '.' && has_decimal {
			// multiple decimal points
			if !has_error {
				l.errors.append(ScanErrorDesc{MultipleDecimalPoints, l.pos, l.pos + 1})
				has_error = true
			}
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
	l.pos++
}

func (l *Lexer) advanceBy(n uint) {
	l.pos += n
}
