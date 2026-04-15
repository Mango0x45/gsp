package strconv

import "testing"

func TestEscapeString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Empty string",
			input: "",
			want:  "",
		},
		{
			name:  "No escapes needed",
			input: "hello world",
			want:  "hello world",
		},
		{
			name:  "Double quotes",
			input: `hello "world"`,
			want:  `hello \"world\"`,
		},
		{
			name:  "Backslashes",
			input: `C:\path\to\dir`,
			want:  `C:\\path\\to\\dir`,
		},
		{
			name:  "Mixed quotes and backslashes",
			input: `\"`,
			want:  `\\\"`,
		},
		{
			name:  "Consecutive escapes",
			input: `""\\`,
			want:  `\"\"\\\\`,
		},
		{
			name:  "Only double quote",
			input: `"`,
			want:  `\"`,
		},
		{
			name:  "Only backslash",
			input: `\`,
			want:  `\\`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EscapeString(tt.input); got != tt.want {
				t.Errorf("EscapeString(%q) = %q, want %q",
					tt.input, got, tt.want)
			}
		})
	}
}
