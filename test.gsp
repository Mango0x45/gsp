html {
	head {
		title {- Hello, World!}
		/ meta foo="bar" hello world {}

		style {
			.foo {
				background-color: red;
			}
		}
		/script {
			const x = (a, b, c) => {
				{
					submit(a);
					submit(b);
					submit(c);
				}
				console.log(a, b, c);
			};
		}
		x:thing foo foo="hi" foo="bye" {}
		foo {-
			Hello, @em .red { strong {- "World" } }!
			@p {= I am not trimmed! }
		}
	}
}
