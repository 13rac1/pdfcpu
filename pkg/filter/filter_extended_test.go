/*
Copyright 2018 The pdfcpu Authors.

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

package filter_test

import (
	"io"
	"strings"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/filter"
)

func TestSupportsDecodeParms(t *testing.T) {
	tests := []struct {
		filterName string
		want       bool
	}{
		{filter.CCITTFax, true},
		{filter.LZW, true},
		{filter.Flate, true},
		{filter.ASCII85, false},
		{filter.ASCIIHex, false},
		{filter.RunLength, false},
		{filter.DCT, false},
		{filter.JBIG2, false},
		{filter.JPX, false},
		{"InvalidFilter", false},
	}

	for _, tt := range tests {
		t.Run(tt.filterName, func(t *testing.T) {
			got := filter.SupportsDecodeParms(tt.filterName)
			if got != tt.want {
				t.Errorf("SupportsDecodeParms(%q) = %v, want %v", tt.filterName, got, tt.want)
			}
		})
	}
}

func TestList(t *testing.T) {
	list := filter.List()

	// Verify we have the expected filters
	expectedFilters := map[string]bool{
		filter.ASCII85:   true,
		filter.ASCIIHex:  true,
		filter.RunLength: true,
		filter.LZW:       true,
		filter.Flate:     true,
	}

	if len(list) != len(expectedFilters) {
		t.Errorf("List() returned %d filters, want %d", len(list), len(expectedFilters))
	}

	for _, f := range list {
		if !expectedFilters[f] {
			t.Errorf("List() contains unexpected filter %q", f)
		}
	}
}

func TestDecodeLength(t *testing.T) {
	// Test DecodeLength with various filters
	testFilters := []string{filter.ASCII85, filter.ASCIIHex, filter.RunLength, filter.LZW, filter.Flate}

	original := "Hello, World! This is a test message for DecodeLength functionality."

	for _, filterName := range testFilters {
		t.Run(filterName, func(t *testing.T) {
			f, err := filter.NewFilter(filterName, nil)
			if err != nil {
				t.Fatalf("NewFilter(%q) error = %v", filterName, err)
			}

			// Encode the data
			encoded, err := f.Encode(strings.NewReader(original))
			if err != nil {
				t.Fatalf("Encode() error = %v", err)
			}

			// Read encoded bytes
			encodedBytes, err := io.ReadAll(encoded)
			if err != nil {
				t.Fatalf("ReadAll(encoded) error = %v", err)
			}

			// Test DecodeLength with -1 (decode all)
			decoded, err := f.DecodeLength(strings.NewReader(string(encodedBytes)), -1)
			if err != nil {
				t.Fatalf("DecodeLength(-1) error = %v", err)
			}

			decodedBytes, err := io.ReadAll(decoded)
			if err != nil {
				t.Fatalf("ReadAll(decoded) error = %v", err)
			}

			if string(decodedBytes) != original {
				t.Errorf("DecodeLength(-1) = %q, want %q", string(decodedBytes), original)
			}

			// Test DecodeLength with a positive limit
			decoded2, err := f.DecodeLength(strings.NewReader(string(encodedBytes)), 10)
			if err != nil {
				t.Fatalf("DecodeLength(10) error = %v", err)
			}

			decodedBytes2, err := io.ReadAll(decoded2)
			if err != nil {
				t.Fatalf("ReadAll(decoded2) error = %v", err)
			}

			// Should have at least 10 bytes decoded
			if len(decodedBytes2) < 10 {
				t.Errorf("DecodeLength(10) returned %d bytes, want at least 10", len(decodedBytes2))
			}
		})
	}
}

func TestNewFilterWithParms(t *testing.T) {
	// Test filters that accept parameters
	parms := map[string]int{
		"Predictor": 12,
		"Columns":   100,
	}

	// Flate with parms should work
	f, err := filter.NewFilter(filter.Flate, parms)
	if err != nil {
		t.Errorf("NewFilter(Flate, parms) error = %v", err)
	}
	if f == nil {
		t.Error("NewFilter(Flate, parms) returned nil filter")
	}

	// LZW with parms should work
	f, err = filter.NewFilter(filter.LZW, parms)
	if err != nil {
		t.Errorf("NewFilter(LZW, parms) error = %v", err)
	}
	if f == nil {
		t.Error("NewFilter(LZW, parms) returned nil filter")
	}
}

func TestFilterConstants(t *testing.T) {
	// Verify filter constants have expected values
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"ASCII85", filter.ASCII85, "ASCII85Decode"},
		{"ASCIIHex", filter.ASCIIHex, "ASCIIHexDecode"},
		{"RunLength", filter.RunLength, "RunLengthDecode"},
		{"LZW", filter.LZW, "LZWDecode"},
		{"Flate", filter.Flate, "FlateDecode"},
		{"CCITTFax", filter.CCITTFax, "CCITTFaxDecode"},
		{"JBIG2", filter.JBIG2, "JBIG2Decode"},
		{"DCT", filter.DCT, "DCTDecode"},
		{"JPX", filter.JPX, "JPXDecode"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestErrUnsupportedFilter(t *testing.T) {
	if filter.ErrUnsupportedFilter == nil {
		t.Error("ErrUnsupportedFilter should not be nil")
	}
	if filter.ErrUnsupportedFilter.Error() == "" {
		t.Error("ErrUnsupportedFilter.Error() should not be empty")
	}
}

func TestASCII85DecodeMissingEOD(t *testing.T) {
	f, err := filter.NewFilter(filter.ASCII85, nil)
	if err != nil {
		t.Fatalf("NewFilter(ASCII85) error = %v", err)
	}

	// Encode some data
	encoded, err := f.Encode(strings.NewReader("Hello"))
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}

	// Read encoded bytes
	encodedBytes, err := io.ReadAll(encoded)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}

	// Remove the EOD marker (~>)
	if len(encodedBytes) >= 2 {
		encodedBytes = encodedBytes[:len(encodedBytes)-2]
	}

	// Attempt to decode - should fail
	_, err = f.Decode(strings.NewReader(string(encodedBytes)))
	if err == nil {
		t.Error("Decode() should return error when EOD marker is missing")
	}
	if !strings.Contains(err.Error(), "missing eod marker") {
		t.Errorf("Decode() error = %q, want error containing 'missing eod marker'", err.Error())
	}
}

func TestASCIIHexDecodeInvalidHex(t *testing.T) {
	f, err := filter.NewFilter(filter.ASCIIHex, nil)
	if err != nil {
		t.Fatalf("NewFilter(ASCIIHex) error = %v", err)
	}

	tests := []struct {
		name  string
		input string
	}{
		{"invalid chars", "GHIJ>"},   // G, H, I, J are not valid hex
		{"mixed invalid", "48GG65>"}, // GG is invalid
		{"only invalid", "XYZ>"},     // X, Y, Z are not valid hex
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := f.Decode(strings.NewReader(tt.input))
			if err == nil {
				t.Errorf("Decode(%q) should return error for invalid hex", tt.input)
			}
		})
	}
}

func TestLZWDecodeUnsupportedPredictor(t *testing.T) {
	// LZW with predictor > 1 should fail
	parms := map[string]int{"Predictor": 12}
	f, err := filter.NewFilter(filter.LZW, parms)
	if err != nil {
		t.Fatalf("NewFilter(LZW, parms) error = %v", err)
	}

	// First encode some data without predictor
	fNoPredictor, _ := filter.NewFilter(filter.LZW, nil)
	encoded, err := fNoPredictor.Encode(strings.NewReader("Hello, World!"))
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}

	encodedBytes, err := io.ReadAll(encoded)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}

	// Try to decode with predictor - should fail
	_, err = f.Decode(strings.NewReader(string(encodedBytes)))
	if err == nil {
		t.Error("Decode() should return error for unsupported predictor")
	}
	if !strings.Contains(err.Error(), "unsupported predictor") {
		t.Errorf("Decode() error = %q, want error containing 'unsupported predictor'", err.Error())
	}
}

func TestFlateDecodeInvalidPredictor(t *testing.T) {
	// Flate with invalid predictor value should fail
	parms := map[string]int{"Predictor": 99} // Invalid predictor
	f, err := filter.NewFilter(filter.Flate, parms)
	if err != nil {
		t.Fatalf("NewFilter(Flate, parms) error = %v", err)
	}

	// First encode some data without predictor
	fNoPredictor, _ := filter.NewFilter(filter.Flate, nil)
	encoded, err := fNoPredictor.Encode(strings.NewReader("Hello, World!"))
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}

	encodedBytes, err := io.ReadAll(encoded)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}

	// Try to decode with invalid predictor - should fail
	_, err = f.Decode(strings.NewReader(string(encodedBytes)))
	if err == nil {
		t.Error("Decode() should return error for invalid predictor")
	}
	if !strings.Contains(err.Error(), "undefined") || !strings.Contains(err.Error(), "Predictor") {
		t.Errorf("Decode() error = %q, want error about undefined Predictor", err.Error())
	}
}

func TestFlateDecodeInvalidBPC(t *testing.T) {
	// Flate with invalid BitsPerComponent should fail
	parms := map[string]int{
		"Predictor":        12,
		"BitsPerComponent": 7, // Invalid - must be 1, 2, 4, 8, or 16
	}
	f, err := filter.NewFilter(filter.Flate, parms)
	if err != nil {
		t.Fatalf("NewFilter(Flate, parms) error = %v", err)
	}

	// First encode some data
	fNoPredictor, _ := filter.NewFilter(filter.Flate, nil)
	encoded, err := fNoPredictor.Encode(strings.NewReader("Hello, World!"))
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}

	encodedBytes, err := io.ReadAll(encoded)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}

	// Try to decode - should fail due to invalid BPC
	_, err = f.Decode(strings.NewReader(string(encodedBytes)))
	if err == nil {
		t.Error("Decode() should return error for invalid BitsPerComponent")
	}
	if !strings.Contains(err.Error(), "BitsPerComponent") {
		t.Errorf("Decode() error = %q, want error about BitsPerComponent", err.Error())
	}
}

func TestFlateDecodeZeroColors(t *testing.T) {
	// Flate with Colors=0 should fail
	parms := map[string]int{
		"Predictor": 12,
		"Colors":    0, // Invalid - must be > 0
	}
	f, err := filter.NewFilter(filter.Flate, parms)
	if err != nil {
		t.Fatalf("NewFilter(Flate, parms) error = %v", err)
	}

	// First encode some data
	fNoPredictor, _ := filter.NewFilter(filter.Flate, nil)
	encoded, err := fNoPredictor.Encode(strings.NewReader("Hello, World!"))
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}

	encodedBytes, err := io.ReadAll(encoded)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}

	// Try to decode - should fail due to Colors=0
	_, err = f.Decode(strings.NewReader(string(encodedBytes)))
	if err == nil {
		t.Error("Decode() should return error for Colors=0")
	}
	if !strings.Contains(err.Error(), "Colors") {
		t.Errorf("Decode() error = %q, want error about Colors", err.Error())
	}
}
