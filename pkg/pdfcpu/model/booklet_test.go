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
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

func TestBookletTypeString(t *testing.T) {
	tests := []struct {
		bt   BookletType
		want string
	}{
		{Booklet, "booklet"},
		{BookletAdvanced, "booklet advanced"},
		{BookletPerfectBound, "booklet perfect bound"},
	}

	for _, tt := range tests {
		if got := tt.bt.String(); got != tt.want {
			t.Errorf("BookletType(%d).String() = %q, want %q", tt.bt, got, tt.want)
		}
	}
}

func TestBookletTypeConstants(t *testing.T) {
	if Booklet != 0 {
		t.Error("Booklet should be 0")
	}
	if BookletAdvanced != 1 {
		t.Error("BookletAdvanced should be 1")
	}
	if BookletPerfectBound != 2 {
		t.Error("BookletPerfectBound should be 2")
	}
}

func TestBookletBindingString(t *testing.T) {
	tests := []struct {
		bb   BookletBinding
		want string
	}{
		{LongEdge, "long-edge"},
		{ShortEdge, "short-edge"},
	}

	for _, tt := range tests {
		if got := tt.bb.String(); got != tt.want {
			t.Errorf("BookletBinding(%d).String() = %q, want %q", tt.bb, got, tt.want)
		}
	}
}

func TestBookletBindingConstants(t *testing.T) {
	if LongEdge != 0 {
		t.Error("LongEdge should be 0")
	}
	if ShortEdge != 1 {
		t.Error("ShortEdge should be 1")
	}
}

func TestCutOrFoldString(t *testing.T) {
	nup := &NUp{BookletType: Booklet}
	nupAdv := &NUp{BookletType: BookletAdvanced}

	tests := []struct {
		cf   cutOrFold
		nup  *NUp
		want string
	}{
		{cut, nup, "Cut here"},
		{cut, nupAdv, "Fold & Cut here"},
		{fold, nup, "Fold here"},
		{none, nup, ""},
	}

	for _, tt := range tests {
		if got := tt.cf.String(tt.nup); got != tt.want {
			t.Errorf("cutOrFold(%d).String() = %q, want %q", tt.cf, got, tt.want)
		}
	}
}

func TestGetCutFolds(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		binding  BookletBinding
		portrait bool
		wantHorz cutOrFold
		wantVert cutOrFold
	}{
		{"2up", 2, LongEdge, true, fold, none},
		{"4up long-edge portrait", 4, LongEdge, true, cut, fold},
		{"4up short-edge portrait", 4, ShortEdge, true, fold, cut},
		{"6up", 6, LongEdge, true, cut, fold},
		{"8up long-edge", 8, LongEdge, true, cut, cut},
		{"8up short-edge", 8, ShortEdge, true, cut, fold},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pageDim *types.Dim
			if tt.portrait {
				pageDim = &types.Dim{Width: 595, Height: 842}
			} else {
				pageDim = &types.Dim{Width: 842, Height: 595}
			}

			grid := gridForN(tt.n)
			nup := &NUp{
				Grid:           grid,
				PageDim:        pageDim,
				BookletBinding: tt.binding,
				BookletType:    Booklet,
			}

			horz, vert := getCutFolds(nup)
			if horz != tt.wantHorz {
				t.Errorf("horizontal = %d, want %d", horz, tt.wantHorz)
			}
			if vert != tt.wantVert {
				t.Errorf("vertical = %d, want %d", vert, tt.wantVert)
			}
		})
	}
}

func TestGetCutFoldsPerfectBound(t *testing.T) {
	// Perfect bound converts all folds to cuts
	pageDim := &types.Dim{Width: 595, Height: 842}
	grid := &types.Dim{Width: 2, Height: 1} // 2-up
	nup := &NUp{
		Grid:           grid,
		PageDim:        pageDim,
		BookletBinding: LongEdge,
		BookletType:    BookletPerfectBound,
	}

	horz, _ := getCutFolds(nup)
	if horz != cut {
		t.Errorf("Perfect bound should convert fold to cut, got %d", horz)
	}
}

// Helper function to create grid dimensions for n-up
func gridForN(n int) *types.Dim {
	switch n {
	case 2:
		return &types.Dim{Width: 2, Height: 1}
	case 4:
		return &types.Dim{Width: 2, Height: 2}
	case 6:
		return &types.Dim{Width: 3, Height: 2}
	case 8:
		return &types.Dim{Width: 4, Height: 2}
	default:
		return &types.Dim{Width: 1, Height: 1}
	}
}
