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

package model

import (
	"testing"
)

func TestPDFVersion(t *testing.T) {
	tests := []struct {
		input   string
		want    Version
		wantErr bool
	}{
		{"1.0", V10, false},
		{"1.1", V11, false},
		{"1.2", V12, false},
		{"1.3", V13, false},
		{"1.4", V14, false},
		{"1.5", V15, false},
		{"1.6", V16, false},
		{"1.7", V17, false},
		{"2.0", V20, false},
		{"", -1, true},
		{"0.9", -1, true},
		{"1.8", -1, true},
		{"2.1", -1, true},
		{"invalid", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := PDFVersion(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("PDFVersion(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PDFVersion(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestPDFVersionRelaxed(t *testing.T) {
	tests := []struct {
		input   string
		want    Version
		wantErr bool
	}{
		{"1.7.0", V17, false},
		{"1.7", -1, true},
		{"1.7.1", -1, true},
		{"invalid", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := PDFVersionRelaxed(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("PDFVersionRelaxed(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PDFVersionRelaxed(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestVersionString(t *testing.T) {
	tests := []struct {
		v    Version
		want string
	}{
		{V10, "1.0"},
		{V11, "1.1"},
		{V12, "1.2"},
		{V13, "1.3"},
		{V14, "1.4"},
		{V15, "1.5"},
		{V16, "1.6"},
		{V17, "1.7"},
		{V20, "2.0"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.v.String()
			if got != tt.want {
				t.Errorf("Version(%d).String() = %q, want %q", tt.v, got, tt.want)
			}
		})
	}
}

func TestIdenticalMajorAndMinorVersions(t *testing.T) {
	tests := []struct {
		name string
		v1   string
		v2   string
		want bool
	}{
		{"same versions", "1.2.3", "1.2.4", true},
		{"different minor", "1.2.3", "1.3.3", false},
		{"different major", "1.2.3", "2.2.3", false},
		{"no dot v1", "123", "1.2.3", false},
		{"no dot v2", "1.2.3", "123", false},
		{"single dot v1", "1.2", "1.2.3", true},
		{"single dot v2", "1.2.3", "1.2", true},
		{"empty v1", "", "1.2.3", false},
		{"empty v2", "1.2.3", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := identicalMajorAndMinorVersions(tt.v1, tt.v2)
			if got != tt.want {
				t.Errorf("identicalMajorAndMinorVersions(%q, %q) = %v, want %v", tt.v1, tt.v2, got, tt.want)
			}
		})
	}
}

func TestVersionConstants(t *testing.T) {
	// Verify version constants have expected values
	if V10 != 0 {
		t.Errorf("V10 = %d, want 0", V10)
	}
	if V17 != 7 {
		t.Errorf("V17 = %d, want 7", V17)
	}
	if V20 != 8 {
		t.Errorf("V20 = %d, want 8", V20)
	}
}

func TestVersionStr(t *testing.T) {
	// VersionStr should be set
	if VersionStr == "" {
		t.Error("VersionStr should not be empty")
	}
}
