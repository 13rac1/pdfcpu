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

package pdfcpu

import (
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

func TestParseResizeConfig(t *testing.T) {
	tests := []struct {
		name       string
		config     string
		unit       types.DisplayUnit
		wantScale  float64
		wantBorder bool
		wantBgCol  bool
	}{
		{
			name:      "scale factor 2.0",
			config:    "scalefactor:2.0",
			unit:      types.POINTS,
			wantScale: 2.0,
		},
		{
			name:      "scale factor 0.5",
			config:    "scalefactor:0.5",
			unit:      types.POINTS,
			wantScale: 0.5,
		},
		{
			name:       "scale with border",
			config:     "scalefactor:2.0, border:on",
			unit:       types.POINTS,
			wantScale:  2.0,
			wantBorder: true,
		},
		{
			name:      "scale with bgcolor",
			config:    "scalefactor:0.5, bgcolor:#00FF00",
			unit:      types.POINTS,
			wantScale: 0.5,
			wantBgCol: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ParseResizeConfig(tt.config, tt.unit)
			if err != nil {
				t.Fatalf("ParseResizeConfig(%q) error = %v", tt.config, err)
			}
			if res.Scale != tt.wantScale {
				t.Errorf("Scale = %v, want %v", res.Scale, tt.wantScale)
			}
			if res.Border != tt.wantBorder {
				t.Errorf("Border = %v, want %v", res.Border, tt.wantBorder)
			}
			if (res.BgColor != nil) != tt.wantBgCol {
				t.Errorf("BgColor set = %v, want %v", res.BgColor != nil, tt.wantBgCol)
			}
		})
	}
}

func TestParseResizeConfigFormSize(t *testing.T) {
	tests := []struct {
		name       string
		config     string
		wantWidth  float64
		wantHeight float64
	}{
		{
			name:       "A4 paper",
			config:     "formsize:A4",
			wantWidth:  595,
			wantHeight: 842,
		},
		{
			name:       "Letter paper",
			config:     "papersize:Letter",
			wantWidth:  612,
			wantHeight: 792,
		},
		{
			name:       "A4 landscape",
			config:     "formsize:A4L",
			wantWidth:  842,
			wantHeight: 595,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ParseResizeConfig(tt.config, types.POINTS)
			if err != nil {
				t.Fatalf("ParseResizeConfig(%q) error = %v", tt.config, err)
			}
			if res.PageDim == nil {
				t.Fatal("PageDim is nil")
			}
			// Check with some tolerance for floating point
			if abs(res.PageDim.Width-tt.wantWidth) > 0.01 {
				t.Errorf("Width = %v, want %v", res.PageDim.Width, tt.wantWidth)
			}
			if abs(res.PageDim.Height-tt.wantHeight) > 0.01 {
				t.Errorf("Height = %v, want %v", res.PageDim.Height, tt.wantHeight)
			}
		})
	}
}

func TestParseResizeConfigDimensions(t *testing.T) {
	res, err := ParseResizeConfig("dimensions:200 300", types.POINTS)
	if err != nil {
		t.Fatalf("ParseResizeConfig(dimensions) error = %v", err)
	}
	if res.PageDim == nil {
		t.Fatal("PageDim is nil")
	}
	if res.PageDim.Width != 200 {
		t.Errorf("Width = %v, want 200", res.PageDim.Width)
	}
	if res.PageDim.Height != 300 {
		t.Errorf("Height = %v, want 300", res.PageDim.Height)
	}
	if !res.UserDim {
		t.Error("UserDim should be true")
	}
}

func TestParseResizeConfigEnforce(t *testing.T) {
	res, err := ParseResizeConfig("formsize:A4, enforce:on", types.POINTS)
	if err != nil {
		t.Fatalf("ParseResizeConfig(enforce) error = %v", err)
	}
	if !res.EnforceOrient {
		t.Error("EnforceOrient should be true")
	}
}

func TestParseResizeConfigErrors(t *testing.T) {
	tests := []struct {
		name   string
		config string
	}{
		{"empty string", ""},
		{"no colon", "invalid"},
		{"scale zero", "scalefactor:0"},
		{"scale one", "scalefactor:1"},
		{"scale negative", "scalefactor:-1"},
		{"scale non-numeric", "scalefactor:abc"},
		{"unknown parameter", "unknown:value"},
		{"scale and formsize conflict", "scalefactor:2.0, formsize:A4"},
		{"dimensions and formsize conflict", "dimensions:200 300, formsize:A4"},
		{"invalid formsize", "formsize:INVALID"},
		{"invalid border value", "scalefactor:2.0, border:maybe"},
		{"invalid bgcolor", "scalefactor:2.0, bgcolor:invalid"},
		{"invalid enforce value", "formsize:A4, enforce:maybe"},
		{"invalid dimensions format", "dimensions:200"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseResizeConfig(tt.config, types.POINTS)
			if err == nil {
				t.Errorf("ParseResizeConfig(%q) expected error, got nil", tt.config)
			}
		})
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
