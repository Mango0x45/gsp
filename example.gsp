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
.Ed
