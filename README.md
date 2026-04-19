GSP (pronounced _gee ess pee_) is an HTML-compatible markup language.
The standard GSP distribution includes a GSP to HTML transpiler as
well as utility programs.

Writing HTML can be made more bearable using things like Emmet, but
it’s still not great, and the syntax is far too bloated, visually
polluting your documents.  It is for this reason that GSP exists.

On top of standard markup, GSP has native support for templating via
external processes.  This allows the macro implementation to remain
extremely small and simple, while retaining a high degree of power.

## Source Installation

Installation depends on the Go compiler and `make`.

First, clone the repository or fetch a release tarball and move into
it:

```
$ git clone https://git.thomasvoss.com/gsp
$ cd gsp
```

Then you can compile the binaries:

```
$ make
```

Finally, you can install the binaries and documentation with the
following:

```
# make install
```

## Documentation

Various manuals ship with the standard GSP distribution:

```
$ man 1 gsp                     # transpiler documentation
$ man 1 gspesc                  # input escaping documentation
$ man 5 gsp                     # language documentation
$ man 7 gsp-macros              # macro system documentation
```

The `example.gsp` example document is also provided at the root of the
repository, and is installed along with the other manual pages
(typically at `/usr/share/gsp/doc`, but the exact location will vary
per system).

## Syntax Example

```gsp
html lang="en" {
	head {
		meta charset="UTF-8" {}
		meta
			name="viewport"
			content="width=device-width, initial-scale=1.0"
		{}
		link href="/style.css" rel="stylesheet" {}
		title {- GSP Language Reference}
	}
	body {
		p #first-p .red {-
			GSP allows us to define IDs and classes on
			tags in a manner that matches CSS selector
			syntax.
		}

		/ div {
			p {-
				It also supports comments that operate
				on a syntactical level.
			}
			p {- Isn’t that neat? }
		}

		p {-
			GSP also features a powerful macro system,
			allowing us to perform tasks like syntax
			highlighting source code from within a
			document.
		}

		$$syntax_highlight lang="c" {-
#include <stdio.h>
#include <stdlib.h>

int
main(int argc, char **argv)
{
	puts("Hello, World!");
	return EXIT_SUCCESS;
}
		}
	}
}
```

## Why The Name GSP?

I was originally inspired by Pug, but my dog is a GSP, not a pug.