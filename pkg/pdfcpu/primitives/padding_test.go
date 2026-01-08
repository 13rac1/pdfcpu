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
)

func TestPaddingValidate(t *testing.T) {
	tests := []struct {
		name    string
		padding *Padding
		wantErr bool
		errMsg  string
	}{
		{
			name:    "invalid $ reference",
			padding: &Padding{Name: "$"},
			wantErr: true,
			errMsg:  "invalid padding reference $",
		},
		{
			name: "positive width expands to all sides",
			padding: &Padding{
				Width: 10.0,
			},
			wantErr: false,
		},
		{
			name: "negative width with individual paddings error",
			padding: &Padding{
				Width:  -1.0,
				Top:    5.0,
				Right:  5.0,
				Bottom: 5.0,
				Left:   5.0,
			},
			wantErr: true,
			errMsg:  "invalid padding width",
		},
		{
			name: "negative width with zero individual paddings ok",
			padding: &Padding{
				Width:  -1.0,
				Top:    0,
				Right:  0,
				Bottom: 0,
				Left:   0,
			},
			wantErr: false,
		},
		{
			name: "zero width with individual paddings ok",
			padding: &Padding{
				Width:  0,
				Top:    5.0,
				Right:  10.0,
				Bottom: 15.0,
				Left:   20.0,
			},
			wantErr: false,
		},
		{
			name: "all zeros valid",
			padding: &Padding{
				Width:  0,
				Top:    0,
				Right:  0,
				Bottom: 0,
				Left:   0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.padding.validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Padding.validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" {
				if err.Error() != tt.errMsg && !contains(err.Error(), tt.errMsg) {
					t.Errorf("Padding.validate() error = %q, want to contain %q", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestPaddingValidateWidthExpansion(t *testing.T) {
	p := &Padding{
		Name:  "test",
		Width: 10.0,
	}

	err := p.validate()
	if err != nil {
		t.Fatalf("Padding.validate() unexpected error = %v", err)
	}

	// After validation, Width > 0 should expand to all sides
	if p.Top != 10.0 || p.Right != 10.0 || p.Bottom != 10.0 || p.Left != 10.0 {
		t.Errorf("Padding Width expansion: got (%.1f, %.1f, %.1f, %.1f), want (10.0, 10.0, 10.0, 10.0)",
			p.Top, p.Right, p.Bottom, p.Left)
	}
}

func TestPaddingMergeIn(t *testing.T) {
	tests := []struct {
		name       string
		padding    *Padding
		p0         *Padding
		wantTop    float64
		wantRight  float64
		wantBottom float64
		wantLeft   float64
	}{
		{
			name: "positive width - no merge",
			padding: &Padding{
				Width: 15.0,
				Top:   1.0,
				Right: 2.0,
			},
			p0: &Padding{
				Top:    10.0,
				Right:  10.0,
				Bottom: 10.0,
				Left:   10.0,
			},
			wantTop:    1.0,
			wantRight:  2.0,
			wantBottom: 0,
			wantLeft:   0,
		},
		{
			name: "negative width - reset all to zero",
			padding: &Padding{
				Width:  -1.0,
				Top:    5.0,
				Right:  5.0,
				Bottom: 5.0,
				Left:   5.0,
			},
			p0: &Padding{
				Top:    10.0,
				Right:  10.0,
				Bottom: 10.0,
				Left:   10.0,
			},
			wantTop:    0,
			wantRight:  0,
			wantBottom: 0,
			wantLeft:   0,
		},
		{
			name: "zero values inherit from p0",
			padding: &Padding{
				Width:  0,
				Top:    0,
				Right:  20.0,
				Bottom: 0,
				Left:   30.0,
			},
			p0: &Padding{
				Top:    10.0,
				Right:  15.0,
				Bottom: 12.0,
				Left:   18.0,
			},
			wantTop:    10.0, // inherited from p0
			wantRight:  20.0, // kept from padding
			wantBottom: 12.0, // inherited from p0
			wantLeft:   30.0, // kept from padding
		},
		{
			name: "negative values reset to zero",
			padding: &Padding{
				Width:  0,
				Top:    -5.0,
				Right:  -10.0,
				Bottom: 8.0,
				Left:   -2.0,
			},
			p0: &Padding{
				Top:    10.0,
				Right:  10.0,
				Bottom: 10.0,
				Left:   10.0,
			},
			wantTop:    0,   // reset from negative
			wantRight:  0,   // reset from negative
			wantBottom: 8.0, // kept from padding
			wantLeft:   0,   // reset from negative
		},
		{
			name: "all zero padding inherits all from p0",
			padding: &Padding{
				Width:  0,
				Top:    0,
				Right:  0,
				Bottom: 0,
				Left:   0,
			},
			p0: &Padding{
				Top:    5.0,
				Right:  6.0,
				Bottom: 7.0,
				Left:   8.0,
			},
			wantTop:    5.0,
			wantRight:  6.0,
			wantBottom: 7.0,
			wantLeft:   8.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.padding.mergeIn(tt.p0)

			if tt.padding.Top != tt.wantTop {
				t.Errorf("Top = %.1f, want %.1f", tt.padding.Top, tt.wantTop)
			}
			if tt.padding.Right != tt.wantRight {
				t.Errorf("Right = %.1f, want %.1f", tt.padding.Right, tt.wantRight)
			}
			if tt.padding.Bottom != tt.wantBottom {
				t.Errorf("Bottom = %.1f, want %.1f", tt.padding.Bottom, tt.wantBottom)
			}
			if tt.padding.Left != tt.wantLeft {
				t.Errorf("Left = %.1f, want %.1f", tt.padding.Left, tt.wantLeft)
			}
		})
	}
}
