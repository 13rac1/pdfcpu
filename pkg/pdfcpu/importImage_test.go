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

package pdfcpu

import (
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

func TestParsePageDim(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		unit      types.DisplayUnit
		wantWidth float64
		wantHeight float64
		wantErr   bool
	}{
		// Valid dimensions
		{
			name:       "points - integer values",
			input:      "100 200",
			unit:       types.POINTS,
			wantWidth:  100.0,
			wantHeight: 200.0,
			wantErr:    false,
		},
		{
			name:       "points - decimal values",
			input:      "100.5 200.75",
			unit:       types.POINTS,
			wantWidth:  100.5,
			wantHeight: 200.75,
			wantErr:    false,
		},
		{
			name:       "inches - converts to points",
			input:      "1 2",
			unit:       types.INCHES,
			wantWidth:  72.0,  // 1 inch = 72 points
			wantHeight: 144.0, // 2 inches = 144 points
			wantErr:    false,
		},
		{
			name:       "centimeters - converts to points",
			input:      "1 1",
			unit:       types.CENTIMETRES,
			wantWidth:  28.346456692913385, // 1 cm in points
			wantHeight: 28.346456692913385,
			wantErr:    false,
		},
		{
			name:       "millimeters - converts to points",
			input:      "10 20",
			unit:       types.MILLIMETRES,
			wantWidth:  28.346456692913385, // 10 mm in points
			wantHeight: 56.69291338582677,  // 20 mm in points
			wantErr:    false,
		},

		// Error cases
		{
			name:    "missing second dimension",
			input:   "100",
			unit:    types.POINTS,
			wantErr: true,
		},
		{
			name:    "too many dimensions",
			input:   "100 200 300",
			unit:    types.POINTS,
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			unit:    types.POINTS,
			wantErr: true,
		},
		{
			name:    "non-numeric width",
			input:   "abc 200",
			unit:    types.POINTS,
			wantErr: true,
		},
		{
			name:    "non-numeric height",
			input:   "100 xyz",
			unit:    types.POINTS,
			wantErr: true,
		},
		{
			name:    "zero width",
			input:   "0 200",
			unit:    types.POINTS,
			wantErr: true,
		},
		{
			name:    "zero height",
			input:   "100 0",
			unit:    types.POINTS,
			wantErr: true,
		},
		{
			name:    "negative width",
			input:   "-100 200",
			unit:    types.POINTS,
			wantErr: true,
		},
		{
			name:    "negative height",
			input:   "100 -200",
			unit:    types.POINTS,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dim, _, err := ParsePageDim(tt.input, tt.unit)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePageDim(%q, %v) error = %v, wantErr %v", tt.input, tt.unit, err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if dim == nil {
					t.Error("ParsePageDim() returned nil dimension for valid input")
					return
				}

				// Use approximate comparison for floating point
				const epsilon = 0.001
				if abs(dim.Width-tt.wantWidth) > epsilon {
					t.Errorf("ParsePageDim() Width = %v, want %v", dim.Width, tt.wantWidth)
				}
				if abs(dim.Height-tt.wantHeight) > epsilon {
					t.Errorf("ParsePageDim() Height = %v, want %v", dim.Height, tt.wantHeight)
				}
			}
		})
	}
}
