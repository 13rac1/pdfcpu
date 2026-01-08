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

package filter

import (
	"bytes"
	"testing"
)

func TestIntMemberOf(t *testing.T) {
	tests := []struct {
		name string
		i    int
		list []int
		want bool
	}{
		{"empty list", 5, []int{}, false},
		{"single element found", 5, []int{5}, true},
		{"single element not found", 5, []int{3}, false},
		{"multiple elements found first", 1, []int{1, 2, 3}, true},
		{"multiple elements found middle", 2, []int{1, 2, 3}, true},
		{"multiple elements found last", 3, []int{1, 2, 3}, true},
		{"multiple elements not found", 4, []int{1, 2, 3}, false},
		{"predictor values found", PredictorNone, []int{PredictorNone, PredictorSub, PredictorUp}, true},
		{"predictor values not found", PredictorNo, []int{PredictorNone, PredictorSub, PredictorUp}, false},
		{"negative value found", -1, []int{-2, -1, 0}, true},
		{"zero found", 0, []int{-1, 0, 1}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := intMemberOf(tt.i, tt.list)
			if got != tt.want {
				t.Errorf("intMemberOf(%d, %v) = %v, want %v", tt.i, tt.list, got, tt.want)
			}
		})
	}
}

func TestApplyHorDiff(t *testing.T) {
	tests := []struct {
		name   string
		row    []byte
		colors int
		want   []byte
	}{
		{
			name:   "single color simple",
			row:    []byte{10, 5, 3, 2},
			colors: 1,
			want:   []byte{10, 15, 18, 20},
		},
		{
			name:   "single color zeros",
			row:    []byte{0, 0, 0, 0},
			colors: 1,
			want:   []byte{0, 0, 0, 0},
		},
		{
			name:   "two colors",
			row:    []byte{10, 20, 5, 5},
			colors: 2,
			want:   []byte{10, 20, 15, 25},
		},
		{
			name:   "three colors (RGB)",
			row:    []byte{100, 150, 200, 10, 20, 30},
			colors: 3,
			want:   []byte{100, 150, 200, 110, 170, 230},
		},
		{
			name:   "four colors (RGBA)",
			row:    []byte{50, 100, 150, 200, 5, 10, 15, 20},
			colors: 4,
			want:   []byte{50, 100, 150, 200, 55, 110, 165, 220},
		},
		{
			name:   "overflow wraps",
			row:    []byte{200, 100},
			colors: 1,
			want:   []byte{200, 44}, // 200 + 100 = 300, wraps to 44
		},
		{
			name:   "empty row",
			row:    []byte{},
			colors: 1,
			want:   []byte{},
		},
		{
			name:   "single pixel",
			row:    []byte{42},
			colors: 1,
			want:   []byte{42},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy since applyHorDiff modifies in place
			row := make([]byte, len(tt.row))
			copy(row, tt.row)

			got, err := applyHorDiff(row, tt.colors)
			if err != nil {
				t.Fatalf("applyHorDiff() error = %v", err)
			}
			if !bytes.Equal(got, tt.want) {
				t.Errorf("applyHorDiff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessRow(t *testing.T) {
	tests := []struct {
		name          string
		pr            []byte // previous row
		cr            []byte // current row (first byte is filter type for PNG)
		predictor     int
		colors        int
		bytesPerPixel int
		want          []byte
	}{
		{
			name:          "TIFF predictor",
			pr:            []byte{0, 0, 0, 0},
			cr:            []byte{10, 5, 3, 2},
			predictor:     PredictorTIFF,
			colors:        1,
			bytesPerPixel: 1,
			want:          []byte{10, 15, 18, 20},
		},
		{
			name:          "PNG None filter",
			pr:            []byte{PNGNone, 0, 0, 0, 0},
			cr:            []byte{PNGNone, 1, 2, 3, 4},
			predictor:     PredictorNone,
			colors:        1,
			bytesPerPixel: 1,
			want:          []byte{1, 2, 3, 4},
		},
		{
			name:          "PNG Sub filter",
			pr:            []byte{PNGSub, 0, 0, 0, 0},
			cr:            []byte{PNGSub, 10, 5, 3, 2},
			predictor:     PredictorSub,
			colors:        1,
			bytesPerPixel: 1,
			want:          []byte{10, 15, 18, 20},
		},
		{
			name:          "PNG Up filter",
			pr:            []byte{PNGUp, 10, 20, 30, 40},
			cr:            []byte{PNGUp, 1, 2, 3, 4},
			predictor:     PredictorUp,
			colors:        1,
			bytesPerPixel: 1,
			want:          []byte{11, 22, 33, 44},
		},
		{
			name:          "PNG Average filter simple",
			pr:            []byte{PNGAverage, 0, 0, 0, 0},
			cr:            []byte{PNGAverage, 10, 5, 3, 2},
			predictor:     PredictorAverage,
			colors:        1,
			bytesPerPixel: 1,
			want:          []byte{10, 10, 8, 6},
		},
		{
			name:          "PNG Paeth filter",
			pr:            []byte{PNGPaeth, 0, 0, 0, 0},
			cr:            []byte{PNGPaeth, 1, 2, 3, 4},
			predictor:     PredictorPaeth,
			colors:        1,
			bytesPerPixel: 1,
			want:          []byte{1, 3, 6, 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make copies since processRow modifies in place
			pr := make([]byte, len(tt.pr))
			copy(pr, tt.pr)
			cr := make([]byte, len(tt.cr))
			copy(cr, tt.cr)

			got, err := processRow(pr, cr, tt.predictor, tt.colors, tt.bytesPerPixel)
			if err != nil {
				t.Fatalf("processRow() error = %v", err)
			}
			if !bytes.Equal(got, tt.want) {
				t.Errorf("processRow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlateParameters(t *testing.T) {
	tests := []struct {
		name        string
		parms       map[string]int
		wantColors  int
		wantBPC     int
		wantColumns int
		wantErr     bool
	}{
		{
			name:        "all defaults",
			parms:       nil,
			wantColors:  1,
			wantBPC:     8,
			wantColumns: 1,
			wantErr:     false,
		},
		{
			name:        "empty parms",
			parms:       map[string]int{},
			wantColors:  1,
			wantBPC:     8,
			wantColumns: 1,
			wantErr:     false,
		},
		{
			name:        "custom colors",
			parms:       map[string]int{"Colors": 3},
			wantColors:  3,
			wantBPC:     8,
			wantColumns: 1,
			wantErr:     false,
		},
		{
			name:        "custom BPC 1",
			parms:       map[string]int{"BitsPerComponent": 1},
			wantColors:  1,
			wantBPC:     1,
			wantColumns: 1,
			wantErr:     false,
		},
		{
			name:        "custom BPC 16",
			parms:       map[string]int{"BitsPerComponent": 16},
			wantColors:  1,
			wantBPC:     16,
			wantColumns: 1,
			wantErr:     false,
		},
		{
			name:        "custom columns",
			parms:       map[string]int{"Columns": 100},
			wantColors:  1,
			wantBPC:     8,
			wantColumns: 100,
			wantErr:     false,
		},
		{
			name:        "all custom",
			parms:       map[string]int{"Colors": 4, "BitsPerComponent": 8, "Columns": 640},
			wantColors:  4,
			wantBPC:     8,
			wantColumns: 640,
			wantErr:     false,
		},
		{
			name:        "zero colors error",
			parms:       map[string]int{"Colors": 0},
			wantColors:  0,
			wantBPC:     0,
			wantColumns: 0,
			wantErr:     true,
		},
		{
			name:        "invalid BPC error",
			parms:       map[string]int{"BitsPerComponent": 3},
			wantColors:  0,
			wantBPC:     0,
			wantColumns: 0,
			wantErr:     true,
		},
		{
			name:        "invalid BPC 5",
			parms:       map[string]int{"BitsPerComponent": 5},
			wantColors:  0,
			wantBPC:     0,
			wantColumns: 0,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := flate{baseFilter{parms: tt.parms}}
			colors, bpc, columns, err := f.parameters()

			if (err != nil) != tt.wantErr {
				t.Errorf("parameters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return // Error case, don't check values
			}
			if colors != tt.wantColors {
				t.Errorf("parameters() colors = %d, want %d", colors, tt.wantColors)
			}
			if bpc != tt.wantBPC {
				t.Errorf("parameters() bpc = %d, want %d", bpc, tt.wantBPC)
			}
			if columns != tt.wantColumns {
				t.Errorf("parameters() columns = %d, want %d", columns, tt.wantColumns)
			}
		})
	}
}

func TestCheckBufLen(t *testing.T) {
	tests := []struct {
		name   string
		bufLen int
		maxLen int64
		want   bool
	}{
		{"negative maxLen always true", 100, -1, true},
		{"negative maxLen with empty buf", 0, -1, true},
		{"buf smaller than max", 50, 100, true},
		{"buf equal to max", 100, 100, false},
		{"buf larger than max", 150, 100, false},
		{"empty buf with zero max", 0, 0, false},
		{"empty buf with positive max", 0, 10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b bytes.Buffer
			b.Write(make([]byte, tt.bufLen))
			got := checkBufLen(b, tt.maxLen)
			if got != tt.want {
				t.Errorf("checkBufLen(buf.Len=%d, %d) = %v, want %v", tt.bufLen, tt.maxLen, got, tt.want)
			}
		})
	}
}

func TestPredictorConstants(t *testing.T) {
	// Verify predictor constants have expected values
	tests := []struct {
		name string
		got  int
		want int
	}{
		{"PredictorNo", PredictorNo, 1},
		{"PredictorTIFF", PredictorTIFF, 2},
		{"PredictorNone", PredictorNone, 10},
		{"PredictorSub", PredictorSub, 11},
		{"PredictorUp", PredictorUp, 12},
		{"PredictorAverage", PredictorAverage, 13},
		{"PredictorPaeth", PredictorPaeth, 14},
		{"PredictorOptimum", PredictorOptimum, 15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("%s = %d, want %d", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestPNGFilterConstants(t *testing.T) {
	// Verify PNG filter constants have expected values
	tests := []struct {
		name string
		got  int
		want int
	}{
		{"PNGNone", PNGNone, 0x00},
		{"PNGSub", PNGSub, 0x01},
		{"PNGUp", PNGUp, 0x02},
		{"PNGAverage", PNGAverage, 0x03},
		{"PNGPaeth", PNGPaeth, 0x04},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("%s = %d, want %d", tt.name, tt.got, tt.want)
			}
		})
	}
}
