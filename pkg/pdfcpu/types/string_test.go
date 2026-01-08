/*
Copyright 2022 The pdfcpu Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package types

import (
	"bytes"
	"testing"
)

func TestByteForOctalString(t *testing.T) {
	tests := []struct {
		input    string
		expected byte
	}{
		{
			"001",
			0x1,
		},
		{
			"01",
			0x1,
		},
		{
			"1",
			0x1,
		},
		{
			"010",
			0x8,
		},
		{
			"020",
			0x10,
		},
		{
			"377",
			0xff,
		},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			got := ByteForOctalString(test.input)
			if got != test.expected {
				t.Errorf("got %x; want %x", got, test.expected)
			}
		})
	}
}

func TestUnescapeStringWithOctal(t *testing.T) {
	tests := []struct {
		input    string
		expected []byte
	}{
		{
			"\\5",
			[]byte{0x05},
		},
		{
			"\\5a",
			[]byte{0x05, 'a'},
		},
		{
			"\\5\\5",
			[]byte{0x05, 0x05},
		},
		{
			"\\53",
			[]byte{'+'},
		},
		{
			"\\53a",
			[]byte{'+', 'a'},
		},
		{
			"\\053",
			[]byte{'+'},
		},
		{
			"\\0053",
			[]byte{0x05, '3'},
		},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			got, err := Unescape(test.input)
			if err != nil {
				t.Fail()
			}
			if !bytes.Equal(got, test.expected) {
				t.Errorf("got %x; want %x", got, test.expected)
			}
		})
	}
}

func TestDecodeName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"",
			"",
		},
		{
			"Size",
			"Size",
		},
		{
			"S#69#7a#65",
			"Size",
		},
		{
			"#52#6f#6f#74",
			"Root",
		},
		{
			"#4f#75t#6c#69#6e#65#73",
			"Outlines",
		},
		{
			"C#6fu#6et",
			"Count",
		},
		{
			"K#69#64s",
			"Kids",
		},
		{
			"#50a#72e#6et",
			"Parent",
		},
		{
			"#4d#65di#61#42#6f#78",
			"MediaBox",
		},
		{
			"#46#69#6c#74er",
			"Filter",
		},
		{
			"#46#6ca#74e#44#65c#6fde",
			"FlateDecode",
		},
		{
			"A#53#43#49I#48e#78D#65code",
			"ASCIIHexDecode",
		},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			got, err := DecodeName(test.input)
			if err != nil {
				t.Fail()
			}
			if got != test.expected {
				t.Errorf("got %x; want %x", got, test.expected)
			}
		})
	}
}

func TestEncodeName(t *testing.T) {
	testcases := []struct {
		Input    string
		Expected string
	}{
		{"Foo", "Foo"},
		{"A#", "A#23"},
		{"F#o", "F#23o"},
		{"A;Name_With-Various***Characters?", "A;Name_With-Various***Characters?"},
		{"1.2", "1.2"},
		{"$$", "$$"},
		{"@pattern", "@pattern"},
		{".notdef", ".notdef"},
		{"Lime Green", "Lime#20Green"},
		{"paired()parentheses", "paired#28#29parentheses"},
		{"The_Key_of_F#_Minor", "The_Key_of_F#23_Minor"},
	}
	for _, tc := range testcases {
		if encoded := EncodeName(tc.Input); encoded != tc.Expected {
			t.Errorf("expected %s for %s, got %s", tc.Expected, tc.Input, encoded)
		}
	}
}

func TestRemoveControlChars(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"hello\nworld", "helloworld"},
		{"hello\r\nworld", "helloworld"},
		{"hello\tworld", "helloworld"},
		{"hello\bworld", "helloworld"},
		{"hello\fworld", "helloworld"},
		{"\n\r\t\b\f", ""},
		{"no control chars", "no control chars"},
		{"", ""},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			got := RemoveControlChars(test.input)
			if got != test.expected {
				t.Errorf("got %q; want %q", got, test.expected)
			}
		})
	}
}

func TestNewStringSet(t *testing.T) {
	// Test with nil slice
	set := NewStringSet(nil)
	if len(set) != 0 {
		t.Error("NewStringSet(nil) should return empty set")
	}

	// Test with slice
	slice := []string{"a", "b", "c"}
	set = NewStringSet(slice)
	if len(set) != 3 {
		t.Errorf("NewStringSet() returned set of length %d, want 3", len(set))
	}
	for _, s := range slice {
		if !set[s] {
			t.Errorf("set should contain %q", s)
		}
	}

	// Test with duplicates
	slice = []string{"a", "a", "b"}
	set = NewStringSet(slice)
	if len(set) != 2 {
		t.Errorf("NewStringSet() with duplicates returned set of length %d, want 2", len(set))
	}
}

func TestEscape(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"hello\nworld", "hello\\nworld"},
		{"hello\rworld", "hello\\rworld"},
		{"hello\tworld", "hello\\tworld"},
		{"hello\bworld", "hello\\bworld"},
		{"hello\fworld", "hello\\fworld"},
		{"hello\\world", "hello\\\\world"},
		{"hello(world)", "hello\\(world\\)"},
		{"", ""},
		{"no special", "no special"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			got, err := Escape(test.input)
			if err != nil {
				t.Fatalf("Escape() error: %v", err)
			}
			if *got != test.expected {
				t.Errorf("got %q; want %q", *got, test.expected)
			}
		})
	}
}

func TestUnescapeSpecialChars(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []byte
	}{
		{"newline", "\\n", []byte{0x0A}},
		{"carriage return", "\\r", []byte{0x0D}},
		{"tab", "\\t", []byte{0x09}},
		{"backspace", "\\b", []byte{0x08}},
		{"form feed", "\\f", []byte{0x0C}},
		{"escaped backslash", "\\\\", []byte{'\\'}},
		{"escaped parens", "\\(\\)", []byte{'(', ')'}},
		{"mixed", "hello\\nworld", []byte("hello\nworld")},
		{"line continuation LF", "hello\\\nworld", []byte("helloworld")},
		{"line continuation CRLF", "hello\\\r\nworld", []byte("helloworld")},
		{"regular text", "hello world", []byte("hello world")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Unescape(test.input)
			if err != nil {
				t.Fatalf("Unescape() error: %v", err)
			}
			if !bytes.Equal(got, test.expected) {
				t.Errorf("got %x; want %x", got, test.expected)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "olleh"},
		{"", ""},
		{"a", "a"},
		{"ab", "ba"},
		{"12345", "54321"},
		{"世界", "界世"},   // Unicode
		{"a🎉b", "b🎉a"}, // Emoji
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			got := Reverse(test.input)
			if got != test.expected {
				t.Errorf("got %q; want %q", got, test.expected)
			}
		})
	}
}

func TestTrimLeadingComment(t *testing.T) {
	// TrimLeadingComment returns the original string if non-comment content is found
	// after skipping whitespace. It returns "" if only comments or whitespace.
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{" hello", " hello"},    // Returns original string, doesn't trim
		{"\thello", "\thello"},  // Returns original string
		{"\nhello", "\nhello"},  // Returns original string
		{"%comment", ""},        // Comment found first = empty
		{" %comment", ""},       // Whitespace then comment = empty
		{"  %comment line", ""}, // Whitespace then comment = empty
		{"", ""},                // Empty = empty
		{"   ", ""},             // Only whitespace = empty
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			got := TrimLeadingComment(test.input)
			if got != test.expected {
				t.Errorf("got %q; want %q", got, test.expected)
			}
		})
	}
}

func TestDecodeNameErrors(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"null byte", "hello\x00world", true},
		{"incomplete hex", "A#1", true},
		{"invalid hex", "A#GG", true},
		{"null via hex", "A#00", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := DecodeName(test.input)
			if test.wantErr && err == nil {
				t.Error("expected error but got none")
			}
			if !test.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestEscapeRoundTrip(t *testing.T) {
	tests := []string{
		"hello",
		"hello\nworld",
		"hello\tworld",
		"hello\\world",
		"hello(world)",
		"complex\n\r\t\b\f\\()",
	}

	for _, input := range tests {
		escaped, err := Escape(input)
		if err != nil {
			t.Fatalf("Escape(%q) error: %v", input, err)
		}
		unescaped, err := Unescape(*escaped)
		if err != nil {
			t.Fatalf("Unescape(%q) error: %v", *escaped, err)
		}
		if string(unescaped) != input {
			t.Errorf("roundtrip failed: input=%q escaped=%q unescaped=%q", input, *escaped, string(unescaped))
		}
	}
}

func TestEncodeDecodeNameRoundTrip(t *testing.T) {
	tests := []string{
		"Normal",
		"With Space",
		"With#Hash",
		"With(Parens)",
		"With<Angle>",
		"With[Brackets]",
		"With{Braces}",
		"With/Slash",
		"With%Percent",
	}

	for _, input := range tests {
		encoded := EncodeName(input)
		decoded, err := DecodeName(encoded)
		if err != nil {
			t.Fatalf("DecodeName(%q) error: %v", encoded, err)
		}
		if decoded != input {
			t.Errorf("roundtrip failed: input=%q encoded=%q decoded=%q", input, encoded, decoded)
		}
	}
}
