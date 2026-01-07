// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filter

import (
	"testing"
)

func TestAbs(t *testing.T) {
	tests := []struct {
		input int
		want  int
	}{
		{0, 0},
		{1, 1},
		{-1, 1},
		{100, 100},
		{-100, 100},
		{1000000, 1000000},
		{-1000000, 1000000},
	}

	for _, tt := range tests {
		got := abs(tt.input)
		if got != tt.want {
			t.Errorf("abs(%d) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestPaeth(t *testing.T) {
	tests := []struct {
		name    string
		a, b, c uint8
		want    uint8
	}{
		// When all values are equal, return a
		{"all equal", 100, 100, 100, 100},
		// When a is closest to prediction (p = a + b - c = 10 + 200 - 200 = 10, pa = 0)
		{"a closest", 10, 200, 200, 10},
		// When b is closest to prediction (p = 200 + 10 - 200 = 10, pb = 0)
		{"b closest", 200, 10, 200, 10},
		// p = 200 + 200 - 10 = 390, pa = 190, pb = 190, pc = 380, so returns a or b (tie goes to a->b order)
		{"c closest", 200, 200, 10, 200},
		// Edge cases
		{"all zero", 0, 0, 0, 0},
		{"all max", 255, 255, 255, 255},
		// p = 1 + 2 - 0 = 3, pa = 2, pb = 1, pc = 3, pb <= pc so return b
		{"png example 1", 1, 2, 0, 2},
		{"png example 2", 0, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := paeth(tt.a, tt.b, tt.c)
			if got != tt.want {
				t.Errorf("paeth(%d, %d, %d) = %d, want %d", tt.a, tt.b, tt.c, got, tt.want)
			}
		})
	}
}

func TestFilterPaeth(t *testing.T) {
	tests := []struct {
		name          string
		cdat          []byte
		pdat          []byte
		bytesPerPixel int
		want          []byte
	}{
		{
			name:          "single byte per pixel",
			cdat:          []byte{1, 2, 3, 4},
			pdat:          []byte{0, 0, 0, 0},
			bytesPerPixel: 1,
			want:          []byte{1, 3, 6, 10},
		},
		{
			name:          "two bytes per pixel",
			cdat:          []byte{1, 2, 1, 2},
			pdat:          []byte{0, 0, 0, 0},
			bytesPerPixel: 2,
			want:          []byte{1, 2, 2, 4},
		},
		{
			name:          "with previous row data",
			cdat:          []byte{1, 1, 1, 1},
			pdat:          []byte{5, 5, 5, 5},
			bytesPerPixel: 1,
			want:          []byte{6, 7, 8, 9},
		},
		{
			name:          "empty slice",
			cdat:          []byte{},
			pdat:          []byte{},
			bytesPerPixel: 1,
			want:          []byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy since filterPaeth modifies in place
			cdat := make([]byte, len(tt.cdat))
			copy(cdat, tt.cdat)

			filterPaeth(cdat, tt.pdat, tt.bytesPerPixel)

			for i, v := range cdat {
				if v != tt.want[i] {
					t.Errorf("filterPaeth result[%d] = %d, want %d", i, v, tt.want[i])
				}
			}
		})
	}
}

func TestIntSizeConstant(t *testing.T) {
	// intSize should be 32 as defined in the code
	if intSize != 32 {
		t.Errorf("intSize = %d, want 32", intSize)
	}
}
