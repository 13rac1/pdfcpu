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

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func TestParseOrientation(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		checkValue func(*model.NUp) bool
		wantErr    bool
	}{
		{
			name:  "right down orientation",
			input: "rd",
			checkValue: func(nup *model.NUp) bool {
				return nup.Orient == model.RightDown
			},
			wantErr: false,
		},
		{
			name:  "down right orientation",
			input: "dr",
			checkValue: func(nup *model.NUp) bool {
				return nup.Orient == model.DownRight
			},
			wantErr: false,
		},
		{
			name:  "left down orientation",
			input: "ld",
			checkValue: func(nup *model.NUp) bool {
				return nup.Orient == model.LeftDown
			},
			wantErr: false,
		},
		{
			name:  "down left orientation",
			input: "dl",
			checkValue: func(nup *model.NUp) bool {
				return nup.Orient == model.DownLeft
			},
			wantErr: false,
		},
		{
			name:    "invalid orientation",
			input:   "invalid",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "uppercase",
			input:   "RD",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nup := model.DefaultNUpConfig()
			err := parseOrientation(tt.input, nup)

			if (err != nil) != tt.wantErr {
				t.Errorf("parseOrientation(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.checkValue != nil && !tt.checkValue(nup) {
				t.Errorf("parseOrientation(%q) Orient = %v, want correct value", tt.input, nup.Orient)
			}
		})
	}
}

func TestParseEnforce(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantEnforce bool
		wantErr     bool
	}{
		// True variants
		{
			name:        "on",
			input:       "on",
			wantEnforce: true,
			wantErr:     false,
		},
		{
			name:        "true",
			input:       "true",
			wantEnforce: true,
			wantErr:     false,
		},
		{
			name:        "t",
			input:       "t",
			wantEnforce: true,
			wantErr:     false,
		},
		{
			name:        "ON (uppercase)",
			input:       "ON",
			wantEnforce: true,
			wantErr:     false,
		},
		{
			name:        "True (mixed case)",
			input:       "True",
			wantEnforce: true,
			wantErr:     false,
		},
		// False variants
		{
			name:        "off",
			input:       "off",
			wantEnforce: false,
			wantErr:     false,
		},
		{
			name:        "false",
			input:       "false",
			wantEnforce: false,
			wantErr:     false,
		},
		{
			name:        "f",
			input:       "f",
			wantEnforce: false,
			wantErr:     false,
		},
		{
			name:        "OFF (uppercase)",
			input:       "OFF",
			wantEnforce: false,
			wantErr:     false,
		},
		// Invalid
		{
			name:    "invalid",
			input:   "invalid",
			wantErr: true,
		},
		{
			name:    "empty",
			input:   "",
			wantErr: true,
		},
		{
			name:    "yes",
			input:   "yes",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nup := model.DefaultNUpConfig()
			err := parseEnforce(tt.input, nup)

			if (err != nil) != tt.wantErr {
				t.Errorf("parseEnforce(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}

			if !tt.wantErr && nup.Enforce != tt.wantEnforce {
				t.Errorf("parseEnforce(%q) Enforce = %v, want %v", tt.input, nup.Enforce, tt.wantEnforce)
			}
		})
	}
}

func TestParseElementBorder(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantBorder bool
		wantErr    bool
	}{
		// True variants
		{
			name:       "on",
			input:      "on",
			wantBorder: true,
			wantErr:    false,
		},
		{
			name:       "true",
			input:      "true",
			wantBorder: true,
			wantErr:    false,
		},
		{
			name:       "t",
			input:      "t",
			wantBorder: true,
			wantErr:    false,
		},
		// False variants
		{
			name:       "off",
			input:      "off",
			wantBorder: false,
			wantErr:    false,
		},
		{
			name:       "false",
			input:      "false",
			wantBorder: false,
			wantErr:    false,
		},
		{
			name:       "f",
			input:      "f",
			wantBorder: false,
			wantErr:    false,
		},
		// Case insensitive
		{
			name:       "ON (uppercase)",
			input:      "ON",
			wantBorder: true,
			wantErr:    false,
		},
		{
			name:       "OFF (uppercase)",
			input:      "OFF",
			wantBorder: false,
			wantErr:    false,
		},
		// Invalid
		{
			name:    "invalid",
			input:   "invalid",
			wantErr: true,
		},
		{
			name:    "empty",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nup := model.DefaultNUpConfig()
			err := parseElementBorder(tt.input, nup)

			if (err != nil) != tt.wantErr {
				t.Errorf("parseElementBorder(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}

			if !tt.wantErr && nup.Border != tt.wantBorder {
				t.Errorf("parseElementBorder(%q) Border = %v, want %v", tt.input, nup.Border, tt.wantBorder)
			}
		})
	}
}

func TestParseBookletGuides(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantGuides bool
		wantErr    bool
	}{
		// True variants
		{
			name:       "on",
			input:      "on",
			wantGuides: true,
			wantErr:    false,
		},
		{
			name:       "true",
			input:      "true",
			wantGuides: true,
			wantErr:    false,
		},
		{
			name:       "t",
			input:      "t",
			wantGuides: true,
			wantErr:    false,
		},
		// False variants
		{
			name:       "off",
			input:      "off",
			wantGuides: false,
			wantErr:    false,
		},
		{
			name:       "false",
			input:      "false",
			wantGuides: false,
			wantErr:    false,
		},
		{
			name:       "f",
			input:      "f",
			wantGuides: false,
			wantErr:    false,
		},
		// Case insensitive
		{
			name:       "ON (uppercase)",
			input:      "ON",
			wantGuides: true,
			wantErr:    false,
		},
		// Invalid
		{
			name:    "invalid",
			input:   "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nup := model.DefaultNUpConfig()
			err := parseBookletGuides(tt.input, nup)

			if (err != nil) != tt.wantErr {
				t.Errorf("parseBookletGuides(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}

			if !tt.wantErr && nup.BookletGuides != tt.wantGuides {
				t.Errorf("parseBookletGuides(%q) BookletGuides = %v, want %v", tt.input, nup.BookletGuides, tt.wantGuides)
			}
		})
	}
}

func TestParseBookletMultifolio(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		wantMultiFolio bool
		wantErr        bool
	}{
		// True variants
		{
			name:           "on",
			input:          "on",
			wantMultiFolio: true,
			wantErr:        false,
		},
		{
			name:           "true",
			input:          "true",
			wantMultiFolio: true,
			wantErr:        false,
		},
		{
			name:           "t",
			input:          "t",
			wantMultiFolio: true,
			wantErr:        false,
		},
		// False variants
		{
			name:           "off",
			input:          "off",
			wantMultiFolio: false,
			wantErr:        false,
		},
		{
			name:           "false",
			input:          "false",
			wantMultiFolio: false,
			wantErr:        false,
		},
		{
			name:           "f",
			input:          "f",
			wantMultiFolio: false,
			wantErr:        false,
		},
		// Case insensitive
		{
			name:           "TRUE (uppercase)",
			input:          "TRUE",
			wantMultiFolio: true,
			wantErr:        false,
		},
		// Invalid
		{
			name:    "invalid",
			input:   "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nup := model.DefaultNUpConfig()
			err := parseBookletMultifolio(tt.input, nup)

			if (err != nil) != tt.wantErr {
				t.Errorf("parseBookletMultifolio(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}

			if !tt.wantErr && nup.MultiFolio != tt.wantMultiFolio {
				t.Errorf("parseBookletMultifolio(%q) MultiFolio = %v, want %v", tt.input, nup.MultiFolio, tt.wantMultiFolio)
			}
		})
	}
}

func TestParseBookletFolioSize(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		wantFolioSize int
		wantErr       bool
	}{
		{
			name:          "valid size 8",
			input:         "8",
			wantFolioSize: 8,
			wantErr:       false,
		},
		{
			name:          "valid size 4",
			input:         "4",
			wantFolioSize: 4,
			wantErr:       false,
		},
		{
			name:          "valid size 16",
			input:         "16",
			wantFolioSize: 16,
			wantErr:       false,
		},
		{
			name:    "non-numeric",
			input:   "abc",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "decimal not allowed",
			input:   "8.5",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nup := model.DefaultNUpConfig()
			err := parseBookletFolioSize(tt.input, nup)

			if (err != nil) != tt.wantErr {
				t.Errorf("parseBookletFolioSize(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}

			if !tt.wantErr && nup.FolioSize != tt.wantFolioSize {
				t.Errorf("parseBookletFolioSize(%q) FolioSize = %v, want %v", tt.input, nup.FolioSize, tt.wantFolioSize)
			}
		})
	}
}
