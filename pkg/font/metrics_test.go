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

package font

import (
	"strings"
	"testing"
)

func TestUserSpaceUnits(t *testing.T) {
	tests := []struct {
		name       string
		glyphUnits float64
		fontSize   int
		want       float64
	}{
		{"zero glyph units", 0, 12, 0},
		{"1000 glyph units at 12pt", 1000, 12, 12},
		{"500 glyph units at 12pt", 500, 12, 6},
		{"1000 glyph units at 24pt", 1000, 24, 24},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UserSpaceUnits(tt.glyphUnits, tt.fontSize)
			if got != tt.want {
				t.Errorf("UserSpaceUnits(%v, %d) = %v, want %v", tt.glyphUnits, tt.fontSize, got, tt.want)
			}
		})
	}
}

func TestGlyphSpaceUnits(t *testing.T) {
	tests := []struct {
		name      string
		userUnits float64
		fontSize  int
		want      float64
	}{
		{"zero user units", 0, 12, 0},
		{"12 user units at 12pt", 12, 12, 1000},
		{"6 user units at 12pt", 6, 12, 500},
		{"24 user units at 24pt", 24, 24, 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GlyphSpaceUnits(tt.userUnits, tt.fontSize)
			if got != tt.want {
				t.Errorf("GlyphSpaceUnits(%v, %d) = %v, want %v", tt.userUnits, tt.fontSize, got, tt.want)
			}
		})
	}
}

func TestIsCoreFont(t *testing.T) {
	tests := []struct {
		fontName string
		want     bool
	}{
		{"Helvetica", true},
		{"Helvetica-Bold", true},
		{"Helvetica-Oblique", true},
		{"Helvetica-BoldOblique", true},
		{"Times-Roman", true},
		{"Times-Bold", true},
		{"Courier", true},
		{"Symbol", true},
		{"ZapfDingbats", true},
		{"NotAFont", false},
		{"Arial", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.fontName, func(t *testing.T) {
			got := IsCoreFont(tt.fontName)
			if got != tt.want {
				t.Errorf("IsCoreFont(%q) = %v, want %v", tt.fontName, got, tt.want)
			}
		})
	}
}

func TestCoreFontNames(t *testing.T) {
	names := CoreFontNames()

	// Should have 14 core fonts
	if len(names) != 14 {
		t.Errorf("CoreFontNames() returned %d fonts, want 14", len(names))
	}

	// Check that Helvetica is in the list
	found := false
	for _, name := range names {
		if name == "Helvetica" {
			found = true
			break
		}
	}
	if !found {
		t.Error("CoreFontNames() should include Helvetica")
	}
}

func TestSupportedFont(t *testing.T) {
	tests := []struct {
		fontName string
		want     bool
	}{
		{"Helvetica", true},
		{"Courier", true},
		{"NotAFont", false},
	}

	for _, tt := range tests {
		t.Run(tt.fontName, func(t *testing.T) {
			got := SupportedFont(tt.fontName)
			if got != tt.want {
				t.Errorf("SupportedFont(%q) = %v, want %v", tt.fontName, got, tt.want)
			}
		})
	}
}

func TestIsSupportedFontFile(t *testing.T) {
	tests := []struct {
		filename string
		want     bool
	}{
		{"font.gob", true},
		{"FONT.GOB", true},
		{"font.GOB", true},
		{"font.ttf", false},
		{"font.otf", false},
		{"font", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			got := isSupportedFontFile(tt.filename)
			if got != tt.want {
				t.Errorf("isSupportedFontFile(%q) = %v, want %v", tt.filename, got, tt.want)
			}
		})
	}
}

func TestTTFLightString(t *testing.T) {
	ttf := TTFLight{
		PostscriptName:  "TestFont",
		Protected:       false,
		UnitsPerEm:      1000,
		Ascent:          800,
		Descent:         -200,
		CapHeight:       700,
		FirstChar:       32,
		LastChar:        255,
		LLx:             -100,
		LLy:             -200,
		URx:             1000,
		URy:             800,
		ItalicAngle:     0,
		FixedPitch:      false,
		Bold:            false,
		HorMetricsCount: 256,
		GlyphCount:      256,
		GlyphWidths:     make([]int, 256),
	}

	s := ttf.String()
	if !strings.Contains(s, "TestFont") {
		t.Error("TTFLight.String() should contain PostscriptName")
	}
	if !strings.Contains(s, "1000") {
		t.Error("TTFLight.String() should contain UnitsPerEm")
	}
}

func TestTTFLightSupportsUnicodeBlock(t *testing.T) {
	ttf := TTFLight{
		UnicodeRange: [4]uint32{0x00000001, 0, 0, 0}, // bit 0 set (Basic Latin)
	}

	if !ttf.supportsUnicodeBlock(0) {
		t.Error("supportsUnicodeBlock(0) should return true when bit 0 is set")
	}
	if ttf.supportsUnicodeBlock(1) {
		t.Error("supportsUnicodeBlock(1) should return false when bit 1 is not set")
	}
}

func TestTTFLightUnicodeRangeBits(t *testing.T) {
	ttf := TTFLight{}

	tests := []struct {
		id        string
		wantLen   int
		wantFirst int
	}{
		{"LATN", 4, 0},  // Latin: bits 0, 1, 2, 3
		{"GREK", 1, 7},  // Greek: bit 7
		{"CYRL", 1, 9},  // Cyrillic: bit 9
		{"ARMN", 1, 10}, // Armenian: bit 10
		{"HEBR", 1, 11}, // Hebrew: bit 11
		{"ARAB", 1, 13}, // Arabic: bit 13
		{"DEVA", 1, 15}, // Devanagari: bit 15
		{"BENG", 1, 16}, // Bengali: bit 16
		{"THAI", 1, 24}, // Thai: bit 24
		{"HIRA", 1, 49}, // Hiragana: bit 49
		{"KANA", 1, 50}, // Katakana: bit 50
		{"JPAN", 3, 59}, // Japanese: bits 59, 49, 50
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			bits := ttf.unicodeRangeBits(tt.id)
			if len(bits) != tt.wantLen {
				t.Errorf("unicodeRangeBits(%q) len = %d, want %d", tt.id, len(bits), tt.wantLen)
			}
			if len(bits) > 0 && bits[0] != tt.wantFirst {
				t.Errorf("unicodeRangeBits(%q)[0] = %d, want %d", tt.id, bits[0], tt.wantFirst)
			}
		})
	}

	// Test unknown script
	bits := ttf.unicodeRangeBits("UNKNOWN")
	if bits != nil {
		t.Error("unicodeRangeBits for unknown script should return nil")
	}
}

func TestTTFLightSupportsScript(t *testing.T) {
	ttf := TTFLight{
		UnicodeRange: [4]uint32{0x00000001, 0, 0, 0}, // bit 0 set (Basic Latin)
	}

	t.Run("valid script supported", func(t *testing.T) {
		ok, err := ttf.SupportsScript("LATN")
		if err != nil {
			t.Errorf("SupportsScript(LATN) error = %v", err)
		}
		if !ok {
			t.Error("SupportsScript(LATN) should return true for Latin support")
		}
	})

	t.Run("valid script not supported", func(t *testing.T) {
		ok, err := ttf.SupportsScript("GREK")
		if err != nil {
			t.Errorf("SupportsScript(GREK) error = %v", err)
		}
		if ok {
			t.Error("SupportsScript(GREK) should return false without Greek support")
		}
	})

	t.Run("invalid script length", func(t *testing.T) {
		_, err := ttf.SupportsScript("LA")
		if err == nil {
			t.Error("SupportsScript with invalid length should return error")
		}
	})

	t.Run("unknown script", func(t *testing.T) {
		_, err := ttf.SupportsScript("XXXX")
		if err == nil {
			t.Error("SupportsScript with unknown script should return error")
		}
	})
}

func TestTTFLightGids(t *testing.T) {
	ttf := TTFLight{
		Chars: map[uint32]uint16{
			'A': 1,
			'B': 2,
			'C': 3,
		},
	}

	gids := ttf.Gids()
	if len(gids) != 3 {
		t.Errorf("Gids() len = %d, want 3", len(gids))
	}
}

func TestBoundingBox(t *testing.T) {
	// Test with a core font
	bbox := BoundingBox("Helvetica")
	if bbox == nil {
		t.Error("BoundingBox(Helvetica) should not return nil")
	}
}

func TestTextWidth(t *testing.T) {
	// Test with a core font
	width := TextWidth("Hello", "Helvetica", 12)
	if width <= 0 {
		t.Error("TextWidth should return positive value for non-empty text")
	}

	// Empty text should have zero width
	width = TextWidth("", "Helvetica", 12)
	if width != 0 {
		t.Errorf("TextWidth for empty string = %v, want 0", width)
	}
}

func TestCharWidth(t *testing.T) {
	// Test with a core font
	width := CharWidth("Helvetica", 'A')
	if width <= 0 {
		t.Error("CharWidth(Helvetica, 'A') should return positive value")
	}
}

func TestDescent(t *testing.T) {
	descent := Descent("Helvetica", 12)
	// Descent should be non-negative (it's the absolute value of the negative descent)
	if descent < 0 {
		t.Errorf("Descent should be non-negative, got %v", descent)
	}
}

func TestAscent(t *testing.T) {
	ascent := Ascent("Helvetica", 12)
	if ascent <= 0 {
		t.Error("Ascent should be positive")
	}
}

func TestLineHeight(t *testing.T) {
	lh := LineHeight("Helvetica", 12)
	if lh <= 0 {
		t.Error("LineHeight should be positive")
	}
}

func TestSize(t *testing.T) {
	// Calculate what size is needed to render "Hello" in 100 user units
	size := Size("Hello", "Helvetica", 100)
	if size <= 0 {
		t.Error("Size should return positive value")
	}
}

func TestSizeForLineHeight(t *testing.T) {
	size := SizeForLineHeight("Helvetica", 14)
	if size <= 0 {
		t.Error("SizeForLineHeight should return positive value")
	}
}

func TestUserSpaceFontBBox(t *testing.T) {
	bbox := UserSpaceFontBBox("Helvetica", 12)
	if bbox == nil {
		t.Error("UserSpaceFontBBox should not return nil")
	}
	if bbox.Width() <= 0 {
		t.Error("UserSpaceFontBBox width should be positive")
	}
}
