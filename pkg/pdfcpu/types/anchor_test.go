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

package types

import "testing"

func TestAnchorString(t *testing.T) {
	tests := []struct {
		anchor Anchor
		want   string
	}{
		{TopLeft, "top left"},
		{TopCenter, "top center"},
		{TopRight, "top right"},
		{Left, "left"},
		{Center, "center"},
		{Right, "right"},
		{BottomLeft, "bottom left"},
		{BottomCenter, "bottom center"},
		{BottomRight, "bottom right"},
		{Full, "full"},
		{Anchor(-1), ""}, // Unknown anchor
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.anchor.String()
			if got != tt.want {
				t.Errorf("Anchor(%d).String() = %q, want %q", tt.anchor, got, tt.want)
			}
		})
	}
}

func TestAnchorConstants(t *testing.T) {
	// Verify anchor constants have expected values
	if TopLeft != 0 {
		t.Errorf("TopLeft = %d, want 0", TopLeft)
	}
	if Center != 4 {
		t.Errorf("Center = %d, want 4", Center)
	}
	if Full != 9 {
		t.Errorf("Full = %d, want 9", Full)
	}
}
