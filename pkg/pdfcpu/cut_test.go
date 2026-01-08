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

func TestParseCutConfigForPoster(t *testing.T) {
	tests := []struct {
		name       string
		config     string
		wantScale  float64
		wantBorder bool
		wantBgCol  bool
	}{
		{
			name:      "formsize A4",
			config:    "formsize:A4",
			wantScale: 1.0,
		},
		{
			name:      "formsize with scale",
			config:    "formsize:A4, scalefactor:2.0",
			wantScale: 2.0,
		},
		{
			name:       "formsize with border",
			config:     "formsize:A4, border:on",
			wantScale:  1.0,
			wantBorder: true,
		},
		{
			name:      "formsize with margin",
			config:    "formsize:A4, margin:10",
			wantScale: 1.0,
		},
		{
			name:      "formsize with bgcolor",
			config:    "formsize:A4, bgcolor:#0000FF",
			wantScale: 1.0,
			wantBgCol: true,
		},
		{
			name:      "dimensions",
			config:    "dimensions:200 300",
			wantScale: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cut, err := ParseCutConfigForPoster(tt.config, types.POINTS)
			if err != nil {
				t.Fatalf("ParseCutConfigForPoster(%q) error = %v", tt.config, err)
			}
			if cut.Scale != tt.wantScale {
				t.Errorf("Scale = %v, want %v", cut.Scale, tt.wantScale)
			}
			if cut.Border != tt.wantBorder {
				t.Errorf("Border = %v, want %v", cut.Border, tt.wantBorder)
			}
			if (cut.BgColor != nil) != tt.wantBgCol {
				t.Errorf("BgColor set = %v, want %v", cut.BgColor != nil, tt.wantBgCol)
			}
		})
	}
}

func TestParseCutConfigForPosterErrors(t *testing.T) {
	tests := []struct {
		name   string
		config string
	}{
		{"empty string", ""},
		{"no colon", "invalid"},
		{"unknown parameter", "unknown:value"},
		{"invalid formsize", "formsize:INVALID"},
		{"invalid scale", "formsize:A4, scalefactor:abc"},
		{"scale less than 1", "formsize:A4, scalefactor:0.5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseCutConfigForPoster(tt.config, types.POINTS)
			if err == nil {
				t.Errorf("ParseCutConfigForPoster(%q) expected error, got nil", tt.config)
			}
		})
	}
}

func TestParseCutConfigForN(t *testing.T) {
	validNs := []int{2, 3, 4, 6, 8, 9, 12, 16}

	for _, n := range validNs {
		t.Run("valid n", func(t *testing.T) {
			cut, err := ParseCutConfigForN(n, "", types.POINTS)
			if err != nil {
				t.Fatalf("ParseCutConfigForN(%d, \"\") error = %v", n, err)
			}
			if cut == nil {
				t.Fatal("cut is nil")
			}
		})
	}

	t.Run("with border config", func(t *testing.T) {
		cut, err := ParseCutConfigForN(4, "border:on", types.POINTS)
		if err != nil {
			t.Fatalf("ParseCutConfigForN(4, border:on) error = %v", err)
		}
		if !cut.Border {
			t.Error("Border should be true")
		}
	})

	t.Run("with margin config", func(t *testing.T) {
		cut, err := ParseCutConfigForN(4, "margin:5", types.POINTS)
		if err != nil {
			t.Fatalf("ParseCutConfigForN(4, margin:5) error = %v", err)
		}
		if cut.Margin != 5 {
			t.Errorf("Margin = %v, want 5", cut.Margin)
		}
	})

	t.Run("with bgcolor config", func(t *testing.T) {
		cut, err := ParseCutConfigForN(4, "bgcolor:#FF0000", types.POINTS)
		if err != nil {
			t.Fatalf("ParseCutConfigForN(4, bgcolor) error = %v", err)
		}
		if cut.BgColor == nil {
			t.Error("BgColor should not be nil")
		}
	})
}

func TestParseCutConfigForNErrors(t *testing.T) {
	invalidNs := []int{0, 1, 5, 7, 10, 11, 13, 14, 15, 17, 100}

	for _, n := range invalidNs {
		t.Run("invalid n", func(t *testing.T) {
			_, err := ParseCutConfigForN(n, "", types.POINTS)
			if err == nil {
				t.Errorf("ParseCutConfigForN(%d, \"\") expected error, got nil", n)
			}
		})
	}

	t.Run("invalid config string", func(t *testing.T) {
		_, err := ParseCutConfigForN(4, "invalid", types.POINTS)
		if err == nil {
			t.Error("ParseCutConfigForN(4, invalid) expected error, got nil")
		}
	})
}

func TestParseCutConfig(t *testing.T) {
	tests := []struct {
		name     string
		config   string
		wantHor  int
		wantVert int
	}{
		{
			name:     "horizontal cut",
			config:   "hor:0.5",
			wantHor:  1,
			wantVert: 0,
		},
		{
			name:     "vertical cut",
			config:   "ver:0.5",
			wantHor:  0,
			wantVert: 1,
		},
		{
			name:     "multiple horizontal cuts",
			config:   "hor:0.3 0.6",
			wantHor:  2,
			wantVert: 0,
		},
		{
			name:     "both cuts",
			config:   "hor:0.5, ver:0.5",
			wantHor:  1,
			wantVert: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cut, err := ParseCutConfig(tt.config, types.POINTS)
			if err != nil {
				t.Fatalf("ParseCutConfig(%q) error = %v", tt.config, err)
			}
			if len(cut.Hor) != tt.wantHor {
				t.Errorf("len(Hor) = %v, want %v", len(cut.Hor), tt.wantHor)
			}
			if len(cut.Vert) != tt.wantVert {
				t.Errorf("len(Vert) = %v, want %v", len(cut.Vert), tt.wantVert)
			}
		})
	}
}

func TestParseCutConfigWithOptions(t *testing.T) {
	t.Run("with border", func(t *testing.T) {
		cut, err := ParseCutConfig("hor:0.5, border:on", types.POINTS)
		if err != nil {
			t.Fatalf("error = %v", err)
		}
		if !cut.Border {
			t.Error("Border should be true")
		}
	})

	t.Run("with margin", func(t *testing.T) {
		cut, err := ParseCutConfig("hor:0.5, margin:10", types.POINTS)
		if err != nil {
			t.Fatalf("error = %v", err)
		}
		if cut.Margin != 10 {
			t.Errorf("Margin = %v, want 10", cut.Margin)
		}
	})

	t.Run("with bgcolor", func(t *testing.T) {
		cut, err := ParseCutConfig("hor:0.5, bgcolor:#00FF00", types.POINTS)
		if err != nil {
			t.Fatalf("error = %v", err)
		}
		if cut.BgColor == nil {
			t.Error("BgColor should not be nil")
		}
	})
}

func TestParseCutConfigErrors(t *testing.T) {
	tests := []struct {
		name   string
		config string
	}{
		{"empty string", ""},
		{"no colon", "invalid"},
		{"unknown parameter", "unknown:value"},
		{"horizontal cut zero", "hor:0"},
		{"horizontal cut one", "hor:1"},
		{"horizontal cut negative", "hor:-0.5"},
		{"horizontal cut greater than one", "hor:1.5"},
		{"vertical cut zero", "ver:0"},
		{"vertical cut one", "ver:1"},
		{"cut non-numeric", "hor:abc"},
		{"invalid border", "hor:0.5, border:maybe"},
		{"negative margin", "hor:0.5, margin:-10"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseCutConfig(tt.config, types.POINTS)
			if err == nil {
				t.Errorf("ParseCutConfig(%q) expected error, got nil", tt.config)
			}
		})
	}
}

func TestParseCutConfigCutValues(t *testing.T) {
	cut, err := ParseCutConfig("hor:0.25 0.5 0.75", types.POINTS)
	if err != nil {
		t.Fatalf("error = %v", err)
	}

	expected := []float64{0.25, 0.5, 0.75}
	if len(cut.Hor) != len(expected) {
		t.Fatalf("len(Hor) = %d, want %d", len(cut.Hor), len(expected))
	}
	for i, v := range expected {
		if cut.Hor[i] != v {
			t.Errorf("Hor[%d] = %v, want %v", i, cut.Hor[i], v)
		}
	}
}
