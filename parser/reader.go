package parser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"
)

// position represents the current position of the reader in the buffer.  It has
// both a row and column, as well as a previous column.  We need that previous
// column to figure out our position when backtracking with reader.unreadRune().
type position struct {
	col     uint
	row     uint
	prevCol uint
}

// String returns the position of the last-processed rune in the buffer in the
// standard ‘row:col’ format, while starting at 1:1 to be vim-compatible.
func (p position) String() string {
	return fmt.Sprintf("%d:%d", p.row+1, p.col)
}

// reader represents the actual parser.  It has both an underlying buffered
// reader, and a position.
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

// unreadRune moves the parser one rune back, allowing for basic backtracking.
// You can only safely call reader.unreadRune() once before reading again.  If
// not you will risk messing up the parsers position tracking if you unread
// multiple newlines.
func (reader *reader) unreadRune() error {
	if reader.pos.col == 0 {
		reader.pos.col = reader.pos.prevCol
		reader.pos.row--
	} else {
		reader.pos.col--
	}

	return reader.r.UnreadRune()
}

// readRune reads and returns the next rune in the readers buffer.  It is just a
// wrapper around ‘reader.r.ReadRune()’ that updates the readers position.
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

// readNonSpaceRune is identical to reader.readRune(), except it skips any runes
// representing spaces (as defined by unicode).
func (reader *reader) readNonSpaceRune() (rune, error) {
	if err := reader.skipSpaces(); err != nil {
		return 0, err
	}

	return reader.readRune()
}

// skipSpaces moves the parser forwards so that the next call to
// ‘reader.readRune()’ is guaranteed to read a non-space rune as defined by
// Unicode.
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
