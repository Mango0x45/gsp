package formatter

import (
	"io"
	"os"
	"strings"
	"sync"
	"testing"

	"git.thomasvoss.com/gsp/parser"
)

var (
	stdoutMutex sync.Mutex
	r           *os.File
	w           *os.File
	stdout      *os.File
)

func redirectStdout() {
	stdoutMutex.Lock()
	stdout = os.Stdout
	r, w, _ = os.Pipe()
	os.Stdout = w
}

func restoreAndCapture() string {
	defer stdoutMutex.Unlock()
	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = stdout
	return string(out)
}

func TestPrintAst(t *testing.T) {
	s := `
	html lang="en" {
	  >head attr {
	    >title {-
	      My Website
	    }
	    meta x="y"{}
	  }
	  >body {
	    >div #some-id{}
	    div key="val" .class-1 .class-2 {
	      p {- This is some @em{-emphatic} text	  }
	    }

	    tags key  = "Some long value" {}
	  }
	}`
	result := `<html lang="en"><head attr><title>
	      My Website
	    </title>
<meta x="y"></head>
<body><div id="some-id">
<div class="class-1 class-2" key="val"><p> This is some <em>emphatic</em> text	  </p></div><tags key="Some long value"></body>
</html>`

	// Write the source to a temp file
	r := strings.NewReader(s)
	f, _ := os.CreateTemp("", "tmp*")
	defer f.Close()
	io.Copy(f, r)
	f.Seek(0, 0)
	ast, _ := parser.ParseFile(f)

	redirectStdout()
	PrintAst(ast)

	out := restoreAndCapture()
	if out != result {
		t.Fatalf("PrintAst() printed unexpected string ‘%s’", out)
	}
}

func TestPrintAttrs_nClasses(t *testing.T) {
	redirectStdout()
	printAttrs([]parser.Attr{
		{
			Key:   "foo",
			Value: "b&r",
		},
		{
			Key: "baz",
		},
		{
			Key:   "hello",
			Value: "<world\"",
		},
	})

	out := restoreAndCapture()
	if out != " foo=\"b&amp;r\" baz hello=\"&lt;world&quot;\"" {
		t.Fatalf("printAttrs() printed unexpected string ‘%s’", out)
	}
}

func TestPrintAttrs_classes(t *testing.T) {
	redirectStdout()
	printAttrs([]parser.Attr{
		{
			Key:   "class",
			Value: "\"foo\"",
		},
		{
			Key:   "class",
			Value: "<bar>",
		},
		{
			Key:   "class",
			Value: "b&z",
		},
	})

	out := restoreAndCapture()
	if out != " class=\"&quot;foo&quot; &lt;bar> b&amp;z\"" {
		t.Fatalf("printAttrs() printed unexpected string ‘%s’", out)
	}
}

func TestPrintAttrs_mixedAttrs(t *testing.T) {
	redirectStdout()
	printAttrs([]parser.Attr{
		{
			Key:   "foo",
			Value: "bar",
		},
		{
			Key: "baz",
		},
		{
			Key:   "class",
			Value: "foo",
		},
		{
			Key:   "hello",
			Value: "world",
		},
		{
			Key:   "class",
			Value: "bar",
		},
		{
			Key:   "class",
			Value: "baz",
		},
	})

	out := restoreAndCapture()
	if out != " class=\"foo bar baz\" foo=\"bar\" baz hello=\"world\"" {
		t.Fatalf("printAttrs() printed unexpected string ‘%s’", out)
	}
}

func TestPrintText(t *testing.T) {
	redirectStdout()
	printText("'Hello' <em>there</em> to you & \"world\"!")

	out := restoreAndCapture()
	if out != "&apos;Hello&apos; &lt;em&gt;there&lt;/em&gt; to you &amp; &quot;world&quot;!" {
		t.Fatalf("printText() printed unexpected string ‘%s’", out)
	}
}

func TestPrintChildrenTextOnly(t *testing.T) {
	redirectStdout()
	printChildren([]parser.AstNode{
		{
			Type: parser.Text,
			Text: "   \t\x0A Hello ",
		},
		{
			Type: parser.Text,
			Text: " There ",
		},
		{
			Type: parser.Text,
			Text: " World \x0A\t   ",
		},
	})

	out := restoreAndCapture()
	if out != "   \t\x0A Hello  There  World \x0A\t   " {
		t.Fatalf("printChildren() printed unexpected string ‘%s’", out)
	}
}

func TestPrintChildrenTrimMixed(t *testing.T) {
	redirectStdout()
	children := []parser.AstNode{
		{
			Type: parser.Tagless,
			Children: []parser.AstNode{
				{
					Type: parser.Text,
					Text: " Hello World",
				},
			},
		},
	}
	printChildrenTrim([]parser.AstNode{
		{
			Type: parser.Normal,
			Text: "em",
			Children: children,
		},
		{
			Type: parser.Text,
			Text: "Foo Bar  ",
		},
	})

	out := restoreAndCapture()
	if out != "<em> Hello World</em>Foo Bar" {
		t.Fatalf("printChildrenTrim() printed unexpected string ‘%s’", out)
	}
}

func TestPrintChildrenTrimTextOnly(t *testing.T) {
	redirectStdout()
	printChildrenTrim([]parser.AstNode{
		{
			Type: parser.Text,
			Text: "   \t\x0A Hello ",
		},
		{
			Type: parser.Text,
			Text: " There ",
		},
		{
			Type: parser.Text,
			Text: " World \x0A\t   ",
		},
	})

	out := restoreAndCapture()
	if out != "Hello  There  World" {
		t.Fatalf("printChildrenTrim() printed unexpected string ‘%s’", out)
	}
}

func TestTrimLeftSpaces(t *testing.T) {
	sx := "   \t \x0AHello World\x0A \t   "
	sy := "Hello World\x0A \t   "
	if sz := trimLeftSpaces(sx); sz != sy {
		t.Fatalf("trimLeftSpaces() returned ‘%s’", sy)
	}
}

func TestTrimRightSpaces(t *testing.T) {
	sx := "   \t \x0AHello World\x0A \t   "
	sy := "   \t \x0AHello World"
	if sz := trimRightSpaces(sx); sz != sy {
		t.Fatalf("trimRightSpaces() returned ‘%s’", sy)
	}
}
