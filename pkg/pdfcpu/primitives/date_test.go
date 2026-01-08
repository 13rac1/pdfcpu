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
	"strings"
	"testing"
)

func TestDateFormatForFmtInt(t *testing.T) {
	tests := []struct {
		name    string
		fmtInt  string
		wantExt string
		wantErr bool
	}{
		// Dash separator formats
		{"dash yyyy-mm-dd", "2006-01-02", "yyyy-mm-dd", false},
		{"dash yyyy-dd-mm", "2006-02-01", "yyyy-dd-mm", false},
		{"dash dd-mm-yyyy", "02-01-2006", "dd-mm-yyyy", false},
		{"dash mm-dd-yyyy", "01-02-2006", "mm-dd-yyyy", false},
		{"dash yyyy-m-d", "2006-1-2", "yyyy-m-d", false},
		{"dash yyyy-d-m", "2006-2-1", "yyyy-d-m", false},
		{"dash d-m-yyyy", "2-1-2006", "d-m-yyyy", false},
		{"dash m-d-yyyy", "1-2-2006", "m-d-yyyy", false},

		// Slash separator formats
		{"slash yyyy/mm/dd", "2006/01/02", "yyyy/mm/dd", false},
		{"slash yyyy/dd/mm", "2006/02/01", "yyyy/dd/mm", false},
		{"slash dd/mm/yyyy", "02/01/2006", "dd/mm/yyyy", false},
		{"slash mm/dd/yyyy", "01/02/2006", "mm/dd/yyyy", false},
		{"slash yyyy/m/d", "2006/1/2", "yyyy/m/d", false},
		{"slash yyyy/d/m", "2006/2/1", "yyyy/d/m", false},
		{"slash d/m/yyyy", "2/1/2006", "d/m/yyyy", false},
		{"slash m/d/yyyy", "1/2/2006", "m/d/yyyy", false},

		// Dot separator formats
		{"dot yyyy.mm.dd", "2006.01.02", "yyyy.mm.dd", false},
		{"dot yyyy.dd.mm", "2006.02.01", "yyyy.dd.mm", false},
		{"dot dd.mm.yyyy", "02.01.2006", "dd.mm.yyyy", false},
		{"dot mm.dd.yyyy", "01.02.2006", "mm.dd.yyyy", false},
		{"dot yyyy.m.d", "2006.1.2", "yyyy.m.d", false},
		{"dot yyyy.d.m", "2006.2.1", "yyyy.d.m", false},
		{"dot d.m.yyyy", "2.1.2006", "d.m.yyyy", false},
		{"dot m.d.yyyy", "1.2.2006", "m.d.yyyy", false},

		// Invalid formats
		{"invalid format", "invalid", "", true},
		{"empty string", "", "", true},
		{"wrong separator", "2006_01_02", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			df, err := DateFormatForFmtInt(tt.fmtInt)
			if (err != nil) != tt.wantErr {
				t.Errorf("DateFormatForFmtInt(%q) error = %v, wantErr %v", tt.fmtInt, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if df == nil {
					t.Error("DateFormatForFmtInt() returned nil DateFormat")
					return
				}
				if df.Ext != tt.wantExt {
					t.Errorf("DateFormatForFmtInt(%q) Ext = %q, want %q", tt.fmtInt, df.Ext, tt.wantExt)
				}
				if df.Int != tt.fmtInt {
					t.Errorf("DateFormatForFmtInt(%q) Int = %q, want %q", tt.fmtInt, df.Int, tt.fmtInt)
				}
			}
		})
	}
}

func TestDateFormatForFmtExt(t *testing.T) {
	tests := []struct {
		name    string
		fmtExt  string
		wantInt string
		wantErr bool
	}{
		// Lowercase formats
		{"lowercase yyyy-mm-dd", "yyyy-mm-dd", "2006-01-02", false},
		{"lowercase dd/mm/yyyy", "dd/mm/yyyy", "02/01/2006", false},
		{"lowercase mm.dd.yyyy", "mm.dd.yyyy", "01.02.2006", false},

		// Uppercase formats (should work - case insensitive)
		{"uppercase YYYY-MM-DD", "YYYY-MM-DD", "2006-01-02", false},
		{"uppercase DD/MM/YYYY", "DD/MM/YYYY", "02/01/2006", false},
		{"uppercase MM.DD.YYYY", "MM.DD.YYYY", "01.02.2006", false},

		// Mixed case formats
		{"mixed Yyyy-Mm-Dd", "Yyyy-Mm-Dd", "2006-01-02", false},
		{"mixed Dd/Mm/Yyyy", "Dd/Mm/Yyyy", "02/01/2006", false},

		// All 24 formats (lowercase)
		{"dash yyyy-m-d", "yyyy-m-d", "2006-1-2", false},
		{"dash yyyy-d-m", "yyyy-d-m", "2006-2-1", false},
		{"dash d-m-yyyy", "d-m-yyyy", "2-1-2006", false},
		{"dash m-d-yyyy", "m-d-yyyy", "1-2-2006", false},
		{"slash yyyy/m/d", "yyyy/m/d", "2006/1/2", false},
		{"slash yyyy/d/m", "yyyy/d/m", "2006/2/1", false},
		{"slash d/m/yyyy", "d/m/yyyy", "2/1/2006", false},
		{"slash m/d/yyyy", "m/d/yyyy", "1/2/2006", false},
		{"dot yyyy.m.d", "yyyy.m.d", "2006.1.2", false},
		{"dot yyyy.d.m", "yyyy.d.m", "2006.2.1", false},
		{"dot d.m.yyyy", "d.m.yyyy", "2.1.2006", false},
		{"dot m.d.yyyy", "m.d.yyyy", "1.2.2006", false},

		// Invalid formats
		{"invalid format", "invalid", "", true},
		{"empty string", "", "", true},
		{"wrong separator", "yyyy_mm_dd", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			df, err := DateFormatForFmtExt(tt.fmtExt)
			if (err != nil) != tt.wantErr {
				t.Errorf("DateFormatForFmtExt(%q) error = %v, wantErr %v", tt.fmtExt, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if df == nil {
					t.Error("DateFormatForFmtExt() returned nil DateFormat")
					return
				}
				if df.Int != tt.wantInt {
					t.Errorf("DateFormatForFmtExt(%q) Int = %q, want %q", tt.fmtExt, df.Int, tt.wantInt)
				}
			}
		})
	}
}

func TestDateFormatForDate(t *testing.T) {
	tests := []struct {
		name    string
		date    string
		wantErr bool
	}{
		// Dash separator dates
		{"dash 2023-12-25", "2023-12-25", false},
		{"dash 25-12-2023", "25-12-2023", false},
		{"dash 12-25-2023", "12-25-2023", false},
		{"dash 2023-1-5", "2023-1-5", false},
		{"dash 5-1-2023", "5-1-2023", false},

		// Slash separator dates
		{"slash 2023/12/25", "2023/12/25", false},
		{"slash 25/12/2023", "25/12/2023", false},
		{"slash 12/25/2023", "12/25/2023", false},
		{"slash 2023/1/5", "2023/1/5", false},
		{"slash 5/1/2023", "5/1/2023", false},

		// Dot separator dates
		{"dot 2023.12.25", "2023.12.25", false},
		{"dot 25.12.2023", "25.12.2023", false},
		{"dot 12.25.2023", "12.25.2023", false},
		{"dot 2023.1.5", "2023.1.5", false},
		{"dot 5.1.2023", "5.1.2023", false},

		// Invalid dates
		{"invalid format", "invalid", true},
		{"empty string", "", true},
		{"wrong separator", "2023_12_25", true},
		{"invalid date", "2023-13-45", true},
		{"nonsense", "abc/def/ghi", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			df, err := DateFormatForDate(tt.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("DateFormatForDate(%q) error = %v, wantErr %v", tt.date, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if df == nil {
					t.Error("DateFormatForDate() returned nil DateFormat")
					return
				}
				// Verify the returned format can parse the date
				if err := df.validate(tt.date); err != nil {
					t.Errorf("DateFormatForDate(%q) returned format %q that cannot parse the date: %v",
						tt.date, df.Int, err)
				}
			}
		})
	}
}

func TestDateFormatValidate(t *testing.T) {
	tests := []struct {
		name    string
		format  DateFormat
		date    string
		wantErr bool
	}{
		// Valid dates with different formats
		{
			name:    "valid yyyy-mm-dd",
			format:  DateFormat{Int: "2006-01-02", Ext: "yyyy-mm-dd"},
			date:    "2023-12-25",
			wantErr: false,
		},
		{
			name:    "valid dd-mm-yyyy",
			format:  DateFormat{Int: "02-01-2006", Ext: "dd-mm-yyyy"},
			date:    "25-12-2023",
			wantErr: false,
		},
		{
			name:    "valid mm/dd/yyyy",
			format:  DateFormat{Int: "01/02/2006", Ext: "mm/dd/yyyy"},
			date:    "12/25/2023",
			wantErr: false,
		},
		{
			name:    "valid dd.mm.yyyy",
			format:  DateFormat{Int: "02.01.2006", Ext: "dd.mm.yyyy"},
			date:    "25.12.2023",
			wantErr: false,
		},
		{
			name:    "valid yyyy/m/d",
			format:  DateFormat{Int: "2006/1/2", Ext: "yyyy/m/d"},
			date:    "2023/1/5",
			wantErr: false,
		},
		{
			name:    "valid leap year date",
			format:  DateFormat{Int: "2006-01-02", Ext: "yyyy-mm-dd"},
			date:    "2020-02-29",
			wantErr: false,
		},

		// Invalid dates - wrong format
		{
			name:    "wrong format - expected dash, got slash",
			format:  DateFormat{Int: "2006-01-02", Ext: "yyyy-mm-dd"},
			date:    "2023/12/25",
			wantErr: true,
		},
		{
			name:    "wrong format - expected mm/dd, got dd/mm",
			format:  DateFormat{Int: "01/02/2006", Ext: "mm/dd/yyyy"},
			date:    "25/12/2023",
			wantErr: true,
		},

		// Invalid dates - invalid values
		{
			name:    "invalid month 13",
			format:  DateFormat{Int: "2006-01-02", Ext: "yyyy-mm-dd"},
			date:    "2023-13-25",
			wantErr: true,
		},
		{
			name:    "invalid day 32",
			format:  DateFormat{Int: "2006-01-02", Ext: "yyyy-mm-dd"},
			date:    "2023-12-32",
			wantErr: true,
		},
		{
			name:    "invalid leap year",
			format:  DateFormat{Int: "2006-01-02", Ext: "yyyy-mm-dd"},
			date:    "2023-02-29",
			wantErr: true,
		},
		{
			name:    "empty date string",
			format:  DateFormat{Int: "2006-01-02", Ext: "yyyy-mm-dd"},
			date:    "",
			wantErr: true,
		},
		{
			name:    "nonsense date",
			format:  DateFormat{Int: "2006-01-02", Ext: "yyyy-mm-dd"},
			date:    "abc-def-ghi",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.format.validate(tt.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("DateFormat.validate(%q) error = %v, wantErr %v", tt.date, err, tt.wantErr)
			}
		})
	}
}

func TestDateFormatsCount(t *testing.T) {
	// Verify we have all 24 expected formats
	if len(dateFormats) != 24 {
		t.Errorf("dateFormats length = %d, want 24 (8 patterns × 3 separators)", len(dateFormats))
	}

	// Verify we have 8 formats per separator
	dashCount := 0
	slashCount := 0
	dotCount := 0

	for _, df := range dateFormats {
		if strings.Contains(df.Ext, "-") {
			dashCount++
		} else if strings.Contains(df.Ext, "/") {
			slashCount++
		} else if strings.Contains(df.Ext, ".") {
			dotCount++
		}
	}

	if dashCount != 8 {
		t.Errorf("dash separator formats = %d, want 8", dashCount)
	}
	if slashCount != 8 {
		t.Errorf("slash separator formats = %d, want 8", slashCount)
	}
	if dotCount != 8 {
		t.Errorf("dot separator formats = %d, want 8", dotCount)
	}
}
