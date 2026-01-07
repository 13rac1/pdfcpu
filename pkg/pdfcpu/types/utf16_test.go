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
	"testing"
)

func TestIsStringUTF16BE(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"with BOM", "\xFE\xFFtest", true},
		{"without BOM", "test", false},
		{"empty string", "", false},
		{"just BOM", "\xFE\xFF", true},
		{"wrong BOM order", "\xFF\xFE", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsStringUTF16BE(tt.input)
			if got != tt.want {
				t.Errorf("IsStringUTF16BE(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsUTF16BE(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  bool
	}{
		{"with BOM even length", []byte{0xFE, 0xFF, 0x00, 0x41}, true},
		{"without BOM", []byte{0x00, 0x41, 0x00, 0x42}, false},
		{"empty slice", []byte{}, false},
		{"just BOM", []byte{0xFE, 0xFF}, true},
		{"odd length", []byte{0xFE, 0xFF, 0x00}, false},
		{"wrong BOM order", []byte{0xFF, 0xFE, 0x00, 0x41}, false},
		{"nil slice", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsUTF16BE(tt.input)
			if got != tt.want {
				t.Errorf("IsUTF16BE(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestEncodeUTF16String(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", "\xFE\xFF"},
		{"ASCII A", "A", "\xFE\xFF\x00A"},
		{"ASCII AB", "AB", "\xFE\xFF\x00A\x00B"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeUTF16String(tt.input)
			if got != tt.want {
				t.Errorf("EncodeUTF16String(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestDecodeUTF16String(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"ASCII A", "\xFE\xFF\x00A", "A", false},
		{"ASCII AB", "\xFE\xFF\x00A\x00B", "AB", false},
		{"empty with BOM", "\xFE\xFF", "", false},
		{"without BOM", "AB", "", true},
		{"empty", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeUTF16String(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeUTF16String(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DecodeUTF16String(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestEncodeDecodeRoundTrip(t *testing.T) {
	tests := []string{
		"Hello",
		"Hello World",
		"Test 123",
		"",
	}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			encoded := EncodeUTF16String(tt)
			decoded, err := DecodeUTF16String(encoded)
			if err != nil {
				t.Errorf("DecodeUTF16String error: %v", err)
				return
			}
			if decoded != tt {
				t.Errorf("Round trip failed: got %q, want %q", decoded, tt)
			}
		})
	}
}

func TestEscapedUTF16String(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"simple", "test", false},
		{"empty", "", false},
		{"with parens", "test()", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EscapedUTF16String(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("EscapedUTF16String(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got == nil && !tt.wantErr {
				t.Error("EscapedUTF16String returned nil without error")
			}
		})
	}
}

func TestStringLiteralToString(t *testing.T) {
	tests := []struct {
		name    string
		input   StringLiteral
		want    string
		wantErr bool
	}{
		{"simple ASCII", StringLiteral("Hello"), "Hello", false},
		{"empty", StringLiteral(""), "", false},
		{"with escape", StringLiteral("Hello\\nWorld"), "Hello\nWorld", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StringLiteralToString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringLiteralToString(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StringLiteralToString(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestHexLiteralToString(t *testing.T) {
	tests := []struct {
		name    string
		input   HexLiteral
		want    string
		wantErr bool
	}{
		{"simple hex", HexLiteral("48656C6C6F"), "Hello", false},
		{"empty", HexLiteral(""), "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HexLiteralToString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("HexLiteralToString(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HexLiteralToString(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestStringOrHexLiteral(t *testing.T) {
	t.Run("StringLiteral", func(t *testing.T) {
		sl := StringLiteral("Hello")
		got, err := StringOrHexLiteral(sl)
		if err != nil {
			t.Errorf("StringOrHexLiteral error: %v", err)
			return
		}
		if got == nil || *got != "Hello" {
			t.Errorf("StringOrHexLiteral(StringLiteral) = %v, want Hello", got)
		}
	})

	t.Run("HexLiteral", func(t *testing.T) {
		hl := HexLiteral("48656C6C6F")
		got, err := StringOrHexLiteral(hl)
		if err != nil {
			t.Errorf("StringOrHexLiteral error: %v", err)
			return
		}
		if got == nil || *got != "Hello" {
			t.Errorf("StringOrHexLiteral(HexLiteral) = %v, want Hello", got)
		}
	})

	t.Run("invalid type", func(t *testing.T) {
		_, err := StringOrHexLiteral(Integer(1))
		if err == nil {
			t.Error("StringOrHexLiteral(Integer) should return error")
		}
	})
}

func TestErrInvalidUTF16BE(t *testing.T) {
	if ErrInvalidUTF16BE == nil {
		t.Error("ErrInvalidUTF16BE should not be nil")
	}
	if ErrInvalidUTF16BE.Error() == "" {
		t.Error("ErrInvalidUTF16BE.Error() should not be empty")
	}
}
