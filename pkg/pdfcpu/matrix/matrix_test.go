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

package matrix

import (
	"math"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

const tolerance = 1e-9

func floatEquals(a, b float64) bool {
	return math.Abs(a-b) < tolerance
}

func matrixEquals(a, b Matrix) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if !floatEquals(a[i][j], b[i][j]) {
				return false
			}
		}
	}
	return true
}

func TestConstants(t *testing.T) {
	t.Run("DegToRad", func(t *testing.T) {
		// 180 degrees should equal Pi radians
		got := 180 * DegToRad
		want := math.Pi
		if !floatEquals(got, want) {
			t.Errorf("180 * DegToRad = %v, want %v", got, want)
		}
	})

	t.Run("RadToDeg", func(t *testing.T) {
		// Pi radians should equal 180 degrees
		got := math.Pi * RadToDeg
		want := 180.0
		if !floatEquals(got, want) {
			t.Errorf("Pi * RadToDeg = %v, want %v", got, want)
		}
	})

	t.Run("roundtrip", func(t *testing.T) {
		// Converting degrees to radians and back should be identity
		deg := 45.0
		got := deg * DegToRad * RadToDeg
		if !floatEquals(got, deg) {
			t.Errorf("roundtrip %v = %v, want %v", deg, got, deg)
		}
	})
}

func TestIdentMatrix(t *testing.T) {
	want := Matrix{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}}
	if !matrixEquals(IdentMatrix, want) {
		t.Errorf("IdentMatrix = %v, want %v", IdentMatrix, want)
	}
}

func TestMatrixMultiply(t *testing.T) {
	tests := []struct {
		name string
		m, n Matrix
		want Matrix
	}{
		{
			name: "identity times identity",
			m:    IdentMatrix,
			n:    IdentMatrix,
			want: IdentMatrix,
		},
		{
			name: "matrix times identity",
			m:    Matrix{{2, 0, 0}, {0, 3, 0}, {0, 0, 1}},
			n:    IdentMatrix,
			want: Matrix{{2, 0, 0}, {0, 3, 0}, {0, 0, 1}},
		},
		{
			name: "identity times matrix",
			m:    IdentMatrix,
			n:    Matrix{{2, 0, 0}, {0, 3, 0}, {0, 0, 1}},
			want: Matrix{{2, 0, 0}, {0, 3, 0}, {0, 0, 1}},
		},
		{
			name: "scale matrices",
			m:    Matrix{{2, 0, 0}, {0, 2, 0}, {0, 0, 1}},
			n:    Matrix{{3, 0, 0}, {0, 3, 0}, {0, 0, 1}},
			want: Matrix{{6, 0, 0}, {0, 6, 0}, {0, 0, 1}},
		},
		{
			name: "translation matrices",
			m:    Matrix{{1, 0, 0}, {0, 1, 0}, {10, 20, 1}},
			n:    Matrix{{1, 0, 0}, {0, 1, 0}, {5, 7, 1}},
			want: Matrix{{1, 0, 0}, {0, 1, 0}, {15, 27, 1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.Multiply(tt.n)
			if !matrixEquals(got, tt.want) {
				t.Errorf("Multiply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatrixTransform(t *testing.T) {
	tests := []struct {
		name string
		m    Matrix
		p    types.Point
		want types.Point
	}{
		{
			name: "identity transform",
			m:    IdentMatrix,
			p:    types.Point{X: 10, Y: 20},
			want: types.Point{X: 10, Y: 20},
		},
		{
			name: "origin with identity",
			m:    IdentMatrix,
			p:    types.Point{X: 0, Y: 0},
			want: types.Point{X: 0, Y: 0},
		},
		{
			name: "scale by 2",
			m:    Matrix{{2, 0, 0}, {0, 2, 0}, {0, 0, 1}},
			p:    types.Point{X: 5, Y: 10},
			want: types.Point{X: 10, Y: 20},
		},
		{
			name: "translate",
			m:    Matrix{{1, 0, 0}, {0, 1, 0}, {100, 200, 1}},
			p:    types.Point{X: 5, Y: 10},
			want: types.Point{X: 105, Y: 210},
		},
		{
			name: "rotate 90 degrees",
			m:    Matrix{{0, 1, 0}, {-1, 0, 0}, {0, 0, 1}},
			p:    types.Point{X: 1, Y: 0},
			want: types.Point{X: 0, Y: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.Transform(tt.p)
			if !floatEquals(got.X, tt.want.X) || !floatEquals(got.Y, tt.want.Y) {
				t.Errorf("Transform() = (%v, %v), want (%v, %v)", got.X, got.Y, tt.want.X, tt.want.Y)
			}
		})
	}
}

func TestMatrixString(t *testing.T) {
	m := Matrix{{1.5, 2.5, 3.5}, {4.5, 5.5, 6.5}, {7.5, 8.5, 9.5}}
	got := m.String()
	// Just verify it doesn't panic and returns something
	if len(got) == 0 {
		t.Error("String() returned empty string")
	}
	// Verify it contains expected formatted values
	want := "1.50 2.50 3.50\n4.50 5.50 6.50\n7.50 8.50 9.50\n"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestCalcTransformMatrix(t *testing.T) {
	tests := []struct {
		name                     string
		sx, sy, sin, cos, dx, dy float64
		checkFn                  func(Matrix) bool
	}{
		{
			name: "identity transform",
			sx:   1, sy: 1, sin: 0, cos: 1, dx: 0, dy: 0,
			checkFn: func(m Matrix) bool { return matrixEquals(m, IdentMatrix) },
		},
		{
			name: "scale only",
			sx:   2, sy: 3, sin: 0, cos: 1, dx: 0, dy: 0,
			checkFn: func(m Matrix) bool {
				// Verify scaling is applied
				p := m.Transform(types.Point{X: 1, Y: 1})
				return floatEquals(p.X, 2) && floatEquals(p.Y, 3)
			},
		},
		{
			name: "translate only",
			sx:   1, sy: 1, sin: 0, cos: 1, dx: 10, dy: 20,
			checkFn: func(m Matrix) bool {
				p := m.Transform(types.Point{X: 0, Y: 0})
				return floatEquals(p.X, 10) && floatEquals(p.Y, 20)
			},
		},
		{
			name: "rotate 90 degrees",
			sx:   1, sy: 1, sin: 1, cos: 0, dx: 0, dy: 0,
			checkFn: func(m Matrix) bool {
				p := m.Transform(types.Point{X: 1, Y: 0})
				return floatEquals(p.X, 0) && floatEquals(p.Y, 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalcTransformMatrix(tt.sx, tt.sy, tt.sin, tt.cos, tt.dx, tt.dy)
			if !tt.checkFn(got) {
				t.Errorf("CalcTransformMatrix() produced unexpected result: %v", got)
			}
		})
	}
}

func TestCalcRotateAndTranslateTransformMatrix(t *testing.T) {
	tests := []struct {
		name    string
		r       float64
		dx, dy  float64
		checkFn func(Matrix) bool
	}{
		{
			name: "no rotation no translation",
			r:    0, dx: 0, dy: 0,
			checkFn: func(m Matrix) bool {
				p := m.Transform(types.Point{X: 1, Y: 0})
				return floatEquals(p.X, 1) && floatEquals(p.Y, 0)
			},
		},
		{
			name: "90 degree rotation",
			r:    90, dx: 0, dy: 0,
			checkFn: func(m Matrix) bool {
				p := m.Transform(types.Point{X: 1, Y: 0})
				return floatEquals(p.X, 0) && floatEquals(p.Y, 1)
			},
		},
		{
			name: "180 degree rotation",
			r:    180, dx: 0, dy: 0,
			checkFn: func(m Matrix) bool {
				p := m.Transform(types.Point{X: 1, Y: 0})
				return floatEquals(p.X, -1) && floatEquals(p.Y, 0)
			},
		},
		{
			name: "270 degree rotation",
			r:    270, dx: 0, dy: 0,
			checkFn: func(m Matrix) bool {
				p := m.Transform(types.Point{X: 1, Y: 0})
				return floatEquals(p.X, 0) && floatEquals(p.Y, -1)
			},
		},
		{
			name: "translation only",
			r:    0, dx: 100, dy: 200,
			checkFn: func(m Matrix) bool {
				p := m.Transform(types.Point{X: 0, Y: 0})
				return floatEquals(p.X, 100) && floatEquals(p.Y, 200)
			},
		},
		{
			name: "rotation and translation",
			r:    90, dx: 10, dy: 20,
			checkFn: func(m Matrix) bool {
				// Point (1,0) rotated 90 becomes (0,1), then translated to (10, 21)
				p := m.Transform(types.Point{X: 1, Y: 0})
				return floatEquals(p.X, 10) && floatEquals(p.Y, 21)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalcRotateAndTranslateTransformMatrix(tt.r, tt.dx, tt.dy)
			if !tt.checkFn(got) {
				t.Errorf("CalcRotateAndTranslateTransformMatrix(%v, %v, %v) produced unexpected result", tt.r, tt.dx, tt.dy)
			}
		})
	}
}

func TestCalcRotateTransformMatrix(t *testing.T) {
	tests := []struct {
		name string
		rot  float64
		bb   *types.Rectangle
	}{
		{
			name: "no rotation",
			rot:  0,
			bb:   types.RectForDim(100, 100),
		},
		{
			name: "90 degree rotation",
			rot:  90,
			bb:   types.RectForDim(100, 100),
		},
		{
			name: "180 degree rotation",
			rot:  180,
			bb:   types.RectForDim(100, 100),
		},
		{
			name: "270 degree rotation",
			rot:  270,
			bb:   types.RectForDim(100, 100),
		},
		{
			name: "45 degree rotation",
			rot:  45,
			bb:   types.RectForDim(200, 100),
		},
		{
			name: "non-origin bounding box",
			rot:  90,
			bb:   types.RectForWidthAndHeight(50, 50, 100, 100),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify it doesn't panic and returns a valid matrix
			got := CalcRotateTransformMatrix(tt.rot, tt.bb)

			// Verify the matrix has reasonable values (not NaN or Inf)
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					if math.IsNaN(got[i][j]) || math.IsInf(got[i][j], 0) {
						t.Errorf("CalcRotateTransformMatrix() contains invalid value at [%d][%d]: %v", i, j, got[i][j])
					}
				}
			}
		})
	}
}
