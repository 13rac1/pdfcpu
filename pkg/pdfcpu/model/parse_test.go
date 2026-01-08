/*
Copyright 2024 The pdfcpu Authors.

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

package model

import (
	"testing"
)

func TestDecodeNameHexInvalid(t *testing.T) {
	testcases := []string{
		"#",
		"#A",
		"#a",
		"#G0",
		"#00",
		"Fo\x00",
	}
	for _, tc := range testcases {
		if decoded, err := decodeNameHexSequence(tc); err == nil {
			t.Errorf("expected error decoding %s, got %s", tc, decoded)
		}
	}
}

func TestDecodeNameHexValid(t *testing.T) {
	testcases := []struct {
		Input    string
		Expected string
	}{
		{"", ""},
		{"Foo", "Foo"},
		{"A#23", "A#"},
		// Examples from "7.3.5 Name Objects"
		{"Name1", "Name1"},
		{"ASomewhatLongerName", "ASomewhatLongerName"},
		{"A;Name_With-Various***Characters?", "A;Name_With-Various***Characters?"},
		{"1.2", "1.2"},
		{"$$", "$$"},
		{"@pattern", "@pattern"},
		{".notdef", ".notdef"},
		{"Lime#20Green", "Lime Green"},
		{"paired#28#29parentheses", "paired()parentheses"},
		{"The_Key_of_F#23_Minor", "The_Key_of_F#_Minor"},
		{"A#42", "AB"},
	}
	for _, tc := range testcases {
		decoded, err := decodeNameHexSequence(tc.Input)
		if err != nil {
			t.Errorf("decoding %s failed: %s", tc.Input, err)
		} else if decoded != tc.Expected {
			t.Errorf("expected %s when decoding %s, got %s", tc.Expected, tc.Input, decoded)
		}
	}
}

func TestDetectNonEscaped(t *testing.T) {
	testcases := []struct {
		input string
		want  int
	}{
		{"", -1},
		{" ( ", 1},
		{" \\( )", -1},
		{"\\(", -1},
		{"   \\(   ", -1},
		{"\\()(", 3},
		{" \\(\\((abc)", 5},
	}
	for _, tc := range testcases {
		got := detectNonEscaped(tc.input, "(")
		if tc.want != got {
			t.Errorf("%s, want: %d, got: %d", tc.input, tc.want, got)
		}
	}
}

func TestDetectKeywords(t *testing.T) {
	msg := "detectKeywords"

	// process: # gen obj ... obj dict ... {stream ... data ... endstream} endobj
	//                                    streamInd                        endInd
	//                                  -1 if absent                    -1 if absent

	//s := "5 0 obj\n<</Title (xxxxendobjxxxxx)\n/Parent 4 0 R\n/Dest [3 0 R /XYZ 0 738 0]>>\nendobj\n" //78

	s := "1 0 obj\n<<\n /Lang (en-endobject-stream-UK%)  % comment \n>>\nendobj\n\n2 0 obj\n"
	//    0....... ..1 .........2.........3.........4.........5..... ... .6
	endInd, _, err := DetectKeywords(s)
	if err != nil {
		t.Errorf("%s failed: %v", msg, err)
	}
	if endInd != 59 {
		t.Errorf("%s failed: want %d, got %d", msg, 59, endInd)
	}

	// negative test
	s = "1 0 obj\n<<\n /Lang (en-endobject-stream-UK%)  % endobject"
	endInd, _, err = DetectKeywords(s)
	if err != nil {
		t.Errorf("%s failed: %v", msg, err)
	}
	if endInd > 0 {
		t.Errorf("%s failed: want %d, got %d", msg, 0, endInd)
	}

}

func TestHexString(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		want   string
		wantOk bool
	}{
		{"empty string", "", "", true},
		{"valid hex lowercase", "abc123", "ABC123", true},
		{"valid hex uppercase", "ABC123", "ABC123", true},
		{"valid hex mixed", "AbC123", "ABC123", true},
		{"valid with spaces", "AB CD", "ABCD", true},
		{"valid with tabs", "AB\tCD", "ABCD", true},
		{"valid with newlines", "AB\nCD", "ABCD", true},
		{"odd length", "ABC", "ABC0", true},
		{"odd length complex", "A B C", "A0B0C0", true},
		{"all zeros", "000000", "000000", true},
		{"all Fs", "FFFFFF", "FFFFFF", true},
		{"invalid char G", "ABG123", "", false},
		{"invalid char Z", "123Z", "", false},
		{"invalid special char", "AB@CD", "", false},
		{"single digit", "A", "A0", true},
		{"space before odd", "A ", "A0", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := hexString(tt.input)
			if ok != tt.wantOk {
				t.Errorf("hexString(%q) ok = %v, want %v", tt.input, ok, tt.wantOk)
				return
			}
			if tt.wantOk && got != nil && *got != tt.want {
				t.Errorf("hexString(%q) = %q, want %q", tt.input, *got, tt.want)
			}
		})
	}
}

func TestBalancedParenthesesPrefix(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"simple", "(hello)", 6},
		{"empty parens", "()", 1},
		{"nested", "((nested))", 9},
		{"double nested", "(a(b)c)", 6},
		{"escaped open", "(\\()", 3},
		{"escaped close", "(\\))", 3},
		{"escaped backslash", "(\\\\)", 3},
		{"multiple escapes", "(\\(\\))", 5},
		{"unbalanced open", "(()", -1},
		{"unbalanced close", "())", 1},
		{"no parens", "hello", 0}, // Returns 0 when j decrements to 0 immediately (no opening paren)
		{"content with newlines", "(hello\nworld)", 12},
		{"complex nested", "(a(b(c)d)e)", 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := balancedParenthesesPrefix(tt.input)
			if got != tt.want {
				t.Errorf("balancedParenthesesPrefix(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestDelimiter(t *testing.T) {
	tests := []struct {
		char byte
		want bool
	}{
		{'<', true},
		{'>', true},
		{'[', true},
		{']', true},
		{'(', true},
		{')', true},
		{'/', true},
		{' ', false},
		{'\t', false},
		{'\n', false},
		{'a', false},
		{'0', false},
		{'{', false},
		{'}', false},
	}

	for _, tt := range tests {
		t.Run(string(tt.char), func(t *testing.T) {
			got := delimiter(tt.char)
			if got != tt.want {
				t.Errorf("delimiter(%q) = %v, want %v", tt.char, got, tt.want)
			}
		})
	}
}

func TestPositionToNextWhitespaceOrChar(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		chars    string
		wantIdx  int
		wantRest string
	}{
		{"finds space", "hello world", "", 5, " world"},
		{"finds tab", "hello\tworld", "", 5, "\tworld"},
		{"finds char /", "hello/world", "/", 5, "/world"},
		{"finds char <", "hello<world", "<>", 5, "<world"},
		{"no match", "helloworld", "/", -1, "helloworld"},
		{"empty chars", "hello world", "", 5, " world"},
		{"empty input", "", "/", -1, ""}, // Returns -1 for empty string
		{"immediate match", " hello", "", 0, " hello"},
		{"immediate char match", "/hello", "/", 0, "/hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIdx, gotRest := positionToNextWhitespaceOrChar(tt.input, tt.chars)
			if gotIdx != tt.wantIdx {
				t.Errorf("positionToNextWhitespaceOrChar(%q, %q) index = %d, want %d", tt.input, tt.chars, gotIdx, tt.wantIdx)
			}
			if gotRest != tt.wantRest {
				t.Errorf("positionToNextWhitespaceOrChar(%q, %q) rest = %q, want %q", tt.input, tt.chars, gotRest, tt.wantRest)
			}
		})
	}
}

func TestPositionToNextWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantIdx  int
		wantRest string
	}{
		{"finds space", "hello world", 5, " world"},
		{"finds tab", "hello\tworld", 5, "\tworld"},
		{"finds newline", "hello\nworld", 5, "\nworld"},
		{"no whitespace", "helloworld", 0, "helloworld"}, // Returns 0 and original string when no match
		{"empty input", "", 0, ""},                       // Returns 0 and empty string when input is empty
		{"immediate space", " hello", 0, " hello"},
		{"finds null byte", "hello\x00world", 5, "\x00world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIdx, gotRest := positionToNextWhitespace(tt.input)
			if gotIdx != tt.wantIdx {
				t.Errorf("positionToNextWhitespace(%q) index = %d, want %d", tt.input, gotIdx, tt.wantIdx)
			}
			if gotRest != tt.wantRest {
				t.Errorf("positionToNextWhitespace(%q) rest = %q, want %q", tt.input, gotRest, tt.wantRest)
			}
		})
	}
}

func TestPositionToNextEOL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantRest string
		wantIdx  int
	}{
		{"finds LF", "hello\nworld", "\nworld", 5},
		{"finds CR", "hello\rworld", "\rworld", 5},
		{"no EOL", "helloworld", "", 0},
		{"empty input", "", "", 0},
		{"immediate LF", "\nhello", "\nhello", 0},
		{"immediate CR", "\rhello", "\rhello", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRest, gotIdx := positionToNextEOL(tt.input)
			if gotIdx != tt.wantIdx {
				t.Errorf("positionToNextEOL(%q) index = %d, want %d", tt.input, gotIdx, tt.wantIdx)
			}
			if gotRest != tt.wantRest {
				t.Errorf("positionToNextEOL(%q) rest = %q, want %q", tt.input, gotRest, tt.wantRest)
			}
		})
	}
}

func TestTrimLeftSpace(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		relaxed bool
		want    string
		wantEOL bool
	}{
		{"no whitespace", "hello", false, "hello", false},
		{"leading spaces", "   hello", false, "hello", false},
		{"leading tabs", "\t\thello", false, "hello", false},
		{"mixed whitespace", "  \t\n  hello", false, "hello", false},
		{"only whitespace", "   ", false, "", false},
		{"empty string", "", false, "", false},
		{"comment stripped", "  %comment\nhello", false, "hello", false},
		{"multiple comments", "  %first\n%second\nhello", false, "hello", false},
		{"relaxed with newline", "\n hello", true, "hello", true},
		{"relaxed with CR", "\r hello", true, "hello", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotEOL := trimLeftSpace(tt.input, tt.relaxed)
			if got != tt.want {
				t.Errorf("trimLeftSpace(%q, %v) = %q, want %q", tt.input, tt.relaxed, got, tt.want)
			}
			if gotEOL != tt.wantEOL {
				t.Errorf("trimLeftSpace(%q, %v) eol = %v, want %v", tt.input, tt.relaxed, gotEOL, tt.wantEOL)
			}
		})
	}
}

func TestForwardParseBuf(t *testing.T) {
	tests := []struct {
		name string
		buf  string
		pos  int
		want string
	}{
		{"normal advance", "hello", 2, "llo"},
		{"advance to end", "hello", 5, ""},
		{"advance past end", "hello", 10, ""},
		{"no advance", "hello", 0, "hello"},
		{"empty buf", "", 0, ""},
		{"empty buf with pos", "", 5, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := forwardParseBuf(tt.buf, tt.pos)
			if got != tt.want {
				t.Errorf("forwardParseBuf(%q, %d) = %q, want %q", tt.buf, tt.pos, got, tt.want)
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}

func TestNoBuf(t *testing.T) {
	tests := []struct {
		name string
		l    *string
		want bool
	}{
		{"nil pointer", nil, true},
		{"empty string", strPtr(""), true},
		{"non-empty string", strPtr("hello"), false},
		{"whitespace only", strPtr("   "), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := noBuf(tt.l)
			if got != tt.want {
				t.Errorf("noBuf(%v) = %v, want %v", tt.l, got, tt.want)
			}
		})
	}
}

func TestParseBooleanOrNull(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantNil bool
		wantStr string
		wantOk  bool
	}{
		{"true lowercase", "true", false, "true", true},
		{"true uppercase", "TRUE", false, "true", true},
		{"true mixed", "TrUe", false, "true", true},
		{"false lowercase", "false", false, "false", true},
		{"false uppercase", "FALSE", false, "false", true},
		{"null lowercase", "null", true, "null", true},
		{"null uppercase", "NULL", true, "null", true},
		{"not matching", "hello", true, "", false},
		{"too short", "tru", true, "", false},
		{"true with suffix", "truething", false, "true", true},
		{"false with suffix", "falsething", false, "false", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotStr, gotOk := parseBooleanOrNull(tt.input)
			if gotOk != tt.wantOk {
				t.Errorf("parseBooleanOrNull(%q) ok = %v, want %v", tt.input, gotOk, tt.wantOk)
				return
			}
			if gotStr != tt.wantStr {
				t.Errorf("parseBooleanOrNull(%q) str = %q, want %q", tt.input, gotStr, tt.wantStr)
			}
			if tt.wantNil && gotVal != nil {
				t.Errorf("parseBooleanOrNull(%q) val = %v, want nil", tt.input, gotVal)
			}
			if !tt.wantNil && gotVal == nil {
				t.Errorf("parseBooleanOrNull(%q) val = nil, want non-nil", tt.input)
			}
		})
	}
}

func TestPosFloor(t *testing.T) {
	tests := []struct {
		name string
		pos1 int
		pos2 int
		want int
	}{
		{"pos1 smaller", 5, 10, 5},
		{"pos2 smaller", 10, 5, 5},
		{"equal values", 5, 5, 5},
		{"pos1 negative", -1, 5, 5},
		{"pos2 negative", 5, -1, 5},
		{"both negative", -1, -1, -1},
		{"pos1 zero", 0, 5, 0},
		{"pos2 zero", 5, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := posFloor(tt.pos1, tt.pos2)
			if got != tt.want {
				t.Errorf("posFloor(%d, %d) = %d, want %d", tt.pos1, tt.pos2, got, tt.want)
			}
		})
	}
}

func TestIsComment(t *testing.T) {
	tests := []struct {
		name       string
		commentPos int
		strLitPos  int
		want       bool
	}{
		{"comment before string", 5, 10, true},
		{"string before comment", 10, 5, false},
		{"equal positions", 5, 5, false},
		{"comment only", 5, -1, true},
		{"string only", -1, 5, false},
		{"neither present", -1, -1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isComment(tt.commentPos, tt.strLitPos)
			if got != tt.want {
				t.Errorf("isComment(%d, %d) = %v, want %v", tt.commentPos, tt.strLitPos, got, tt.want)
			}
		})
	}
}

func TestIsMarkerTerminated(t *testing.T) {
	tests := []struct {
		r    rune
		want bool
	}{
		{' ', true},
		{'\t', true},
		{'\n', true},
		{'\r', true},
		{0x00, true},
		{'a', false},
		{'0', false},
		{'/', false},
	}

	for _, tt := range tests {
		t.Run(string(tt.r), func(t *testing.T) {
			got := isMarkerTerminated(tt.r)
			if got != tt.want {
				t.Errorf("isMarkerTerminated(%q) = %v, want %v", tt.r, got, tt.want)
			}
		})
	}
}
