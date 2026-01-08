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

func TestFormFontValidateISO639(t *testing.T) {
	tests := []struct {
		name    string
		lang    string
		wantErr bool
	}{
		// Valid ISO-639 codes (sample from different parts of the list)
		{"valid en", "en", false},
		{"valid zh", "zh", false},
		{"valid ar", "ar", false},
		{"valid fr", "fr", false},
		{"valid de", "de", false},
		{"valid ja", "ja", false},
		{"valid es", "es", false},
		{"valid ru", "ru", false},
		{"valid pt", "pt", false},
		{"valid it", "it", false},
		{"valid ko", "ko", false},
		{"valid hi", "hi", false},
		{"valid he", "he", false},
		{"valid fa", "fa", false},
		{"valid pl", "pl", false},
		{"valid tr", "tr", false},
		{"valid uk", "uk", false},
		{"valid vi", "vi", false},
		{"valid th", "th", false},
		{"valid sv", "sv", false},

		// Invalid ISO-639 codes
		{"invalid xx", "xx", true},
		{"invalid ZZ", "ZZ", true},
		{"invalid abc", "abc", true},
		{"invalid empty", "", true},
		{"invalid 12", "12", true},
		{"invalid EN (uppercase)", "EN", true}, // codes are lowercase
		{"invalid e", "e", true},              // too short
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FormFont{Lang: tt.lang}
			err := f.validateISO639()
			if (err != nil) != tt.wantErr {
				t.Errorf("FormFont.validateISO639() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.wantErr {
				// Verify error message contains the invalid code
				if !strings.Contains(err.Error(), tt.lang) {
					t.Errorf("FormFont.validateISO639() error = %v, should contain lang code %q", err, tt.lang)
				}
			}
		})
	}
}

func TestFormFontRTL(t *testing.T) {
	tests := []struct {
		name   string
		font   FormFont
		wantRTL bool
	}{
		// RTL scripts
		{"Arab script", FormFont{Script: "Arab"}, true},
		{"Hebr script", FormFont{Script: "Hebr"}, true},

		// RTL languages
		{"Arabic language (ar)", FormFont{Lang: "ar"}, true},
		{"Persian language (fa)", FormFont{Lang: "fa"}, true},
		{"Hebrew language (he)", FormFont{Lang: "he"}, true},

		// RTL script + language
		{"Arab script with ar language", FormFont{Script: "Arab", Lang: "ar"}, true},
		{"Hebr script with he language", FormFont{Script: "Hebr", Lang: "he"}, true},

		// Non-RTL scripts
		{"Latn script", FormFont{Script: "Latn"}, false},
		{"Cyrl script", FormFont{Script: "Cyrl"}, false},
		{"Grek script", FormFont{Script: "Grek"}, false},

		// Non-RTL languages
		{"English (en)", FormFont{Lang: "en"}, false},
		{"Chinese (zh)", FormFont{Lang: "zh"}, false},
		{"French (fr)", FormFont{Lang: "fr"}, false},
		{"German (de)", FormFont{Lang: "de"}, false},
		{"Japanese (ja)", FormFont{Lang: "ja"}, false},
		{"Spanish (es)", FormFont{Lang: "es"}, false},
		{"Russian (ru)", FormFont{Lang: "ru"}, false},

		// Empty values
		{"empty script and language", FormFont{}, false},
		{"empty script", FormFont{Lang: "en"}, false},
		{"empty language", FormFont{Script: "Latn"}, false},

		// Mixed RTL and non-RTL
		{"RTL script with non-RTL language", FormFont{Script: "Arab", Lang: "en"}, true},
		{"non-RTL script with RTL language", FormFont{Script: "Latn", Lang: "ar"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.font.RTL(); got != tt.wantRTL {
				t.Errorf("FormFont.RTL() = %v, want %v (Script=%q, Lang=%q)",
					got, tt.wantRTL, tt.font.Script, tt.font.Lang)
			}
		})
	}
}

func TestISO639CodesCount(t *testing.T) {
	// Verify we have 183 ISO-639 codes
	expectedCount := 183
	actualCount := len(ISO639Codes)

	if actualCount != expectedCount {
		t.Errorf("ISO639Codes length = %d, want %d", actualCount, expectedCount)
	}

	// Verify all codes are 2 characters
	for _, code := range ISO639Codes {
		if len(code) != 2 {
			t.Errorf("ISO639Codes contains invalid code %q (length %d, want 2)", code, len(code))
		}
	}

	// Verify all codes are lowercase
	for _, code := range ISO639Codes {
		if code != strings.ToLower(code) {
			t.Errorf("ISO639Codes contains non-lowercase code %q", code)
		}
	}

	// Verify no duplicates
	seen := make(map[string]bool)
	for _, code := range ISO639Codes {
		if seen[code] {
			t.Errorf("ISO639Codes contains duplicate code %q", code)
		}
		seen[code] = true
	}
}
