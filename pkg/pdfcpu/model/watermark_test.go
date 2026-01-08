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

package model

import (
	"strings"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

func TestWatermarkIsText(t *testing.T) {
	wm := Watermark{Mode: WMText}
	if !wm.IsText() {
		t.Error("IsText() should return true for WMText mode")
	}
	wm.Mode = WMImage
	if wm.IsText() {
		t.Error("IsText() should return false for WMImage mode")
	}
}

func TestWatermarkIsPDF(t *testing.T) {
	wm := Watermark{Mode: WMPDF}
	if !wm.IsPDF() {
		t.Error("IsPDF() should return true for WMPDF mode")
	}
	wm.Mode = WMText
	if wm.IsPDF() {
		t.Error("IsPDF() should return false for WMText mode")
	}
}

func TestWatermarkIsImage(t *testing.T) {
	wm := Watermark{Mode: WMImage}
	if !wm.IsImage() {
		t.Error("IsImage() should return true for WMImage mode")
	}
	wm.Mode = WMText
	if wm.IsImage() {
		t.Error("IsImage() should return false for WMText mode")
	}
}

func TestWatermarkTyp(t *testing.T) {
	tests := []struct {
		mode int
		want string
	}{
		{WMText, "text"},
		{WMImage, "image"},
		{WMPDF, "pdf"},
	}

	for _, tt := range tests {
		wm := Watermark{Mode: tt.mode}
		if got := wm.Typ(); got != tt.want {
			t.Errorf("Typ() for mode %d = %q, want %q", tt.mode, got, tt.want)
		}
	}
}

func TestWatermarkOnTopString(t *testing.T) {
	wm := Watermark{OnTop: false}
	if got := wm.OnTopString(); got != "watermark" {
		t.Errorf("OnTopString() with OnTop=false = %q, want %q", got, "watermark")
	}

	wm.OnTop = true
	if got := wm.OnTopString(); got != "stamp" {
		t.Errorf("OnTopString() with OnTop=true = %q, want %q", got, "stamp")
	}
}

func TestWatermarkMultiStamp(t *testing.T) {
	wm := Watermark{PdfPageNrSrc: 0}
	if !wm.MultiStamp() {
		t.Error("MultiStamp() should return true when PdfPageNrSrc is 0")
	}

	wm.PdfPageNrSrc = 1
	if wm.MultiStamp() {
		t.Error("MultiStamp() should return false when PdfPageNrSrc is not 0")
	}
}

func TestWatermarkString(t *testing.T) {
	wm := DefaultWatermarkConfig()
	wm.TextString = "Test"
	wm.OnTop = true

	s := wm.String()
	if !strings.Contains(s, "Test") {
		t.Error("String() should contain TextString")
	}
	// String() format: "is not on top" or "is  on top" (double space)
	if !strings.Contains(s, "on top") {
		t.Error("String() should contain 'on top'")
	}
}

func TestDefaultWatermarkConfig(t *testing.T) {
	wm := DefaultWatermarkConfig()

	if wm.FontName != "Helvetica" {
		t.Errorf("FontName = %q, want %q", wm.FontName, "Helvetica")
	}
	if wm.FontSize != 24 {
		t.Errorf("FontSize = %d, want %d", wm.FontSize, 24)
	}
	if wm.Pos != types.Center {
		t.Errorf("Pos = %v, want %v", wm.Pos, types.Center)
	}
	if wm.Scale != 0.5 {
		t.Errorf("Scale = %v, want %v", wm.Scale, 0.5)
	}
	if wm.Opacity != 1.0 {
		t.Errorf("Opacity = %v, want %v", wm.Opacity, 1.0)
	}
	if wm.Diagonal != DiagonalLLToUR {
		t.Errorf("Diagonal = %v, want %v", wm.Diagonal, DiagonalLLToUR)
	}
}

func TestLowerLeftCorner(t *testing.T) {
	vp := types.NewRectangle(0, 0, 100, 200)
	bbw, bbh := 20.0, 30.0

	tests := []struct {
		anchor types.Anchor
		wantX  float64
		wantY  float64
	}{
		{types.TopLeft, 0, 170},     // URY - bbh = 200 - 30 = 170
		{types.TopCenter, 40, 170},  // (100/2 - 20/2), 170
		{types.TopRight, 80, 170},   // URX - bbw = 100 - 20 = 80
		{types.Left, 0, 85},         // 0, (200/2 - 30/2) = 85
		{types.Center, 40, 85},      // (100/2 - 20/2), (200/2 - 30/2)
		{types.Right, 80, 85},       // URX - bbw, (200/2 - 30/2)
		{types.BottomLeft, 0, 0},    // LL
		{types.BottomCenter, 40, 0}, // (100/2 - 20/2), 0
		{types.BottomRight, 80, 0},  // URX - bbw, 0
	}

	for _, tt := range tests {
		t.Run(tt.anchor.String(), func(t *testing.T) {
			p := LowerLeftCorner(vp, bbw, bbh, tt.anchor)
			if p.X != tt.wantX {
				t.Errorf("X = %v, want %v", p.X, tt.wantX)
			}
			if p.Y != tt.wantY {
				t.Errorf("Y = %v, want %v", p.Y, tt.wantY)
			}
		})
	}
}

func TestLowerLeftCornerNonZeroOrigin(t *testing.T) {
	// Test with viewport that doesn't start at origin
	vp := types.NewRectangle(10, 20, 110, 220)
	bbw, bbh := 20.0, 30.0

	p := LowerLeftCorner(vp, bbw, bbh, types.BottomLeft)
	if p.X != 10 || p.Y != 20 {
		t.Errorf("BottomLeft = (%v, %v), want (10, 20)", p.X, p.Y)
	}

	p = LowerLeftCorner(vp, bbw, bbh, types.TopRight)
	if p.X != 90 || p.Y != 190 { // 110-20=90, 220-30=190
		t.Errorf("TopRight = (%v, %v), want (90, 190)", p.X, p.Y)
	}
}

func TestWatermarkConstants(t *testing.T) {
	// Test DegToRad and RadToDeg conversions
	deg90 := 90.0
	rad90 := deg90 * DegToRad
	if rad90 < 1.57 || rad90 > 1.58 {
		t.Errorf("90 degrees in radians = %v, expected ~1.5708", rad90)
	}

	deg180 := rad90 * RadToDeg
	if deg180 < 89.9 || deg180 > 90.1 {
		t.Errorf("RadToDeg conversion = %v, expected ~90", deg180)
	}
}

func TestDiagonalConstants(t *testing.T) {
	if NoDiagonal != 0 {
		t.Error("NoDiagonal should be 0")
	}
	if DiagonalLLToUR != 1 {
		t.Error("DiagonalLLToUR should be 1")
	}
	if DiagonalULToLR != 2 {
		t.Error("DiagonalULToLR should be 2")
	}
}

func TestWatermarkModeConstants(t *testing.T) {
	if WMText != 0 {
		t.Error("WMText should be 0")
	}
	if WMImage != 1 {
		t.Error("WMImage should be 1")
	}
	if WMPDF != 2 {
		t.Error("WMPDF should be 2")
	}
}
