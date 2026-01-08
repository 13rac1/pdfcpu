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

func TestNUpN(t *testing.T) {
	tests := []struct {
		width, height float64
		want          int
	}{
		{2, 2, 4},
		{3, 2, 6},
		{2, 3, 6},
		{4, 4, 16},
		{1, 1, 1},
	}

	for _, tt := range tests {
		nup := NUp{Grid: &types.Dim{Width: tt.width, Height: tt.height}}
		if got := nup.N(); got != tt.want {
			t.Errorf("N() for grid %vx%v = %d, want %d", tt.width, tt.height, got, tt.want)
		}
	}
}

func TestNUpIsBooklet(t *testing.T) {
	nup := NUp{BookletType: Booklet}
	if !nup.IsBooklet() {
		t.Error("IsBooklet() should return true for Booklet type")
	}

	nup.BookletType = BookletAdvanced
	if !nup.IsBooklet() {
		t.Error("IsBooklet() should return true for BookletAdvanced type")
	}

	nup.BookletType = BookletPerfectBound
	if nup.IsBooklet() {
		t.Error("IsBooklet() should return false for BookletPerfectBound type")
	}
}

func TestNUpIsTopFoldBinding(t *testing.T) {
	tests := []struct {
		name        string
		portrait    bool
		binding     BookletBinding
		wantTopFold bool
	}{
		{"portrait short-edge", true, ShortEdge, true},
		{"portrait long-edge", true, LongEdge, false},
		{"landscape short-edge", false, ShortEdge, false},
		{"landscape long-edge", false, LongEdge, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pageDim *types.Dim
			if tt.portrait {
				pageDim = &types.Dim{Width: 595, Height: 842} // A4 portrait
			} else {
				pageDim = &types.Dim{Width: 842, Height: 595} // A4 landscape
			}
			nup := NUp{PageDim: pageDim, BookletBinding: tt.binding}
			if got := nup.IsTopFoldBinding(); got != tt.wantTopFold {
				t.Errorf("IsTopFoldBinding() = %v, want %v", got, tt.wantTopFold)
			}
		})
	}
}

func TestNUpString(t *testing.T) {
	nup := NUp{
		PageSize:     "A4",
		PageDim:      &types.Dim{Width: 595, Height: 842},
		Grid:         &types.Dim{Width: 2, Height: 2},
		Orient:       RightDown,
		PageGrid:     false,
		ImgInputFile: true,
	}

	s := nup.String()
	if !strings.Contains(s, "A4") {
		t.Error("String() should contain page size")
	}
	if !strings.Contains(s, "right down") {
		t.Error("String() should contain orientation")
	}
}

func TestDefaultNUpConfig(t *testing.T) {
	nup := DefaultNUpConfig()

	if nup.PageSize != "A4" {
		t.Errorf("PageSize = %q, want %q", nup.PageSize, "A4")
	}
	if nup.Orient != RightDown {
		t.Errorf("Orient = %v, want %v", nup.Orient, RightDown)
	}
	if nup.Margin != 3 {
		t.Errorf("Margin = %v, want %v", nup.Margin, 3)
	}
	if !nup.Border {
		t.Error("Border should be true by default")
	}
	if !nup.Enforce {
		t.Error("Enforce should be true by default")
	}
}

func TestOrientationString(t *testing.T) {
	tests := []struct {
		orient orientation
		want   string
	}{
		{RightDown, "right down"},
		{DownRight, "down right"},
		{LeftDown, "left down"},
		{DownLeft, "down left"},
	}

	for _, tt := range tests {
		if got := tt.orient.String(); got != tt.want {
			t.Errorf("orientation(%d).String() = %q, want %q", tt.orient, got, tt.want)
		}
	}
}

func TestOrientationConstants(t *testing.T) {
	if RightDown != 0 {
		t.Error("RightDown should be 0")
	}
	if DownRight != 1 {
		t.Error("DownRight should be 1")
	}
	if LeftDown != 2 {
		t.Error("LeftDown should be 2")
	}
	if DownLeft != 3 {
		t.Error("DownLeft should be 3")
	}
}

func TestNUpRectsForGrid(t *testing.T) {
	// Test 2x2 grid with RightDown orientation
	nup := NUp{
		PageDim: &types.Dim{Width: 100, Height: 200},
		Grid:    &types.Dim{Width: 2, Height: 2},
		Orient:  RightDown,
	}

	rects := nup.RectsForGrid()
	if len(rects) != 4 {
		t.Fatalf("RectsForGrid() returned %d rects, want 4", len(rects))
	}

	// With RightDown, first row is top row (from left to right)
	// Each cell should be 50x100
	if rects[0].Width() != 50 || rects[0].Height() != 100 {
		t.Errorf("rect[0] dimensions = %vx%v, want 50x100", rects[0].Width(), rects[0].Height())
	}
}

func TestContentBytesForPageRotation(t *testing.T) {
	tests := []struct {
		rot int
		w   float64
		h   float64
	}{
		{0, 100, 200},
		{90, 100, 200},
		{180, 100, 200},
		{270, 100, 200},
		{-90, 100, 200},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			b := ContentBytesForPageRotation(tt.rot, tt.w, tt.h)
			if len(b) == 0 {
				t.Error("ContentBytesForPageRotation() returned empty bytes")
			}
			// Should contain "cm" command
			if !strings.Contains(string(b), "cm") {
				t.Error("ContentBytesForPageRotation() should contain 'cm' command")
			}
		})
	}
}
