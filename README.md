`GSP` (pronounced _gee ess pee_) is a transpiler to convert a nicer to write and
more human-friendly syntax into valid HTML.  Writing HTML can be made more
bearable using things like Emmet, but it’s still not great, and the syntax is
far too bloated, visually polluting your documents.

`GSP` will never support templating or other useless features.  If you need
support for such things, just use a programming- or macro language such as
Python or M4.

## Installation

You need to ensure you have both the Go compiler and the `make` command
available.  If you don’t have Go, you can get it [here][1].  If you’re on a
UNIX-like system such as Linux or MacOS then you should already have `make`.  If
you’re on Windows you should be using WSL anyways, and you can figure out how to
get Make.  You also need Git.

First, clone the repository and move into it:

```
$ git clone https://git.sr.ht/~mango/gsp
$ cd gsp
```

Then you can compile the transpiler with either of the two commands:

```
$ make
$ go build
```

Finally, you can install the transpiler and documentation with the following:

```
$ sudo make install
```

## Documentation

Documentation for the transpiler can be found in the `gsp(1)` manual and
documentation for the language can be found in the `gsp(5)` manual:

```
$ man gsp    # transpiler documentation
$ man 5 gsp  # language documentation
```

## Why Not Use Pug or [INSERT LANGUAGE HERE]

~Simply put, they are all trash.  Pug has decent syntax but requires you use
JavaScript.  All the others fall for the same kind of problem.  As far as I
could find, there is no pre-GSP transpiler from good syntax to HTML that works
as just one binary you call on some files.  All options force you into needing
to write JavaScript/Ruby/etc. scripts, which just isn’t good enough.~

Update 28-05-2024:

Turns out that PugCLI exists, which I actually did search for in the past but
somehow could not find.  As a result my previous statements above do not hold,
although I still find it irritating that all ‘mainstream’ besides Pug that I am
aware of (as of this update) don’t just support basic single-binary usage.

Feel free to use Pug, I simply do not care for the many Pug features I deem
useless or unergonomic (or both).

## Syntax Example

```gsp
html lang="en" {
  head {
    meta charset="UTF-8" {}
    meta name="viewport" content="width=device-width, initial-scale=1.0" {}
    link href="/favicon.svg" rel="shortcut icon" type="image/svg" {}
    link href="/style.svg" rel="stylesheet" {}
    title {-My Website Title}
  }
  body {
    p #my-id  {- This is a paragraph with the id ‘my-id’     }
    p .my-cls {- This is a paragraph with the class ‘my-cls’ }

    p
      #some-id
      .class-1
      .class-2
      key-1="value-1"
      key-2 = "value-2"
    {-
      This paragraph has an ID, two classes, and two additional attributes.  GSP
      allows us to use the ‘#ident’ and ‘.ident’ syntaxes as shorthands for
      applying IDs, and classes.  This is a text node, so nothing is being
      interpreted as GSP nodes, but we can include them inline if we want.  As
      an example, here is some @em {-emphatic} text.  Your inline nodes can also
      have attributes @em #id .cls {-just like a regular node}.
    }
  }
}
```

## Why The Name GSP?

I was originally inspired by Pug, but my dog is a GSP, not a pug.

[1]: https://go.dev/dl/
