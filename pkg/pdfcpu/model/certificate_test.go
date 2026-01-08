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
	"crypto/x509/pkix"
	"strings"
	"testing"
)

func TestIsPEM(t *testing.T) {
	tests := []struct {
		fname string
		want  bool
	}{
		{"cert.pem", true},
		{"cert.PEM", true},
		{"cert.Pem", true},
		{"cert.p7c", false},
		{"cert", false},
		{"", false},
		{"file.pem.txt", false},
		{"path/to/cert.pem", true},
	}

	for _, tt := range tests {
		t.Run(tt.fname, func(t *testing.T) {
			if got := IsPEM(tt.fname); got != tt.want {
				t.Errorf("IsPEM(%q) = %v, want %v", tt.fname, got, tt.want)
			}
		})
	}
}

func TestIsP7C(t *testing.T) {
	tests := []struct {
		fname string
		want  bool
	}{
		{"cert.p7c", true},
		{"cert.P7C", true},
		{"cert.P7c", true},
		{"cert.pem", false},
		{"cert", false},
		{"", false},
		{"file.p7c.txt", false},
		{"path/to/cert.p7c", true},
	}

	for _, tt := range tests {
		t.Run(tt.fname, func(t *testing.T) {
			if got := IsP7C(tt.fname); got != tt.want {
				t.Errorf("IsP7C(%q) = %v, want %v", tt.fname, got, tt.want)
			}
		})
	}
}

func TestStrSliceString(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{"empty", []string{}, ""},
		{"single", []string{"one"}, "one"},
		{"multiple", []string{"one", "two", "three"}, "one,two,three"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := strSliceString(tt.input); got != tt.want {
				t.Errorf("strSliceString(%v) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNameString(t *testing.T) {
	subj := pkix.Name{
		Organization:       []string{"Test Org"},
		OrganizationalUnit: []string{"Test Unit"},
		CommonName:         "Test Name",
		Country:            []string{"US"},
		Province:           []string{"California"},
		Locality:           []string{"San Francisco"},
		StreetAddress:      []string{"123 Test St"},
		PostalCode:         []string{"94102"},
	}

	s := nameString(subj)

	if !strings.Contains(s, "Test Org") {
		t.Error("nameString should contain Organization")
	}
	if !strings.Contains(s, "Test Unit") {
		t.Error("nameString should contain OrganizationalUnit")
	}
	if !strings.Contains(s, "Test Name") {
		t.Error("nameString should contain CommonName")
	}
	if !strings.Contains(s, "US") {
		t.Error("nameString should contain Country")
	}
	if !strings.Contains(s, "California") {
		t.Error("nameString should contain Province")
	}
	if !strings.Contains(s, "San Francisco") {
		t.Error("nameString should contain Locality")
	}
	if !strings.Contains(s, "123 Test St") {
		t.Error("nameString should contain StreetAddress")
	}
	if !strings.Contains(s, "94102") {
		t.Error("nameString should contain PostalCode")
	}
}

func TestNameStringMinimal(t *testing.T) {
	subj := pkix.Name{
		Organization: []string{"Minimal Org"},
	}

	s := nameString(subj)
	if !strings.Contains(s, "Minimal Org") {
		t.Error("nameString should contain Organization")
	}
	// Should not contain optional fields
	if strings.Contains(s, "name") {
		t.Error("nameString should not contain name when CommonName is empty")
	}
}
