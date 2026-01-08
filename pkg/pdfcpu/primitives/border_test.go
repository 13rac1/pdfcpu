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

package primitives

import (
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/color"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

func TestBorderCalc(t *testing.T) {
	tests := []struct {
		name      string
		border    Border
		wantWidth float64
		wantColor *color.SimpleColor
	}{
		{
			name:      "nil color returns Black default",
			border:    Border{Width: 5, col: nil},
			wantWidth: 0,
			wantColor: &color.Black,
		},
		{
			name:      "with color returns width and color",
			border:    Border{Width: 10, col: &color.Red},
			wantWidth: 10.0,
			wantColor: &color.Red,
		},
		{
			name:      "zero width with color",
			border:    Border{Width: 0, col: &color.Blue},
			wantWidth: 0,
			wantColor: &color.Blue,
		},
		{
			name:      "custom color",
			border:    Border{Width: 3, col: &color.Green},
			wantWidth: 3.0,
			wantColor: &color.Green,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWidth, gotColor := tt.border.calc()

			if gotWidth != tt.wantWidth {
				t.Errorf("Border.calc() width = %.1f, want %.1f", gotWidth, tt.wantWidth)
			}

			if gotColor == nil {
				t.Fatal("Border.calc() returned nil color")
			}

			if gotColor.R != tt.wantColor.R || gotColor.G != tt.wantColor.G || gotColor.B != tt.wantColor.B {
				t.Errorf("Border.calc() color = (%f, %f, %f), want (%f, %f, %f)",
					gotColor.R, gotColor.G, gotColor.B,
					tt.wantColor.R, tt.wantColor.G, tt.wantColor.B)
			}
		})
	}
}

func TestBorderMergeIn(t *testing.T) {
	tests := []struct {
		name      string
		border    *Border
		b0        *Border
		wantWidth int
		wantColor *color.SimpleColor
		wantStyle types.LineJoinStyle
	}{
		{
			name:      "zero width inherits from b0",
			border:    &Border{Width: 0, col: &color.Red, style: types.LJRound},
			b0:        &Border{Width: 5, col: &color.Blue, style: types.LJBevel},
			wantWidth: 5,
			wantColor: &color.Red,
			wantStyle: types.LJRound,
		},
		{
			name:      "non-zero width keeps own value",
			border:    &Border{Width: 10, col: &color.Red, style: types.LJRound},
			b0:        &Border{Width: 5, col: &color.Blue, style: types.LJBevel},
			wantWidth: 10,
			wantColor: &color.Red,
			wantStyle: types.LJRound,
		},
		{
			name:      "nil color inherits from b0",
			border:    &Border{Width: 10, col: nil, style: types.LJRound},
			b0:        &Border{Width: 5, col: &color.Green, style: types.LJBevel},
			wantWidth: 10,
			wantColor: &color.Green,
			wantStyle: types.LJRound,
		},
		{
			name:      "non-nil color keeps own value",
			border:    &Border{Width: 10, col: &color.Red, style: types.LJRound},
			b0:        &Border{Width: 5, col: &color.Blue, style: types.LJBevel},
			wantWidth: 10,
			wantColor: &color.Red,
			wantStyle: types.LJRound,
		},
		{
			name:      "LJMiter style inherits from b0",
			border:    &Border{Width: 10, col: &color.Red, style: types.LJMiter},
			b0:        &Border{Width: 5, col: &color.Blue, style: types.LJBevel},
			wantWidth: 10,
			wantColor: &color.Red,
			wantStyle: types.LJBevel,
		},
		{
			name:      "non-LJMiter style keeps own value",
			border:    &Border{Width: 10, col: &color.Red, style: types.LJRound},
			b0:        &Border{Width: 5, col: &color.Blue, style: types.LJBevel},
			wantWidth: 10,
			wantColor: &color.Red,
			wantStyle: types.LJRound,
		},
		{
			name:      "all zero values inherit all from b0",
			border:    &Border{Width: 0, col: nil, style: types.LJMiter},
			b0:        &Border{Width: 8, col: &color.Black, style: types.LJRound},
			wantWidth: 8,
			wantColor: &color.Black,
			wantStyle: types.LJRound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.border.mergeIn(tt.b0)

			if tt.border.Width != tt.wantWidth {
				t.Errorf("Border.mergeIn() Width = %d, want %d", tt.border.Width, tt.wantWidth)
			}

			if tt.border.col == nil && tt.wantColor != nil {
				t.Error("Border.mergeIn() col is nil, expected non-nil")
			} else if tt.border.col != nil && tt.wantColor != nil {
				if tt.border.col.R != tt.wantColor.R || tt.border.col.G != tt.wantColor.G || tt.border.col.B != tt.wantColor.B {
					t.Errorf("Border.mergeIn() col = (%f, %f, %f), want (%f, %f, %f)",
						tt.border.col.R, tt.border.col.G, tt.border.col.B,
						tt.wantColor.R, tt.wantColor.G, tt.wantColor.B)
				}
			}

			if tt.border.style != tt.wantStyle {
				t.Errorf("Border.mergeIn() style = %v, want %v", tt.border.style, tt.wantStyle)
			}
		})
	}
}
