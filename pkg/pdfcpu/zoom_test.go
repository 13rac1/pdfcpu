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

func TestParseZoomConfig(t *testing.T) {
	tests := []struct {
		name       string
		config     string
		unit       types.DisplayUnit
		wantFactor float64
		wantBorder bool
		wantBgCol  bool
	}{
		{
			name:       "zoom factor 2.0",
			config:     "factor:2.0",
			unit:       types.POINTS,
			wantFactor: 2.0,
		},
		{
			name:       "zoom factor 0.5",
			config:     "factor:0.5",
			unit:       types.POINTS,
			wantFactor: 0.5,
		},
		{
			name:       "factor with border on",
			config:     "factor:2.0, border:on",
			unit:       types.POINTS,
			wantFactor: 2.0,
			wantBorder: true,
		},
		{
			name:       "factor with border true",
			config:     "factor:2.0, border:true",
			unit:       types.POINTS,
			wantFactor: 2.0,
			wantBorder: true,
		},
		{
			name:       "factor with bgcolor",
			config:     "factor:0.5, bgcolor:#FF0000",
			unit:       types.POINTS,
			wantFactor: 0.5,
			wantBgCol:  true,
		},
		{
			name:       "prefix completion f -> factor",
			config:     "f:2.0",
			unit:       types.POINTS,
			wantFactor: 2.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zoom, err := ParseZoomConfig(tt.config, tt.unit)
			if err != nil {
				t.Fatalf("ParseZoomConfig(%q) error = %v", tt.config, err)
			}
			if zoom.Factor != tt.wantFactor {
				t.Errorf("Factor = %v, want %v", zoom.Factor, tt.wantFactor)
			}
			if zoom.Border != tt.wantBorder {
				t.Errorf("Border = %v, want %v", zoom.Border, tt.wantBorder)
			}
			if (zoom.BgColor != nil) != tt.wantBgCol {
				t.Errorf("BgColor set = %v, want %v", zoom.BgColor != nil, tt.wantBgCol)
			}
		})
	}
}

func TestParseZoomConfigHMargin(t *testing.T) {
	zoom, err := ParseZoomConfig("hmargin:10", types.POINTS)
	if err != nil {
		t.Fatalf("ParseZoomConfig(hmargin:10) error = %v", err)
	}
	if zoom.HMargin != 10 {
		t.Errorf("HMargin = %v, want 10", zoom.HMargin)
	}
}

func TestParseZoomConfigVMargin(t *testing.T) {
	zoom, err := ParseZoomConfig("vmargin:10", types.POINTS)
	if err != nil {
		t.Fatalf("ParseZoomConfig(vmargin:10) error = %v", err)
	}
	if zoom.VMargin != 10 {
		t.Errorf("VMargin = %v, want 10", zoom.VMargin)
	}
}

func TestParseZoomConfigErrors(t *testing.T) {
	tests := []struct {
		name   string
		config string
	}{
		{"empty string", ""},
		{"no colon", "invalid"},
		{"factor zero", "factor:0"},
		{"factor one", "factor:1"},
		{"factor negative", "factor:-1"},
		{"factor non-numeric", "factor:abc"},
		{"unknown parameter", "unknown:value"},
		{"factor and hmargin conflict", "factor:2.0, hmargin:10"},
		{"factor and vmargin conflict", "factor:2.0, vmargin:10"},
		{"hmargin and vmargin both", "hmargin:10, vmargin:10"},
		{"hmargin zero", "hmargin:0"},
		{"vmargin zero", "vmargin:0"},
		{"invalid border value", "factor:2.0, border:maybe"},
		{"invalid bgcolor", "factor:2.0, bgcolor:invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseZoomConfig(tt.config, types.POINTS)
			if err == nil {
				t.Errorf("ParseZoomConfig(%q) expected error, got nil", tt.config)
			}
		})
	}
}
