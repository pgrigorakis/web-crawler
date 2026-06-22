package main

import "testing"

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name          string
		inputURL      string
		expected      string
		wantErr       bool
		errorContains string
	}{
		{
			name:     "remove scheme",
			inputURL: "https://www.boot.dev/blog/path",
			expected: "www.boot.dev/blog/path",
			wantErr:  false,
		},
		{
			name:     "malformed url",
			inputURL: "://www.boot.dev/blog/path",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "extra slash",
			inputURL: "https://www.boot.dev/blog/path/",
			expected: "www.boot.dev/blog/path",
			wantErr:  false,
		},
		{
			name:     "caps",
			inputURL: "HTTPS://WWW.BOOT.DEV/BLOG/PATH",
			expected: "www.boot.dev/blog/path",
			wantErr:  false,
		},
		{
			name:     "extra slash and caps",
			inputURL: "https://WWW.BOOT.Dev/blog/path/",
			expected: "www.boot.dev/blog/path",
			wantErr:  false,
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if (err != nil) != tc.wantErr {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v, wantErr: %v", i, tc.name, err, tc.wantErr)
				return
			}
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
