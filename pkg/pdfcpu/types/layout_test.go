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
package types

import (
	"math"
	"testing"
)

func TestParsePageFormat(t *testing.T) {
	dim, _, err := ParsePageFormat("A3L")
	if err != nil {
		t.Error(err)
	}
	if (dim.Width != 1191) || (dim.Height != 842) {
		t.Errorf("expected 1191x842. got %s", dim)
	}
	// the original dim should be unmodified
	dimOrig := PaperSize["A3"]
	if (dimOrig.Width != 842) || (dimOrig.Height != 1191) {
		t.Errorf("expected origDim=842x1191x842. got %s", dimOrig)
	}
}

func TestParsePageFormatVariants(t *testing.T) {
	tests := []struct {
		input      string
		wantWidth  float64
		wantHeight float64
		wantErr    bool
	}{
		{"A4", 595, 842, false},
		{"A4P", 595, 842, false},
		{"A4L", 842, 595, false},
		{"Letter", 612, 792, false},
		{"LetterL", 792, 612, false},
		{"Legal", 612, 1008, false},
		{"InvalidFormat", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			dim, _, err := ParsePageFormat(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if dim.Width != tt.wantWidth || dim.Height != tt.wantHeight {
				t.Errorf("got %vx%v, want %vx%v", dim.Width, dim.Height, tt.wantWidth, tt.wantHeight)
			}
		})
	}
}

func TestParseHorAlignment(t *testing.T) {
	tests := []struct {
		input   string
		want    HAlignment
		wantErr bool
	}{
		{"l", AlignLeft, false},
		{"left", AlignLeft, false},
		{"LEFT", AlignLeft, false},
		{"r", AlignRight, false},
		{"right", AlignRight, false},
		{"c", AlignCenter, false},
		{"center", AlignCenter, false},
		{"j", AlignJustify, false},
		{"justify", AlignJustify, false},
		{"invalid", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseHorAlignment(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseOrigin(t *testing.T) {
	tests := []struct {
		input   string
		want    Corner
		wantErr bool
	}{
		{"ll", LowerLeft, false},
		{"lowerleft", LowerLeft, false},
		{"LOWERLEFT", LowerLeft, false},
		{"lr", LowerRight, false},
		{"lowerright", LowerRight, false},
		{"ul", UpperLeft, false},
		{"upperleft", UpperLeft, false},
		{"ur", UpperRight, false},
		{"upperright", UpperRight, false},
		{"invalid", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseOrigin(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseAnchor(t *testing.T) {
	tests := []struct {
		input   string
		want    Anchor
		wantErr bool
	}{
		{"tl", TopLeft, false},
		{"topleft", TopLeft, false},
		{"TOPLEFT", TopLeft, false},
		{"tc", TopCenter, false},
		{"topcenter", TopCenter, false},
		{"tr", TopRight, false},
		{"topright", TopRight, false},
		{"l", Left, false},
		{"left", Left, false},
		{"c", Center, false},
		{"center", Center, false},
		{"r", Right, false},
		{"right", Right, false},
		{"bl", BottomLeft, false},
		{"bottomleft", BottomLeft, false},
		{"bc", BottomCenter, false},
		{"bottomcenter", BottomCenter, false},
		{"br", BottomRight, false},
		{"bottomright", BottomRight, false},
		{"invalid", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseAnchor(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParsePositionAnchor(t *testing.T) {
	tests := []struct {
		input   string
		want    Anchor
		wantErr bool
	}{
		{"tl", TopLeft, false},
		{"topleft", TopLeft, false},
		{"top-left", TopLeft, false},
		{"tc", TopCenter, false},
		{"topcenter", TopCenter, false},
		{"top-center", TopCenter, false},
		{"tr", TopRight, false},
		{"topright", TopRight, false},
		{"top-right", TopRight, false},
		{"l", Left, false},
		{"left", Left, false},
		{"c", Center, false},
		{"center", Center, false},
		{"r", Right, false},
		{"right", Right, false},
		{"bl", BottomLeft, false},
		{"bottomleft", BottomLeft, false},
		{"bottom-left", BottomLeft, false},
		{"bc", BottomCenter, false},
		{"bottomcenter", BottomCenter, false},
		{"bottom-center", BottomCenter, false},
		{"br", BottomRight, false},
		{"bottomright", BottomRight, false},
		{"bottom-right", BottomRight, false},
		{"f", Full, false},
		{"full", Full, false},
		{"invalid", 0, true},
		// ParsePositionAnchor is case-sensitive unlike ParseAnchor
		{"TOPLEFT", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParsePositionAnchor(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseRelPosition(t *testing.T) {
	tests := []struct {
		input   string
		want    RelPosition
		wantErr bool
	}{
		{"l", RelPosLeft, false},
		{"left", RelPosLeft, false},
		{"LEFT", RelPosLeft, false},
		{"r", RelPosRight, false},
		{"right", RelPosRight, false},
		{"t", RelPosTop, false},
		{"top", RelPosTop, false},
		{"b", RelPosBottom, false},
		{"bottom", RelPosBottom, false},
		{"invalid", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseRelPosition(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnchorPosition(t *testing.T) {
	r := NewRectangle(0, 0, 100, 200)
	w, h := 20.0, 30.0

	tests := []struct {
		anchor Anchor
		wantX  float64
		wantY  float64
	}{
		{TopLeft, 0, 170},     // 0, 200-30
		{TopCenter, 40, 170},  // (100-20)/2, 200-30
		{TopRight, 80, 170},   // 100-20, 200-30
		{Left, 0, 85},         // 0, (200-30)/2
		{Center, 40, 85},      // (100-20)/2, (200-30)/2
		{Right, 80, 85},       // 100-20, (200-30)/2
		{BottomLeft, 0, 0},    // 0, 0
		{BottomCenter, 40, 0}, // (100-20)/2, 0
		{BottomRight, 80, 0},  // 100-20, 0
	}

	for _, tt := range tests {
		t.Run(tt.anchor.String(), func(t *testing.T) {
			x, y := AnchorPosition(tt.anchor, r, w, h)
			if x != tt.wantX || y != tt.wantY {
				t.Errorf("got (%v, %v), want (%v, %v)", x, y, tt.wantX, tt.wantY)
			}
		})
	}
}

func TestNormalizeCoord(t *testing.T) {
	r := NewRectangle(10, 20, 110, 220) // 100x200 with origin at (10,20)

	tests := []struct {
		name     string
		x, y     float64
		origin   Corner
		absolute bool
		wantX    float64
		wantY    float64
	}{
		// LowerLeft origin (default PDF space)
		{"LowerLeft relative", 50, 50, LowerLeft, false, 50, 50},
		{"LowerLeft absolute", 50, 50, LowerLeft, true, 60, 70},

		// UpperLeft origin (y inverted)
		{"UpperLeft relative", 50, 50, UpperLeft, false, 50, 150},
		{"UpperLeft absolute", 50, 50, UpperLeft, true, 60, 170},
		{"UpperLeft y exceeds height", 50, 250, UpperLeft, false, 50, 0},

		// LowerRight origin (x inverted)
		{"LowerRight relative", 50, 50, LowerRight, false, 50, 50},
		{"LowerRight absolute", 50, 50, LowerRight, true, 60, 70},
		{"LowerRight x exceeds width", 150, 50, LowerRight, false, 0, 50},

		// UpperRight origin (both inverted)
		{"UpperRight relative", 50, 50, UpperRight, false, 50, 150},
		{"UpperRight absolute", 50, 50, UpperRight, true, 60, 170},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotX, gotY := NormalizeCoord(tt.x, tt.y, r, tt.origin, tt.absolute)
			if gotX != tt.wantX || gotY != tt.wantY {
				t.Errorf("got (%v, %v), want (%v, %v)", gotX, gotY, tt.wantX, tt.wantY)
			}
		})
	}
}

func TestNormalizeOffset(t *testing.T) {
	tests := []struct {
		name   string
		x, y   float64
		origin Corner
		wantX  float64
		wantY  float64
	}{
		{"LowerLeft", 10, 20, LowerLeft, 10, 20},
		{"UpperLeft", 10, 20, UpperLeft, 10, -20},
		{"LowerRight", 10, 20, LowerRight, -10, 20},
		{"UpperRight", 10, 20, UpperRight, -10, -20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotX, gotY := NormalizeOffset(tt.x, tt.y, tt.origin)
			if gotX != tt.wantX || gotY != tt.wantY {
				t.Errorf("got (%v, %v), want (%v, %v)", gotX, gotY, tt.wantX, tt.wantY)
			}
		})
	}
}

func TestBestFitRectIntoRect(t *testing.T) {
	// Test when source fits within dest without scaling
	rSrc := NewRectangle(0, 0, 50, 50)
	rDest := NewRectangle(0, 0, 100, 100)
	w, h, dx, dy, rot := BestFitRectIntoRect(rSrc, rDest, false, false)
	if w != 50 || h != 50 {
		t.Errorf("expected no scaling, got w=%v h=%v", w, h)
	}
	if dx != 25 || dy != 25 {
		t.Errorf("expected centered dx=25 dy=25, got dx=%v dy=%v", dx, dy)
	}
	if rot != 0 {
		t.Errorf("expected no rotation, got %v", rot)
	}

	// Test landscape source into landscape dest
	rSrc = NewRectangle(0, 0, 200, 100) // 2:1 landscape
	rDest = NewRectangle(0, 0, 100, 50) // 2:1 landscape
	w, h, dx, dy, rot = BestFitRectIntoRect(rSrc, rDest, false, true)
	if w != 100 || h != 50 {
		t.Errorf("expected w=100 h=50, got w=%v h=%v", w, h)
	}

	// Test portrait source into portrait dest
	rSrc = NewRectangle(0, 0, 100, 200) // 1:2 portrait
	rDest = NewRectangle(0, 0, 50, 100) // 1:2 portrait
	w, h, dx, dy, rot = BestFitRectIntoRect(rSrc, rDest, false, true)
	if w != 50 || h != 100 {
		t.Errorf("expected w=50 h=100, got w=%v h=%v", w, h)
	}

	// Test landscape source into portrait dest with enforce orientation
	rSrc = NewRectangle(0, 0, 200, 100)  // landscape
	rDest = NewRectangle(0, 0, 100, 200) // portrait
	w, h, dx, dy, rot = BestFitRectIntoRect(rSrc, rDest, true, true)
	if rot != 90 {
		t.Errorf("expected 90 degree rotation, got %v", rot)
	}

	// Test portrait source into landscape dest with enforce orientation
	rSrc = NewRectangle(0, 0, 100, 200)  // portrait
	rDest = NewRectangle(0, 0, 200, 100) // landscape
	w, h, dx, dy, rot = BestFitRectIntoRect(rSrc, rDest, true, true)
	if rot != 90 {
		t.Errorf("expected 90 degree rotation, got %v", rot)
	}

	// Test square source
	rSrc = NewRectangle(0, 0, 100, 100) // square
	rDest = NewRectangle(0, 0, 50, 100) // portrait dest
	w, h, _, _, _ = BestFitRectIntoRect(rSrc, rDest, false, true)
	if w != 50 || h != 50 {
		t.Errorf("expected w=50 h=50 for square, got w=%v h=%v", w, h)
	}
}

func TestCornerConstants(t *testing.T) {
	if LowerLeft != 0 {
		t.Error("LowerLeft should be 0")
	}
	if LowerRight != 1 {
		t.Error("LowerRight should be 1")
	}
	if UpperLeft != 2 {
		t.Error("UpperLeft should be 2")
	}
	if UpperRight != 3 {
		t.Error("UpperRight should be 3")
	}
}

func TestHAlignmentConstants(t *testing.T) {
	if AlignLeft != 0 {
		t.Error("AlignLeft should be 0")
	}
	if AlignCenter != 1 {
		t.Error("AlignCenter should be 1")
	}
	if AlignRight != 2 {
		t.Error("AlignRight should be 2")
	}
	if AlignJustify != 3 {
		t.Error("AlignJustify should be 3")
	}
}

func TestVAlignmentConstants(t *testing.T) {
	if AlignBaseline != 0 {
		t.Error("AlignBaseline should be 0")
	}
	if AlignTop != 1 {
		t.Error("AlignTop should be 1")
	}
	if AlignMiddle != 2 {
		t.Error("AlignMiddle should be 2")
	}
	if AlignBottom != 3 {
		t.Error("AlignBottom should be 3")
	}
}

func TestLineJoinStyleConstants(t *testing.T) {
	if LJMiter != 0 {
		t.Error("LJMiter should be 0")
	}
	if LJRound != 1 {
		t.Error("LJRound should be 1")
	}
	if LJBevel != 2 {
		t.Error("LJBevel should be 2")
	}
}

func TestOrientationConstants(t *testing.T) {
	if Horizontal != 0 {
		t.Error("Horizontal should be 0")
	}
	if Vertical != 1 {
		t.Error("Vertical should be 1")
	}
}

func TestRelPositionConstants(t *testing.T) {
	if RelPosLeft != 0 {
		t.Error("RelPosLeft should be 0")
	}
	if RelPosRight != 1 {
		t.Error("RelPosRight should be 1")
	}
	if RelPosTop != 2 {
		t.Error("RelPosTop should be 2")
	}
	if RelPosBottom != 3 {
		t.Error("RelPosBottom should be 3")
	}
}

func floatEqual(a, b float64) bool {
	return math.Abs(a-b) < 0.001
}
