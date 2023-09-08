package parser

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"
)

type position struct {
	col     uint
	row     uint
	prevCol uint
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d", p.row+1, p.col)
}

type reader struct {
	r   *bufio.Reader
	pos position
}

func (reader *reader) peekRune() (rune, error) {
	bytes := make([]byte, 0, 4)
	var err error

	// Peeking the next rune is annoying.  We want to get the next rune
	// which could be the next 1–4 bytes.  Normally we can just call
	// reader.r.Peek(4) but that doesn’t work here as the last rune in a
	// file could be a 1–3 byte rune, so we would fail with an EOF error.
	for i := 4; i > 0; i-- {
		if bytes, err = reader.r.Peek(i); err == io.EOF {
			continue
		} else if err != nil {
			return 0, err
		} else {
			rune, _ := utf8.DecodeRune(bytes)
			return rune, nil
		}
	}

	return 0, io.EOF
}

func (reader *reader) unreadRune() error {
	if reader.pos.col == 0 {
		reader.pos.col = reader.pos.prevCol
		reader.pos.row--
	} else {
		reader.pos.col--
	}

	return reader.r.UnreadRune()
}

func (reader *reader) readRune() (rune, error) {
	rune, _, err := reader.r.ReadRune()
	if rune == '\n' {
		reader.pos.prevCol = reader.pos.col
		reader.pos.col = 0
		reader.pos.row++
	} else {
		reader.pos.col++
	}
	return rune, err
}

func (reader *reader) readNonSpaceRune() (rune, error) {
	if err := reader.skipSpaces(); err != nil {
		return 0, err
	}

	if r, err := reader.readRune(); err != nil {
		return 0, err
	} else {
		return r, nil
	}
}

func (reader *reader) skipSpaces() error {
	for {
		if rune, err := reader.readRune(); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		} else if !unicode.IsSpace(rune) {
			return reader.unreadRune()
		}
	}
}
