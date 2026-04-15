GSP (pronounced _gee ess pee_) is an alternative syntax for HTML.
Writing HTML can be made more bearable using things like Emmet, but
it’s still not great, and the syntax is far too bloated, visually
polluting your documents.  It is for this reason that GSP exists.

GSP supports templating via macros that execute shell scripts,
allowing GSP to integrate nicely into the UNIX environment.

## Source Installation

Installation depends on the Go compiler and `make`.

First, clone the repository or fetch a release tarball and move into
it:

```
$ git clone https://git.thomasvoss.com/gsp
$ cd gsp
```

Then you can compile the transpiler:

```
$ make
```

Finally, you can install the transpiler and documentation with the
following:

```
$ sudo make install
```

## Documentation

Documentation for the transpiler can be found in the `gsp(1)` manual
and documentation for the language can be found in the `gsp(5)`
manual:

```
$ man gsp    # transpiler documentation
$ man 5 gsp  # language documentation
```

## Syntax Example

```gsp
macro today {
	date '+%A, %d %B %Y'
}

html lang="en" {
	head {
		meta charset="UTF-8" {}
		meta name="viewport" content="width=device-width, initial-scale=1.0" {}
		link href="/favicon.svg" rel="shortcut icon" type="image/svg" {}
		link href="/style.svg" rel="stylesheet" {}
		title {-My Website Title}
		/ style {
			This entire style node is commented out.
		}
		script {
			const myJSVar = 'Hello, World!';
			const x = () => {
				return 'Escaping of characters in JavaScript/CSS'
				+ 'blocks not required!';
			};
		}
	}
	body {
		p #my-id  {- This is a paragraph with the id ‘my-id’     }
		p .my-cls {- This is a paragraph with the class ‘my-cls’ }

		p
			#some-id
			.class-1
			.class-2
			key-1="value-1"
			key-2="value-2"
		{-
			This paragraph has an ID, two classes, and two additional
			attributes.  GSP allows us to use the ‘#ident’ and
			‘.ident’ syntaxes as shorthands for applying IDs, and
			classes.  This is a text node, so nothing is being
			interpreted as GSP nodes, but we can include them inline
			if we want.  As an example, here is some @em {-emphatic}
			text.  Your inline nodes can also have attributes
			@em #id .cls {-just like a regular node}.
		}
	}
	footer {-
		Written on @$today{}
	}
}
```

## Why The Name GSP?

I was originally inspired by Pug, but my dog is a GSP, not a pug.