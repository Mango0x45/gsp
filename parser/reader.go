package parser

import (
	"bufio"
	"errors"
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

// peakRune returns the next rune in the buffer without actually moving the
// parser forwards.
func (reader *reader) peekRune() (rune, error) {
	bytes, _ := reader.r.Peek(4)
	r, size := utf8.DecodeRune(bytes)

	switch {
	case r == utf8.RuneError && size == 0:
		return 0, io.EOF
	case r == utf8.RuneError && size == 1:
		return 0, errors.New("Tried to decode malformed UTF-8")
	}
	return r, nil
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
	r, _, err := reader.r.ReadRune()
	if r == '\n' {
		reader.pos.prevCol = reader.pos.col
		reader.pos.col = 0
		reader.pos.row++
	} else {
		reader.pos.col++
	}
	return r, err
}

func (reader *reader) readNonSpaceRune() (rune, error) {
	if err := reader.skipSpaces(); err != nil {
		return 0, err
	}

	return reader.readRune()
}

func (reader *reader) skipSpaces() error {
	for {
		if r, err := reader.readRune(); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		} else if !unicode.IsSpace(r) {
			return reader.unreadRune()
		}
	}
}
