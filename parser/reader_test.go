package parser

import (
	"bufio"
	"io"
	"strings"
	"testing"
)

func TestPeekRune(t *testing.T) {
	s := "Ta’ Ħaġrat"
	reader := reader{r: bufio.NewReader(strings.NewReader(s))}

	for _, r := range s {
		r1, err := reader.peekRune()
		if err != nil && !(err == io.EOF && r == 't') {
			t.Fatalf("reader.peekRune() failed with ‘%T’ at rune ‘%c’", err, r)
		}

		r2, _ := reader.readRune()
		if r1 != r2 {
			t.Fatalf("reader.peekRune() peaked ‘%c’ but the next read rune was ‘%c’", r1, r2)
		}
	}
}

func TestUnreadRune(t *testing.T) {
	s := "Ta’ Ħaġrat"
	reader := reader{r: bufio.NewReader(strings.NewReader(s))}

	for range s {
		r1, _ := reader.readRune()
		if err := reader.unreadRune(); err != nil {
			t.Fatalf("reader.unreadRune() failed with ‘%T’ when unreading rune ‘%c’", err, r1)
		}
		r2, _ := reader.readRune()
		if r1 != r2 {
			t.Fatalf("reader.readRune() returned ‘%c’ after unreading ‘%c’", r2, r1)
		}
	}
}

func TestUnreadRunePosTracking(t *testing.T) {
	s := "Ta’ Ħaġrat\nĲsselmeer"
	reader := reader{r: bufio.NewReader(strings.NewReader(s))}

	for i := range []rune(s) {
		r, _ := reader.readRune()

		if err := reader.unreadRune(); err != nil {
			t.Fatalf("reader.unreadRune() failed with ‘%T’ when unreading rune ‘%c’", err, r)
		}

		// Ensure we tracked the row correctly
		switch {
		case i < 11 && reader.pos.row != 0:
			t.Fatalf("Expected reader to be on row 0 after unreading rune ‘%c’ but got row %d", r, reader.pos.row)
		case i >= 11 && reader.pos.row != 1:
			t.Fatalf("Expected reader to be on row 1 after unreading rune ‘%c’ but got row %d", r, reader.pos.row)
		}

		// Ensure we tracked the column correctly
		switch {
		case i < 11 && reader.pos.col != uint(i):
			t.Fatalf("Expected reader to be on col %d after unreading rune ‘%c’ but got col %d", i, r, reader.pos.col)
		case i == 11 && reader.pos.col != 0:
			t.Fatalf("Expected reader to be on col 0 after unreading rune ‘\\n’ but got col %d", reader.pos.col)
		case i > 11 && reader.pos.col != uint(i)-11:
			t.Fatalf("Expected reader to be on col %d after unreading rune ‘%c’ but got col %d", i-10, r, reader.pos.col)
		}

		reader.readRune()
	}
}

func TestReadRune(t *testing.T) {
	s := "Ta’ Ħaġrat"
	reader := reader{r: bufio.NewReader(strings.NewReader(s))}

	for _, r1 := range s {
		r2, err := reader.readRune()
		if err != nil {
			t.Fatalf("reader.readRune() failed with ‘%T’", err)
		}
		if r1 != r2 {
			t.Fatalf("reader.readRune() read ‘%c’ but ‘%c’ was expected", r1, r2)
		}
	}

	_, err := reader.readRune()
	switch {
	case err == nil:
		t.Fatal("reader.readRune() expected to fail but didn’t")
	case err != io.EOF:
		t.Fatalf("reader.readRune() expected to fail with io.EOF but got ‘%T’", err)
	}
}

func TestReadRunePosTracking(t *testing.T) {
	s := "Ta’ Ħaġrat\nĲsselmeer"
	line1 := true
	reader := reader{r: bufio.NewReader(strings.NewReader(s))}

	for i := range []rune(s) {
		// Track the line we’re on
		r, _ := reader.readRune()
		if r == '\n' {
			line1 = false
		}

		// Ensure we tracked the row correctly
		switch {
		case line1 && reader.pos.row != 0:
			t.Fatalf("Expected reader to be on row 0 when reading rune ‘%c’ but got row %d", r, reader.pos.row)
		case !line1 && reader.pos.row != 1:
			t.Fatalf("Expected reader to be on row 1 when reading rune ‘%c’ but got row %d", r, reader.pos.row)
		}

		// Ensure we tracked the column correctly
		switch {
		case i < 10 && reader.pos.col != uint(i)+1:
			t.Fatalf("Expected reader to be on col %d when reading rune ‘%c’ but got col %d", i, r, reader.pos.col)
		case i == 10 && reader.pos.col != 0:
			t.Fatalf("Expected reader to be on col 0 when reading rune ‘\\n’ but got col %d", reader.pos.col)
		case i > 10 && reader.pos.col != uint(i)-10:
			t.Fatalf("Expected reader to be on col %d when reading rune ‘%c’ but got col %d", i-10, r, reader.pos.col)
		}
	}
}

func TestReadNonSpaceRune(t *testing.T) {
	s := "Ta’ Ħaġrat  \t\n Ĳsselmeer"
	reader := reader{r: bufio.NewReader(strings.NewReader(s))}
	readUntil := func(target rune) {
		for {
			if r, _ := reader.readRune(); r == target {
				break
			}
		}
	}
	assertNext := func(target rune) {
		if r, err := reader.readNonSpaceRune(); err != nil {
			t.Fatalf("reader.readNonSpaceRune() failed with ‘%T’", err)
		} else if r != target {
			t.Fatalf("reader.readNonSpaceRune() expected to return ‘%c’ but returned ‘%c’", target, r)
		}
	}

	readUntil('’')
	assertNext('Ħ')
	readUntil('t')
	assertNext('Ĳ')
}

func TestSkipSpaces(t *testing.T) {
	s := "Ta’ Ħaġrat  \t\n Ĳsselmeer"
	reader := reader{r: bufio.NewReader(strings.NewReader(s))}
	readUntil := func(target rune) {
		for {
			if r, _ := reader.readRune(); r == target {
				break
			}
		}
	}
	assertNext := func(target rune) {
		if err := reader.skipSpaces(); err != nil {
			t.Fatalf("reader.skipSpaces() failed with ‘%T’", err)
		}
		if r, _ := reader.readRune(); r != target {
			t.Fatalf("reader.readRune() expected to return ‘%c’ but returned ‘%c’", target, r)
		}
	}

	readUntil('’')
	assertNext('Ħ')
	readUntil('t')
	assertNext('Ĳ')
}
