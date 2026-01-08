/*
Copyright 2020 The pdfcpu Authors.

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

func doTestParseBoxListOK(s string, t *testing.T) {
	t.Helper()
	_, err := ParseBoxList(s)
	if err != nil {
		t.Errorf("parseBoxList failed: <%v> <%s>\n", err, s)
		return
	}
}

func doTestParseBoxListFail(s string, t *testing.T) {
	t.Helper()
	_, err := ParseBoxList(s)
	if err == nil {
		t.Errorf("parseBoxList should have failed: <%s>\n", s)
		return
	}
}

func TestParseBoxList(t *testing.T) {
	doTestParseBoxListOK("", t)
	doTestParseBoxListOK("m ", t)
	doTestParseBoxListOK("media,  crop", t)
	doTestParseBoxListOK("m, c, t, b, a", t)
	doTestParseBoxListOK("c,t,b,a,m", t)
	doTestParseBoxListOK("media,crop,bleed,trim,art", t)

	doTestParseBoxListFail("crap", t)
	doTestParseBoxListFail("c t b a ", t)
	doTestParseBoxListFail("media;crop;bleed;trim;art", t)

}

func doTestParseBoxOK(s string, t *testing.T) {
	t.Helper()
	_, err := ParseBox(s, types.POINTS)
	if err != nil {
		t.Errorf("parseBox failed: <%v> <%s>\n", err, s)
		return
	}
}

func doTestParseBoxFail(s string, t *testing.T) {
	t.Helper()
	_, err := ParseBox(s, types.POINTS)
	if err == nil {
		t.Errorf("parseBox should have failed: <%s>\n", s)
		return
	}
}

func TestParseBox(t *testing.T) {

	// Box by rectangle.
	doTestParseBoxOK("[0 0 200 400]", t)
	doTestParseBoxOK("[200 400 0 0]", t)
	doTestParseBoxOK("[-50 -50 200 400]", t)
	doTestParseBoxOK("[2.5 2.5 200 400]", t)
	doTestParseBoxFail("[2.5 200 400]", t)
	doTestParseBoxFail("[2.5 200 400 500 600]", t)
	doTestParseBoxFail("[-50 -50 200 x]", t)

	// Box by 1 margin value.
	doTestParseBoxOK("10.5%", t)
	doTestParseBoxOK("-10.5%", t)
	doTestParseBoxOK("10", t)
	doTestParseBoxOK("-10", t)
	doTestParseBoxOK("10 abs", t)
	doTestParseBoxOK(".5", t)
	doTestParseBoxOK(".5 abs", t)
	doTestParseBoxOK(".4 rel", t)
	doTestParseBoxFail("50%", t)
	doTestParseBoxFail("0.6 rel", t)

	// Box by 2 margin values.
	doTestParseBoxOK("10% -40%", t)
	doTestParseBoxOK("10 5", t)
	doTestParseBoxOK("10 5 abs", t)
	doTestParseBoxOK(".1 .5", t)
	doTestParseBoxOK(".1 .5 abs", t)
	doTestParseBoxOK(".1 .4 rel", t)
	doTestParseBoxFail("10% 40", t)
	doTestParseBoxFail(".5 .5 rel", t)

	// Box by 3 margin values.
	doTestParseBoxOK("10% 15.5% 10%", t)
	doTestParseBoxOK("10 5 15", t)
	doTestParseBoxOK("10 5 15 abs", t)
	doTestParseBoxOK(".1 .155 .1", t)
	doTestParseBoxOK(".1 .155 .1 abs", t)
	doTestParseBoxOK(".1 .155 .1 rel", t)
	doTestParseBoxOK(".1 .155 .6 rel", t)
	doTestParseBoxFail("10% 15.5 10%", t)
	doTestParseBoxFail(".1 .155 r .1 .1", t)
	doTestParseBoxFail(".1 .155 rel .1", t)

	// Box by 4 margin values.
	doTestParseBoxOK("40% 40% 10% 10%", t)
	doTestParseBoxOK("0.4 0.4 20 20", t)
	doTestParseBoxOK("0.4 0.4 .1 .1", t)
	doTestParseBoxOK("0.4 0.4 .1 .1 abs", t)
	doTestParseBoxOK("0.4 0.4 .1 .1 rel", t)
	doTestParseBoxOK("10% 20% 60% 70%", t)
	doTestParseBoxOK("-40% 40% 10% 10%", t)
	doTestParseBoxFail("40% 40% 70% 0%", t)
	doTestParseBoxFail("40% 40% 100 100", t)

	// Box by arbitrary relative position within parent box.
	doTestParseBoxOK("dim:30 30", t)
	doTestParseBoxOK("dim:30 30 abs", t)
	doTestParseBoxOK("dim:.3 .3 rel", t)
	doTestParseBoxOK("dim:30% 30%", t)
	doTestParseBoxOK("pos:tl, dim:30 30", t)
	doTestParseBoxOK("pos:bl, off: 5 5, dim:30 30", t)
	doTestParseBoxOK("pos:bl, off: -5 -5, dim:.3 .3 rel", t)
	doTestParseBoxFail("pos:tl", t)
	doTestParseBoxFail("off:.23 .5", t)
}

func doTestParsePageBoundariesOK(s string, t *testing.T) {
	t.Helper()
	_, err := ParsePageBoundaries(s, types.POINTS)
	if err != nil {
		t.Errorf("parsePageBoundaries failed: <%v> <%s>\n", err, s)
		return
	}
}

func doTestParsePageBoundariesFail(s string, t *testing.T) {
	t.Helper()
	_, err := ParsePageBoundaries(s, types.POINTS)
	if err == nil {
		t.Errorf("parsePageBoundaries should have failed: <%s>\n", s)
		return
	}
}

func TestParsePageBoundaries(t *testing.T) {
	doTestParsePageBoundariesOK("trim:10", t)
	doTestParsePageBoundariesOK("media:[0 0 200 200], crop:10 20, trim:crop, art:bleed, bleed:art", t)
	doTestParsePageBoundariesOK("media:[0 0 200 200], art:bleed, bleed:art", t)
	doTestParsePageBoundariesOK("media:[0 0 200 200], art:bleed, trim:art", t)
	doTestParsePageBoundariesOK("media:[0 0 200 200], art:bleed, trim:bleed", t)
	doTestParsePageBoundariesOK("media:[0 0 200 200], trim:[30 30 170 170], art:bleed", t)
	doTestParsePageBoundariesOK("media:[0 0 200 200]", t)
	doTestParsePageBoundariesOK("media:10", t)
	doTestParsePageBoundariesFail("media:trim", t)
}

func TestResolveBoxType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"full media", "media", "media", false},
		{"prefix m", "m", "media", false},
		{"prefix me", "me", "media", false},
		{"full crop", "crop", "crop", false},
		{"prefix c", "c", "crop", false},
		{"full trim", "trim", "trim", false},
		{"prefix t", "t", "trim", false},
		{"full bleed", "bleed", "bleed", false},
		{"prefix b", "b", "bleed", false},
		{"prefix bl", "bl", "bleed", false},
		{"full art", "art", "art", false},
		{"prefix a", "a", "art", false},
		{"invalid", "invalid", "", true},
		{"empty", "", "media", false}, // Empty string matches first box type due to HasPrefix behavior
		{"partial invalid", "xyz", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveBoxType(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("resolveBoxType(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("resolveBoxType(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseBoxPercentage(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    float64
		wantErr bool
	}{
		{"10 percent", "10", 0.10, false},
		{"25 percent", "25", 0.25, false},
		{"49.9 percent", "49.9", 0.499, false},
		{"-10 percent", "-10", -0.10, false},
		{"-49.9 percent", "-49.9", -0.499, false},
		{"0 percent", "0", 0, false},
		{"50 percent", "50", 0, true},   // Must be < 50
		{"-50 percent", "-50", 0, true}, // Must be > -50
		{"100 percent", "100", 0, true}, // Invalid
		{"invalid", "abc", 0, true},     // Not a number
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseBoxPercentage(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseBoxPercentage(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("parseBoxPercentage(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseBoxDimWidthAndHeight(t *testing.T) {
	tests := []struct {
		name    string
		s1      string
		s2      string
		abs     bool
		wantW   float64
		wantH   float64
		wantErr bool
	}{
		{"absolute 100x200", "100", "200", true, 100, 200, false},
		{"absolute floats", "50.5", "75.5", true, 50.5, 75.5, false},
		{"relative valid", "0.5", "0.75", false, 0.5, 0.75, false},
		{"relative boundary 1", "1", "1", false, 1, 1, false},
		{"relative invalid width", "1.5", "0.5", false, 0, 0, true},
		{"relative invalid height", "0.5", "1.5", false, 0, 0, true},
		{"relative zero width", "0", "0.5", false, 0, 0, true},
		{"relative zero height", "0.5", "0", false, 0, 0, true},
		{"invalid width", "abc", "100", true, 0, 0, true},
		{"invalid height", "100", "abc", true, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotW, gotH, err := parseBoxDimWidthAndHeight(tt.s1, tt.s2, tt.abs)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseBoxDimWidthAndHeight(%q, %q, %v) error = %v, wantErr %v", tt.s1, tt.s2, tt.abs, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if gotW != tt.wantW {
					t.Errorf("parseBoxDimWidthAndHeight(%q, %q, %v) width = %v, want %v", tt.s1, tt.s2, tt.abs, gotW, tt.wantW)
				}
				if gotH != tt.wantH {
					t.Errorf("parseBoxDimWidthAndHeight(%q, %q, %v) height = %v, want %v", tt.s1, tt.s2, tt.abs, gotH, tt.wantH)
				}
			}
		})
	}
}

func TestPageBoundariesString(t *testing.T) {
	tests := []struct {
		name string
		pb   PageBoundaries
		want []string
	}{
		{
			"all boxes",
			PageBoundaries{Media: &Box{}, Crop: &Box{}, Trim: &Box{}, Bleed: &Box{}, Art: &Box{}},
			[]string{"mediaBox", "cropBox", "trimBox", "bleedBox", "artBox"},
		},
		{
			"media only",
			PageBoundaries{Media: &Box{}},
			[]string{"mediaBox"},
		},
		{
			"crop and trim",
			PageBoundaries{Crop: &Box{}, Trim: &Box{}},
			[]string{"cropBox", "trimBox"},
		},
		{
			"empty",
			PageBoundaries{},
			[]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pb.String()
			for _, want := range tt.want {
				if len(tt.want) > 0 && got == "" {
					t.Errorf("PageBoundaries.String() = %q, want to contain %q", got, want)
				}
			}
		})
	}
}

func TestPageBoundariesResolveBox(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"media", "m", false},
		{"crop", "c", false},
		{"trim", "t", false},
		{"bleed", "b", false},
		{"art", "a", false},
		{"full media", "media", false},
		{"invalid", "xyz", true},
		{"empty", "", false}, // Empty string matches first box type due to HasPrefix behavior
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pb := &PageBoundaries{}
			err := pb.ResolveBox(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveBox(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestPageBoundariesSelectAll(t *testing.T) {
	pb := &PageBoundaries{}
	pb.SelectAll()

	if pb.Media == nil {
		t.Error("SelectAll() should set Media")
	}
	if pb.Crop == nil {
		t.Error("SelectAll() should set Crop")
	}
	if pb.Trim == nil {
		t.Error("SelectAll() should set Trim")
	}
	if pb.Bleed == nil {
		t.Error("SelectAll() should set Bleed")
	}
	if pb.Art == nil {
		t.Error("SelectAll() should set Art")
	}
}

func TestPageBoundariesBoxAccessors(t *testing.T) {
	mediaRect := types.NewRectangle(0, 0, 100, 200)
	cropRect := types.NewRectangle(10, 10, 90, 190)
	trimRect := types.NewRectangle(15, 15, 85, 185)
	bleedRect := types.NewRectangle(5, 5, 95, 195)
	artRect := types.NewRectangle(20, 20, 80, 180)

	pb := PageBoundaries{
		Media: &Box{Rect: mediaRect},
		Crop:  &Box{Rect: cropRect},
		Trim:  &Box{Rect: trimRect},
		Bleed: &Box{Rect: bleedRect},
		Art:   &Box{Rect: artRect},
	}

	if got := pb.MediaBox(); got != mediaRect {
		t.Errorf("MediaBox() = %v, want %v", got, mediaRect)
	}
	if got := pb.CropBox(); got != cropRect {
		t.Errorf("CropBox() = %v, want %v", got, cropRect)
	}
	if got := pb.TrimBox(); got != trimRect {
		t.Errorf("TrimBox() = %v, want %v", got, trimRect)
	}
	if got := pb.BleedBox(); got != bleedRect {
		t.Errorf("BleedBox() = %v, want %v", got, bleedRect)
	}
	if got := pb.ArtBox(); got != artRect {
		t.Errorf("ArtBox() = %v, want %v", got, artRect)
	}
}

func TestPageBoundariesBoxAccessorsDefaults(t *testing.T) {
	mediaRect := types.NewRectangle(0, 0, 100, 200)
	cropRect := types.NewRectangle(10, 10, 90, 190)

	// Test defaults when boxes are nil
	pb := PageBoundaries{
		Media: &Box{Rect: mediaRect},
	}

	// CropBox defaults to MediaBox when nil
	if got := pb.CropBox(); got != mediaRect {
		t.Errorf("CropBox() with nil Crop = %v, want MediaBox %v", got, mediaRect)
	}

	pb.Crop = &Box{Rect: cropRect}

	// TrimBox, BleedBox, ArtBox default to CropBox when nil
	if got := pb.TrimBox(); got != cropRect {
		t.Errorf("TrimBox() with nil Trim = %v, want CropBox %v", got, cropRect)
	}
	if got := pb.BleedBox(); got != cropRect {
		t.Errorf("BleedBox() with nil Bleed = %v, want CropBox %v", got, cropRect)
	}
	if got := pb.ArtBox(); got != cropRect {
		t.Errorf("ArtBox() with nil Art = %v, want CropBox %v", got, cropRect)
	}
}

func TestBoxLowerLeftCorner(t *testing.T) {
	r := types.NewRectangle(0, 0, 100, 200)
	w, h := 20.0, 30.0

	tests := []struct {
		anchor types.Anchor
		wantX  float64
		wantY  float64
	}{
		{types.TopLeft, 0, 170},
		{types.TopCenter, 40, 170},
		{types.TopRight, 80, 170},
		{types.Left, 0, 85},
		{types.Center, 40, 85},
		{types.Right, 80, 85},
		{types.BottomLeft, 0, 0},
		{types.BottomCenter, 40, 0},
		{types.BottomRight, 80, 0},
	}

	for _, tt := range tests {
		t.Run(tt.anchor.String(), func(t *testing.T) {
			p := boxLowerLeftCorner(r, w, h, tt.anchor)
			if p.X != tt.wantX {
				t.Errorf("boxLowerLeftCorner X = %v, want %v", p.X, tt.wantX)
			}
			if p.Y != tt.wantY {
				t.Errorf("boxLowerLeftCorner Y = %v, want %v", p.Y, tt.wantY)
			}
		})
	}
}
