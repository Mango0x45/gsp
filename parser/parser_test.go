package parser

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestParseFile(t *testing.T) {
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

	// Write the source to a temp file
	r := strings.NewReader(s)
	f, _ := os.CreateTemp("", "tmp*")
	defer f.Close()
	io.Copy(f, r)
	f.Seek(0, 0)
	ast, _ := ParseFile(f)

	// This is just… easier
	s = fmt.Sprintf("%+v", ast)
	result := `{Type:1 Text: Attrs:[] Children:[{Type:0 Text:html Attrs:[{Key:lang Value:en}] Children:[{Type:0 Text:head Attrs:[{Key:attr Value:}] Children:[{Type:0 Text:title Attrs:[] Children:[{Type:1 Text: Attrs:[] Children:[{Type:3 Text:
	      My Website
	     Attrs:[] Children:[] Newline:false}] Newline:false}] Newline:true} {Type:0 Text:meta Attrs:[{Key:x Value:y}] Children:[] Newline:false}] Newline:true} {Type:0 Text:body Attrs:[] Children:[{Type:0 Text:div Attrs:[{Key:id Value:some-id}] Children:[] Newline:true} {Type:0 Text:div Attrs:[{Key:key Value:val} {Key:class Value:class-1} {Key:class Value:class-2}] Children:[{Type:0 Text:p Attrs:[] Children:[{Type:1 Text: Attrs:[] Children:[{Type:3 Text: This is some  Attrs:[] Children:[] Newline:false} {Type:0 Text:em Attrs:[] Children:[{Type:1 Text: Attrs:[] Children:[{Type:3 Text:emphatic Attrs:[] Children:[] Newline:false}] Newline:false}] Newline:false} {Type:3 Text: text	   Attrs:[] Children:[] Newline:false}] Newline:false}] Newline:false}] Newline:false} {Type:0 Text:tags Attrs:[{Key:key Value:Some long value}] Children:[] Newline:false}] Newline:true}] Newline:false}] Newline:false}`
	if s != result {
		t.Fatalf("ParseFile() parsed unexpected AST ‘%s’", s)
	}
}

func TestValidNameStartChar(t *testing.T) {
	for _, r := range "HELLO_th:re_wörld" {
		if !validNameStartChar(r) {
			t.Fatalf("validNameStartChar() returned false on valid rune ‘%c’", r)
		}
	}
}

func TestValidNameChar(t *testing.T) {
	for _, r := range "hello69-th.re-wörld" {
		if !validNameChar(r) {
			t.Fatalf("validNameChar() returned false on valid rune ‘%c’", r)
		}
	}
}
