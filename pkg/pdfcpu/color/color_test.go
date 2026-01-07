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

package color

import (
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

func TestPredefinedColors(t *testing.T) {
	tests := []struct {
		name    string
		color   SimpleColor
		r, g, b float32
	}{
		{"Black", Black, 0, 0, 0},
		{"White", White, 1, 1, 1},
		{"LightGray", LightGray, 0.9, 0.9, 0.9},
		{"Gray", Gray, 0.5, 0.5, 0.5},
		{"DarkGray", DarkGray, 0.3, 0.3, 0.3},
		{"Red", Red, 1, 0, 0},
		{"Green", Green, 0, 1, 0},
		{"Blue", Blue, 0, 0, 1},
		{"Yellow", Yellow, 0.5, 0.5, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.color.R != tt.r || tt.color.G != tt.g || tt.color.B != tt.b {
				t.Errorf("%s = {%v, %v, %v}, want {%v, %v, %v}",
					tt.name, tt.color.R, tt.color.G, tt.color.B, tt.r, tt.g, tt.b)
			}
		})
	}
}

func TestSimpleColorString(t *testing.T) {
	tests := []struct {
		name  string
		color SimpleColor
		want  string
	}{
		{"black", Black, "r=0.0 g=0.0 b=0.0"},
		{"white", White, "r=1.0 g=1.0 b=1.0"},
		{"red", Red, "r=1.0 g=0.0 b=0.0"},
		{"custom", SimpleColor{0.5, 0.25, 0.75}, "r=0.5 g=0.2 b=0.8"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.color.String()
			if got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSimpleColorArray(t *testing.T) {
	tests := []struct {
		name  string
		color SimpleColor
	}{
		{"black", Black},
		{"white", White},
		{"red", Red},
		{"custom", SimpleColor{0.25, 0.5, 0.75}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arr := tt.color.Array()
			if len(arr) != 3 {
				t.Fatalf("Array() returned %d elements, want 3", len(arr))
			}
			// Verify the values
			r, ok := arr[0].(types.Float)
			if !ok {
				t.Fatalf("arr[0] is not a Float")
			}
			g, ok := arr[1].(types.Float)
			if !ok {
				t.Fatalf("arr[1] is not a Float")
			}
			b, ok := arr[2].(types.Float)
			if !ok {
				t.Fatalf("arr[2] is not a Float")
			}
			if float32(r.Value()) != tt.color.R {
				t.Errorf("R = %v, want %v", r.Value(), tt.color.R)
			}
			if float32(g.Value()) != tt.color.G {
				t.Errorf("G = %v, want %v", g.Value(), tt.color.G)
			}
			if float32(b.Value()) != tt.color.B {
				t.Errorf("B = %v, want %v", b.Value(), tt.color.B)
			}
		})
	}
}

func TestNewSimpleColor(t *testing.T) {
	tests := []struct {
		name    string
		rgb     uint32
		r, g, b float32
	}{
		{"black", 0x000000, 0, 0, 0},
		{"white", 0xFFFFFF, 1, 1, 1},
		{"red", 0xFF0000, 1, 0, 0},
		{"green", 0x00FF00, 0, 1, 0},
		{"blue", 0x0000FF, 0, 0, 1},
		{"mid-gray", 0x808080, 0.5019608, 0.5019608, 0.5019608},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSimpleColor(tt.rgb)
			if got.R != tt.r || got.G != tt.g || got.B != tt.b {
				t.Errorf("NewSimpleColor(0x%06X) = {%v, %v, %v}, want {%v, %v, %v}",
					tt.rgb, got.R, got.G, got.B, tt.r, tt.g, tt.b)
			}
		})
	}
}

func TestNewSimpleColorForArray(t *testing.T) {
	tests := []struct {
		name    string
		arr     types.Array
		r, g, b float32
	}{
		{
			name: "floats",
			arr:  types.Array{types.Float(0.5), types.Float(0.25), types.Float(0.75)},
			r:    0.5, g: 0.25, b: 0.75,
		},
		{
			name: "integers",
			arr:  types.Array{types.Integer(1), types.Integer(0), types.Integer(0)},
			r:    1, g: 0, b: 0,
		},
		{
			name: "mixed",
			arr:  types.Array{types.Float(0.5), types.Integer(1), types.Float(0.25)},
			r:    0.5, g: 1, b: 0.25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSimpleColorForArray(tt.arr)
			if got.R != tt.r || got.G != tt.g || got.B != tt.b {
				t.Errorf("NewSimpleColorForArray() = {%v, %v, %v}, want {%v, %v, %v}",
					got.R, got.G, got.B, tt.r, tt.g, tt.b)
			}
		})
	}
}

func TestNewSimpleColorForHexCode(t *testing.T) {
	tests := []struct {
		name    string
		hexCol  string
		r, g, b float32
		wantErr bool
	}{
		{"white", "#FFFFFF", 1, 1, 1, false},
		{"black", "#000000", 0, 0, 0, false},
		{"red", "#FF0000", 1, 0, 0, false},
		{"green", "#00FF00", 0, 1, 0, false},
		{"blue", "#0000FF", 0, 0, 1, false},
		{"lowercase", "#ffffff", 1, 1, 1, false},
		{"missing hash", "FFFFFF", 0, 0, 0, true},
		{"too short", "#FFF", 0, 0, 0, true},
		{"too long", "#FFFFFFFF", 0, 0, 0, true},
		{"invalid hex", "#GGGGGG", 0, 0, 0, true},
		{"empty", "", 0, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSimpleColorForHexCode(tt.hexCol)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSimpleColorForHexCode(%q) error = %v, wantErr %v", tt.hexCol, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.R != tt.r || got.G != tt.g || got.B != tt.b {
					t.Errorf("NewSimpleColorForHexCode(%q) = {%v, %v, %v}, want {%v, %v, %v}",
						tt.hexCol, got.R, got.G, got.B, tt.r, tt.g, tt.b)
				}
			}
		})
	}
}

func TestInternalSimpleColor(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    SimpleColor
		wantErr bool
	}{
		{"black", "black", Black, false},
		{"Black uppercase", "BLACK", Black, false},
		{"darkgray", "darkgray", DarkGray, false},
		{"gray", "gray", Gray, false},
		{"lightgray", "lightgray", LightGray, false},
		{"white", "white", White, false},
		{"red", "red", Red, false},
		{"green", "green", Green, false},
		{"blue", "blue", Blue, false},
		{"invalid", "purple", SimpleColor{}, true},
		{"empty", "", SimpleColor{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := internalSimpleColor(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("internalSimpleColor(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("internalSimpleColor(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseColor(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    SimpleColor
		wantErr bool
	}{
		// Named colors
		{"named black", "black", Black, false},
		{"named white", "white", White, false},
		{"named red", "red", Red, false},

		// Hex colors
		{"hex white", "#FFFFFF", White, false},
		{"hex black", "#000000", Black, false},
		{"hex red", "#FF0000", Red, false},

		// RGB float values
		{"rgb black", "0 0 0", Black, false},
		{"rgb white", "1 1 1", White, false},
		{"rgb red", "1 0 0", Red, false},
		{"rgb custom", "0.5 0.25 0.75", SimpleColor{0.5, 0.25, 0.75}, false},

		// Error cases
		{"invalid named", "purple", SimpleColor{}, true},
		{"two components", "0.5 0.5", SimpleColor{}, true},
		{"four components", "0.5 0.5 0.5 0.5", SimpleColor{}, true},
		{"invalid red float", "abc 0 0", SimpleColor{}, true},
		{"invalid green float", "0 abc 0", SimpleColor{}, true},
		{"invalid blue float", "0 0 abc", SimpleColor{}, true},
		{"red out of range high", "1.5 0 0", SimpleColor{}, true},
		{"red out of range low", "-0.5 0 0", SimpleColor{}, true},
		{"green out of range", "0 1.5 0", SimpleColor{}, true},
		{"blue out of range", "0 0 1.5", SimpleColor{}, true},
		{"invalid hex", "#GGG", SimpleColor{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseColor(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseColor(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseColor(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestErrInvalidColor(t *testing.T) {
	if ErrInvalidColor == nil {
		t.Error("ErrInvalidColor should not be nil")
	}
	if ErrInvalidColor.Error() == "" {
		t.Error("ErrInvalidColor.Error() should not be empty")
	}
}
