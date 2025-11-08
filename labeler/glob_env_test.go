package labeler

import (
	"testing"
)

func TestMatchGlobWithDisabledExtglob(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		filename string
		disabled bool
		want     bool
	}{
		{
			name:     "Extglob pattern - enabled",
			pattern:  "*.@(js|ts)",
			filename: "file.js",
			disabled: false,
			want:     true,
		},
		{
			name:     "Extglob pattern - disabled (treated as literal)",
			pattern:  "*.@(js|ts)",
			filename: "file.js",
			disabled: true,
			want:     false, // Treated as literal pattern, doesn't match
		},
		{
			name:     "Extglob pattern - disabled (literal match)",
			pattern:  "*.@(js|ts)",
			filename: "file.@(js|ts)",
			disabled: true,
			want:     true, // Literal match
		},
		{
			name:     "Regular glob pattern - enabled",
			pattern:  "*.js",
			filename: "file.js",
			disabled: false,
			want:     true,
		},
		{
			name:     "Regular glob pattern - disabled",
			pattern:  "*.js",
			filename: "file.js",
			disabled: true,
			want:     true, // Regular glob still works
		},
		{
			name:     "Negation extglob - enabled",
			pattern:  "!(test).js",
			filename: "main.js",
			disabled: false,
			want:     true,
		},
		{
			name:     "Negation extglob - disabled",
			pattern:  "!(test).js",
			filename: "main.js",
			disabled: true,
			want:     true, // ! is treated as negation, !(test).js doesn't match, so result is inverted to true
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.disabled {
				t.Setenv("GH_LABEL_KIT_LABELER_DISABLE_EXTGLOB", "1")
			}

			// Test
			got := matchGlob(tt.pattern, tt.filename)
			if got != tt.want {
				t.Errorf("matchGlob(%q, %q) with disabled=%v = %v, want %v",
					tt.pattern, tt.filename, tt.disabled, got, tt.want)
			}
		})
	}
}
