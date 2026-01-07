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

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

func TestDestinationTypeStrings(t *testing.T) {
	tests := []struct {
		typ  DestinationType
		want string
	}{
		{DestXYZ, "XYZ"},
		{DestFit, "Fit"},
		{DestFitH, "FitH"},
		{DestFitV, "FitV"},
		{DestFitR, "FitR"},
		{DestFitB, "FitB"},
		{DestFitBH, "FitBH"},
		{DestFitBV, "FitBV"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := DestinationTypeStrings[tt.typ]
			if got != tt.want {
				t.Errorf("DestinationTypeStrings[%d] = %q, want %q", tt.typ, got, tt.want)
			}
		})
	}
}

func TestDestinationString(t *testing.T) {
	tests := []struct {
		name string
		dest Destination
		want string
	}{
		{"XYZ", Destination{Typ: DestXYZ}, "XYZ"},
		{"Fit", Destination{Typ: DestFit}, "Fit"},
		{"FitH", Destination{Typ: DestFitH}, "FitH"},
		{"FitV", Destination{Typ: DestFitV}, "FitV"},
		{"FitR", Destination{Typ: DestFitR}, "FitR"},
		{"FitB", Destination{Typ: DestFitB}, "FitB"},
		{"FitBH", Destination{Typ: DestFitBH}, "FitBH"},
		{"FitBV", Destination{Typ: DestFitBV}, "FitBV"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.dest.String()
			if got != tt.want {
				t.Errorf("Destination.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDestinationName(t *testing.T) {
	dest := Destination{Typ: DestXYZ}
	name := dest.Name()
	if name.Value() != "XYZ" {
		t.Errorf("Destination.Name() = %q, want %q", name.Value(), "XYZ")
	}
}

func TestDestinationArray(t *testing.T) {
	indRef := *types.NewIndirectRef(1, 0)

	t.Run("DestXYZ", func(t *testing.T) {
		dest := Destination{Typ: DestXYZ, Left: 10, Top: 20, Zoom: 1.5}
		arr := dest.Array(indRef)
		if len(arr) != 5 {
			t.Errorf("DestXYZ array len = %d, want 5", len(arr))
		}
	})

	t.Run("DestFit", func(t *testing.T) {
		dest := Destination{Typ: DestFit}
		arr := dest.Array(indRef)
		if len(arr) != 2 {
			t.Errorf("DestFit array len = %d, want 2", len(arr))
		}
	})

	t.Run("DestFitH", func(t *testing.T) {
		dest := Destination{Typ: DestFitH, Top: 100}
		arr := dest.Array(indRef)
		if len(arr) != 3 {
			t.Errorf("DestFitH array len = %d, want 3", len(arr))
		}
	})

	t.Run("DestFitV", func(t *testing.T) {
		dest := Destination{Typ: DestFitV, Left: 50}
		arr := dest.Array(indRef)
		if len(arr) != 3 {
			t.Errorf("DestFitV array len = %d, want 3", len(arr))
		}
	})

	t.Run("DestFitR", func(t *testing.T) {
		dest := Destination{Typ: DestFitR, Left: 10, Bottom: 20, Right: 100, Top: 200}
		arr := dest.Array(indRef)
		if len(arr) != 6 {
			t.Errorf("DestFitR array len = %d, want 6", len(arr))
		}
	})

	t.Run("DestFitB", func(t *testing.T) {
		dest := Destination{Typ: DestFitB}
		arr := dest.Array(indRef)
		if len(arr) != 2 {
			t.Errorf("DestFitB array len = %d, want 2", len(arr))
		}
	})

	t.Run("DestFitBH", func(t *testing.T) {
		dest := Destination{Typ: DestFitBH, Top: 100}
		arr := dest.Array(indRef)
		if len(arr) != 3 {
			t.Errorf("DestFitBH array len = %d, want 3", len(arr))
		}
	})

	t.Run("DestFitBV", func(t *testing.T) {
		dest := Destination{Typ: DestFitBV, Left: 50}
		arr := dest.Array(indRef)
		if len(arr) != 3 {
			t.Errorf("DestFitBV array len = %d, want 3", len(arr))
		}
	})
}

func TestDestinationTypeConstants(t *testing.T) {
	// Verify constants have expected order
	if DestXYZ != 0 {
		t.Errorf("DestXYZ = %d, want 0", DestXYZ)
	}
	if DestFit != 1 {
		t.Errorf("DestFit = %d, want 1", DestFit)
	}
	if DestFitBV != 7 {
		t.Errorf("DestFitBV = %d, want 7", DestFitBV)
	}
}
