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

package font

import (
	"bytes"
	"testing"
)

func TestUint16ToBigEndianBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    uint16
		expected []byte
	}{
		{"zero", 0x0000, []byte{0x00, 0x00}},
		{"low byte only", 0x00FF, []byte{0x00, 0xFF}},
		{"high byte only", 0xFF00, []byte{0xFF, 0x00}},
		{"mixed bytes", 0x1234, []byte{0x12, 0x34}},
		{"max value", 0xFFFF, []byte{0xFF, 0xFF}},
		{"256", 0x0100, []byte{0x01, 0x00}},
		{"255", 0x00FF, []byte{0x00, 0xFF}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := uint16ToBigEndianBytes(tt.input)
			if !bytes.Equal(got, tt.expected) {
				t.Errorf("uint16ToBigEndianBytes(%#04x) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestUint32ToBigEndianBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    uint32
		expected []byte
	}{
		{"zero", 0x00000000, []byte{0x00, 0x00, 0x00, 0x00}},
		{"low byte only", 0x000000FF, []byte{0x00, 0x00, 0x00, 0xFF}},
		{"second byte", 0x0000FF00, []byte{0x00, 0x00, 0xFF, 0x00}},
		{"third byte", 0x00FF0000, []byte{0x00, 0xFF, 0x00, 0x00}},
		{"high byte only", 0xFF000000, []byte{0xFF, 0x00, 0x00, 0x00}},
		{"mixed bytes", 0x12345678, []byte{0x12, 0x34, 0x56, 0x78}},
		{"max value", 0xFFFFFFFF, []byte{0xFF, 0xFF, 0xFF, 0xFF}},
		{"256", 0x00000100, []byte{0x00, 0x00, 0x01, 0x00}},
		{"65536", 0x00010000, []byte{0x00, 0x01, 0x00, 0x00}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := uint32ToBigEndianBytes(tt.input)
			if !bytes.Equal(got, tt.expected) {
				t.Errorf("uint32ToBigEndianBytes(%#08x) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestUtf16BEToString(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{"empty", []byte{}, ""},
		{"ASCII A", []byte{0x00, 0x41}, "A"},
		{"ASCII Hello", []byte{0x00, 0x48, 0x00, 0x65, 0x00, 0x6C, 0x00, 0x6C, 0x00, 0x6F}, "Hello"},
		{"single char", []byte{0x00, 0x58}, "X"},
		{"space", []byte{0x00, 0x20}, " "},
		{"digit 1", []byte{0x00, 0x31}, "1"},
		{"two chars", []byte{0x00, 0x41, 0x00, 0x42}, "AB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utf16BEToString(tt.input)
			if got != tt.expected {
				t.Errorf("utf16BEToString(%v) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestGetNext32BitAlignedLength(t *testing.T) {
	tests := []struct {
		input    uint32
		expected uint32
	}{
		{0, 0},
		{1, 4},
		{2, 4},
		{3, 4},
		{4, 4},
		{5, 8},
		{6, 8},
		{7, 8},
		{8, 8},
		{9, 12},
		{10, 12},
		{11, 12},
		{12, 12},
		{13, 16},
		{100, 100},
		{101, 104},
		{102, 104},
		{103, 104},
		{104, 104},
	}

	for _, tt := range tests {
		got := getNext32BitAlignedLength(tt.input)
		if got != tt.expected {
			t.Errorf("getNext32BitAlignedLength(%d) = %d, want %d", tt.input, got, tt.expected)
		}
	}
}

func TestPad(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		wantLen   int
		checkZero bool // check that padding bytes are zeros
	}{
		{"empty", []byte{}, 0, false},
		{"1 byte", []byte{0x01}, 4, true},
		{"2 bytes", []byte{0x01, 0x02}, 4, true},
		{"3 bytes", []byte{0x01, 0x02, 0x03}, 4, true},
		{"4 bytes", []byte{0x01, 0x02, 0x03, 0x04}, 4, false},
		{"5 bytes", []byte{0x01, 0x02, 0x03, 0x04, 0x05}, 8, true},
		{"6 bytes", []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}, 8, true},
		{"7 bytes", []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, 8, true},
		{"8 bytes", []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}, 8, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy since pad modifies the slice
			input := make([]byte, len(tt.input))
			copy(input, tt.input)

			got := pad(input)

			if len(got) != tt.wantLen {
				t.Errorf("pad() len = %d, want %d", len(got), tt.wantLen)
			}

			// Check original bytes preserved
			for i := 0; i < len(tt.input); i++ {
				if got[i] != tt.input[i] {
					t.Errorf("pad() modified original byte at %d: got %#02x, want %#02x", i, got[i], tt.input[i])
				}
			}

			// Check padding bytes are zeros
			if tt.checkZero {
				for i := len(tt.input); i < len(got); i++ {
					if got[i] != 0x00 {
						t.Errorf("pad() padding byte at %d = %#02x, want 0x00", i, got[i])
					}
				}
			}
		})
	}
}

func TestCalcTableChecksum(t *testing.T) {
	tests := []struct {
		name     string
		tag      string
		data     []byte
		expected uint32
	}{
		{
			name:     "empty",
			tag:      "test",
			data:     []byte{0x00, 0x00, 0x00, 0x00},
			expected: 0,
		},
		{
			name:     "single word",
			tag:      "test",
			data:     []byte{0x00, 0x00, 0x00, 0x01},
			expected: 1,
		},
		{
			name:     "two words",
			tag:      "test",
			data:     []byte{0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02},
			expected: 3,
		},
		{
			name:     "max single word",
			tag:      "test",
			data:     []byte{0xFF, 0xFF, 0xFF, 0xFF},
			expected: 0xFFFFFFFF,
		},
		{
			name:     "mixed values",
			tag:      "test",
			data:     []byte{0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE, 0xF0},
			expected: 0x12345678 + 0x9ABCDEF0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calcTableChecksum(tt.tag, tt.data)
			if got != tt.expected {
				t.Errorf("calcTableChecksum(%q, %v) = %#08x, want %#08x", tt.tag, tt.data, got, tt.expected)
			}
		})
	}
}

func TestCalcTableChecksumHeadTable(t *testing.T) {
	// The "head" table skips the third uint32 (index 2) when calculating checksum
	// This is bytes 8-11 which contain the checksum adjustment
	data := []byte{
		0x00, 0x00, 0x00, 0x01, // index 0: contributes 1
		0x00, 0x00, 0x00, 0x02, // index 1: contributes 2
		0xFF, 0xFF, 0xFF, 0xFF, // index 2: SKIPPED for "head" table
		0x00, 0x00, 0x00, 0x04, // index 3: contributes 4
	}

	// For "head" tag, index 2 is skipped, so sum = 1 + 2 + 4 = 7
	gotHead := calcTableChecksum("head", data)
	expectedHead := uint32(7)
	if gotHead != expectedHead {
		t.Errorf("calcTableChecksum(\"head\", data) = %d, want %d", gotHead, expectedHead)
	}

	// For other tags, index 2 is included, so sum = 1 + 2 + 0xFFFFFFFF + 4
	// This wraps around due to uint32 overflow: (1 + 2 + 4) + 0xFFFFFFFF = 7 + 0xFFFFFFFF = 6 (with overflow)
	gotOther := calcTableChecksum("test", data)
	expectedOther := uint32(6) // 7 + 0xFFFFFFFF wraps to 6
	if gotOther != expectedOther {
		t.Errorf("calcTableChecksum(\"test\", data) = %#08x, want %#08x", gotOther, expectedOther)
	}
}

func TestMyUint32Sort(t *testing.T) {
	t.Run("Len", func(t *testing.T) {
		m := myUint32{1, 2, 3}
		if m.Len() != 3 {
			t.Errorf("Len() = %d, want 3", m.Len())
		}

		empty := myUint32{}
		if empty.Len() != 0 {
			t.Errorf("Len() on empty = %d, want 0", empty.Len())
		}
	})

	t.Run("Less", func(t *testing.T) {
		m := myUint32{1, 5, 3}
		if !m.Less(0, 1) {
			t.Error("Less(0, 1) should be true (1 < 5)")
		}
		if m.Less(1, 0) {
			t.Error("Less(1, 0) should be false (5 < 1)")
		}
		if m.Less(0, 0) {
			t.Error("Less(0, 0) should be false (1 < 1)")
		}
	})

	t.Run("Swap", func(t *testing.T) {
		m := myUint32{1, 2, 3}
		m.Swap(0, 2)
		if m[0] != 3 || m[2] != 1 {
			t.Errorf("After Swap(0, 2): got %v, want [3, 2, 1]", m)
		}
	})
}

func TestToPDFGlyphSpace(t *testing.T) {
	tests := []struct {
		name       string
		unitsPerEm int
		input      int
		expected   int
	}{
		{"1000 units, value 500", 1000, 500, 500},
		{"1000 units, value 1000", 1000, 1000, 1000},
		{"2048 units, value 1024", 2048, 1024, 500},
		{"2048 units, value 2048", 2048, 2048, 1000},
		{"1000 units, value 0", 1000, 0, 0},
		{"2000 units, value 1000", 2000, 1000, 500},
		{"500 units, value 250", 500, 250, 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fd := ttf{UnitsPerEm: tt.unitsPerEm}
			got := fd.toPDFGlyphSpace(tt.input)
			if got != tt.expected {
				t.Errorf("toPDFGlyphSpace(%d) with UnitsPerEm=%d = %d, want %d",
					tt.input, tt.unitsPerEm, got, tt.expected)
			}
		})
	}
}

func TestTableReaders(t *testing.T) {
	// Create a table with known byte patterns
	data := []byte{
		0x12, 0x34, // uint16 at offset 0: 0x1234
		0x56, 0x78, // int16 at offset 2: 0x5678
		0x12, 0x34, 0x56, 0x78, // uint32 at offset 4: 0x12345678
		0x00, 0x01, 0x00, 0x00, // fixed32 at offset 8: 1.0 (65536/65536)
	}
	tbl := table{data: data}

	t.Run("uint16", func(t *testing.T) {
		got := tbl.uint16(0)
		if got != 0x1234 {
			t.Errorf("uint16(0) = %#04x, want 0x1234", got)
		}
	})

	t.Run("int16", func(t *testing.T) {
		got := tbl.int16(2)
		expected := int16(0x5678)
		if got != expected {
			t.Errorf("int16(2) = %d, want %d", got, expected)
		}
	})

	t.Run("uint32", func(t *testing.T) {
		got := tbl.uint32(4)
		if got != 0x12345678 {
			t.Errorf("uint32(4) = %#08x, want 0x12345678", got)
		}
	})

	t.Run("fixed32", func(t *testing.T) {
		got := tbl.fixed32(8)
		expected := 1.0
		if got != expected {
			t.Errorf("fixed32(8) = %f, want %f", got, expected)
		}
	})
}

func TestTableReadersNegativeInt16(t *testing.T) {
	// Test negative int16 values
	data := []byte{0xFF, 0xFF} // -1 in signed 16-bit
	tbl := table{data: data}

	got := tbl.int16(0)
	if got != -1 {
		t.Errorf("int16(0) for 0xFFFF = %d, want -1", got)
	}
}
