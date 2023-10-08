package frame

import (
	"unicode"

	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-skk"
)

var defaultColor = readline.SGR3(0, 1, 39)

type _Coloring struct {
	skkbits     skk.Coloring
	bits        int
	last        rune
	defaultBits int
}

func (s *_Coloring) Init() readline.ColorSequence {
	s.bits = s.defaultBits
	s.skkbits.Init()
	return defaultColor
}

const (
	backquotedBit = 1
	percentBit    = 2
	quotedBit     = 4
	optionBit     = 8
	backSlash     = 16
)

func (s *_Coloring) Next(codepoint rune) readline.ColorSequence {
	newbits := s.bits &^ backSlash
	if codepoint == '`' {
		newbits ^= backquotedBit
	} else if codepoint == '%' {
		newbits ^= percentBit
	} else if codepoint == '"' && (s.bits&backSlash) == 0 {
		newbits ^= quotedBit
	} else if s.last == ' ' && (codepoint == '/' || codepoint == '-') {
		newbits ^= optionBit
	} else if (s.bits&optionBit) != 0 && !unicode.IsLetter(codepoint) && (codepoint < '0' || codepoint > '9') && codepoint != '-' {
		newbits &^= optionBit
	} else if s.last == '%' && (s.bits&percentBit) != 0 && unicode.IsDigit(codepoint) {
		newbits &^= percentBit
	} else if codepoint == '\\' && (s.bits&backSlash) == 0 {
		newbits |= backSlash
	}
	bits := s.bits | newbits
	color := defaultColor

	if unicode.IsControl(codepoint) {
		color = readline.SGR3(0, 1, 34) // Blue
	} else if codepoint == '\u3000' {
		color = readline.SGR2(0, 41) // RedBack
	} else if (bits & percentBit) != 0 {
		color = readline.SGR3(0, 1, 36) // Cyan
	} else if (bits & backquotedBit) != 0 {
		color = readline.SGR3(0, 1, 31) // Red
	} else if (bits & quotedBit) != 0 {
		color = readline.SGR3(0, 1, 35) // Magenta
	} else if (newbits & optionBit) != 0 {
		color = readline.SGR2(0, 33) // DarkYellow
	} else if codepoint == '&' || codepoint == '|' || codepoint == '<' || codepoint == '>' || (s.last == ' ' && codepoint == ';') {
		color = readline.SGR3(0, 1, 32) // Green
	}

	color = color.Chain(s.skkbits.Next(codepoint))

	s.bits = newbits
	s.last = codepoint
	return color
}
