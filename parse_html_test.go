package main

import "testing"

func TestGetHeadingFromHTMLBasic(t *testing.T) {
	tests := []struct {
		name      string
		inputBody string
		expected  string
	}{
		{
			name:      "valid HTML",
			inputBody: "<html><body><h1>Test Title</h1></body></html>",
			expected:  "Test Title",
		},
		{
			name:      "empty HTML",
			inputBody: "",
			expected:  "",
		},
		{
			name:      "malformed HTML",
			inputBody: "<html><body><h1Test Title</h1></body></html>",
			expected:  "",
		},
		{
			name:      "fallback title HTML",
			inputBody: "<html><body><h2>Fallback Title</h2></body></html>",
			expected:  "Fallback Title",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := getHeadingFromHTML(tc.inputBody)
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetFirstParagraphFromHTMLMainPriority(t *testing.T) {
	tests := []struct {
		name      string
		inputBody string
		expected  string
	}{
		{
			name: "valid HTML",
			inputBody: `<html><body>
				<p>Outside paragraph.</p>
				<main>
					<p>Main paragraph.</p>
				</main>
			</body></html>`,
			expected: "Main paragraph.",
		},
		{
			name: "missing main",
			inputBody: `<html><body>
				<p>Outside paragraph.</p>
			</body></html>`,
			expected: "Outside paragraph.",
		},
		{
			name: "malformed HTML",
			inputBody: `<html><body>
				<main>
					<pMain paragraph.</p>
				</main>
			</body></html>`,
			expected: "",
		},
		{
			name: "empty HTML",
			inputBody: `<html><body>
				<main>
					<p></p>
				</main>
			</body></html>`,
			expected: "",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := getFirstParagraphFromHTML(tc.inputBody)
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
