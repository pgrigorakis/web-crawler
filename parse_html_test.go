package main

import (
	"net/url"
	"reflect"
	"testing"
)

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

func TestGetURLsFromHTMLAbsolute(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		inputBody string
		expected  []string
	}{
		{
			name:      "valid HTML",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><a href="https://crawler-test.com"><span>Boot.dev</span></a></body></html>`,
			expected:  []string{"https://crawler-test.com"},
		},
		{
			name:      "empty HTML",
			inputURL:  "https://crawler-test.com",
			inputBody: "",
			expected:  nil,
		},
		{
			name:      "malformed HTML",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><a f="https://crawler-test.com"><span>Boot.dev</span></a></body></html>`,
			expected:  nil,
		},
		{
			name:     "multiple links ",
			inputURL: "https://crawler-test.com",
			inputBody: `<html><body><a href="https://crawler-test.com"><span>Boot.dev</span></a>
			            <a href="https://crawler-test1.com"><span>Boot.dev</span></a>
			            <a href="https://crawler-test2.com"><span>Boot.dev</span></a>
						</body></html>`,
			expected: []string{"https://crawler-test.com", "https://crawler-test1.com", "https://crawler-test2.com"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			baseURL, err := url.Parse(tc.inputURL)
			if err != nil {
				t.Errorf("couldn't parse input URL: %v", err)
				return
			}

			actual, err := getURLsFromHTML(tc.inputBody, baseURL)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}

}

func TestGetImagesFromHTMLRelative(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		inputBody string
		expected  []string
		wantErr   bool
	}{
		{
			name:      "valid HTML",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><img src="/logo.png" alt="Logo"></body></html>`,
			expected:  []string{"https://crawler-test.com/logo.png"},
			wantErr:   false,
		},
		{
			name:      "empty HTML",
			inputURL:  "https://crawler-test.com",
			inputBody: "",
			expected:  nil,
			wantErr:   false,
		},
		{
			name:      "malformed HTML",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><img "/logo.png" alt="Logo"></body></html>`,
			expected:  nil,
			wantErr:   false,
		},
		{
			name:     "multiple images HTML",
			inputURL: "https://crawler-test.com",
			inputBody: `<html><body><img src="/logo.png" alt="Logo">
						<img src="/logo1.png" alt="Logo">
						<img src="/logo2.png" alt="Logo">
						</body></html>`,
			expected: []string{"https://crawler-test.com/logo.png", "https://crawler-test.com/logo1.png", "https://crawler-test.com/logo2.png"},
			wantErr:  false,
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			baseURL, err := url.Parse(tc.inputURL)
			if err != nil {
				t.Errorf("couldn't parse input URL: %v", err)
				return
			}

			actual, err := getImagesFromHTML(tc.inputBody, baseURL)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if (err != nil) != tc.wantErr {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v, wantErr: %v", i, tc.name, err, tc.wantErr)
				return
			}

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}
