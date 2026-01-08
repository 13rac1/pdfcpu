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

func TestParsePageConfiguration(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		unit         types.DisplayUnit
		wantPageSize string
		wantUserDim  bool
		wantErr      bool
	}{
		// Empty input
		{
			name:    "empty string returns nil",
			input:   "",
			unit:    types.POINTS,
			wantErr: false,
		},

		// Valid formsize/papersize formats
		{
			name:         "formsize A4",
			input:        "formsize:A4",
			unit:         types.POINTS,
			wantPageSize: "A4",
			wantUserDim:  true,
			wantErr:      false,
		},
		{
			name:         "papersize Letter",
			input:        "papersize:Letter",
			unit:         types.POINTS,
			wantPageSize: "Letter",
			wantUserDim:  true,
			wantErr:      false,
		},
		{
			name:         "papersize Legal",
			input:        "papersize:Legal",
			unit:         types.POINTS,
			wantPageSize: "Legal",
			wantUserDim:  true,
			wantErr:      false,
		},

		// Valid dimensions format
		{
			name:         "dimensions with points",
			input:        "dimensions:100 200",
			unit:         types.POINTS,
			wantPageSize: "",
			wantUserDim:  true,
			wantErr:      false,
		},
		{
			name:         "dimensions with inches",
			input:        "dimensions:8.5 11",
			unit:         types.INCHES,
			wantPageSize: "",
			wantUserDim:  true,
			wantErr:      false,
		},

		// Parameter prefix completion
		{
			name:         "form prefix (abbreviation)",
			input:        "form:A4",
			unit:         types.POINTS,
			wantPageSize: "A4",
			wantUserDim:  true,
			wantErr:      false,
		},
		{
			name:         "paper prefix (abbreviation)",
			input:        "paper:Letter",
			unit:         types.POINTS,
			wantPageSize: "Letter",
			wantUserDim:  true,
			wantErr:      false,
		},
		{
			name:         "dim prefix (abbreviation)",
			input:        "dim:100 200",
			unit:         types.POINTS,
			wantPageSize: "",
			wantUserDim:  true,
			wantErr:      false,
		},

		// Error cases
		{
			name:    "invalid format - no colon",
			input:   "formsize A4",
			unit:    types.POINTS,
			wantErr: true,
		},
		{
			name:    "invalid format - multiple colons",
			input:   "formsize:A4:extra",
			unit:    types.POINTS,
			wantErr: true,
		},
		{
			name:    "unknown parameter",
			input:   "unknown:value",
			unit:    types.POINTS,
			wantErr: true,
		},
		{
			name:         "single letter prefix f matches formsize",
			input:        "f:A4",
			unit:         types.POINTS,
			wantPageSize: "A4",
			wantUserDim:  true,
			wantErr:      false,
		},
		{
			name:    "both formsize and dimensions",
			input:   "formsize:A4,dimensions:100 200",
			unit:    types.POINTS,
			wantErr: true, // Only one of formsize or dimensions allowed
		},
		{
			name:    "invalid dimensions",
			input:   "dimensions:abc def",
			unit:    types.POINTS,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := ParsePageConfiguration(tt.input, tt.unit)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePageConfiguration(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Empty input should return nil
			if tt.input == "" {
				if cfg != nil {
					t.Error("ParsePageConfiguration(\"\") should return nil")
				}
				return
			}

			if cfg == nil {
				t.Fatal("ParsePageConfiguration() returned nil for valid input")
			}

			if cfg.UserDim != tt.wantUserDim {
				t.Errorf("ParsePageConfiguration() UserDim = %v, want %v", cfg.UserDim, tt.wantUserDim)
			}

			if tt.wantPageSize != "" && cfg.PageSize != tt.wantPageSize {
				t.Errorf("ParsePageConfiguration() PageSize = %q, want %q", cfg.PageSize, tt.wantPageSize)
			}

			if cfg.PageDim == nil {
				t.Error("ParsePageConfiguration() PageDim is nil")
			}

			if cfg.InpUnit != tt.unit {
				t.Errorf("ParsePageConfiguration() InpUnit = %v, want %v", cfg.InpUnit, tt.unit)
			}
		})
	}
}
