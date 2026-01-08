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

func TestGuideValidate(t *testing.T) {
	tests := []struct {
		name  string
		guide *Guide
		wantX float64
		wantY float64
	}{
		{
			name:  "positive coordinates",
			guide: &Guide{Position: [2]float64{100.0, 200.0}},
			wantX: 100.0,
			wantY: 200.0,
		},
		{
			name:  "zero coordinates",
			guide: &Guide{Position: [2]float64{0, 0}},
			wantX: 0,
			wantY: 0,
		},
		{
			name:  "negative coordinates",
			guide: &Guide{Position: [2]float64{-50.0, -75.0}},
			wantX: -50.0,
			wantY: -75.0,
		},
		{
			name:  "mixed positive and negative",
			guide: &Guide{Position: [2]float64{100.0, -50.0}},
			wantX: 100.0,
			wantY: -50.0,
		},
		{
			name:  "fractional coordinates",
			guide: &Guide{Position: [2]float64{12.5, 37.8}},
			wantX: 12.5,
			wantY: 37.8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.guide.validate()

			if tt.guide.x != tt.wantX {
				t.Errorf("Guide.validate() x = %.1f, want %.1f", tt.guide.x, tt.wantX)
			}
			if tt.guide.y != tt.wantY {
				t.Errorf("Guide.validate() y = %.1f, want %.1f", tt.guide.y, tt.wantY)
			}
		})
	}
}
