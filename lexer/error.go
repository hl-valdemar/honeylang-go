package lexer

import "fmt"

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

type ScanErrorDesc struct {
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

func (e *ScanErrors) append(err ScanErrorDesc) {
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
