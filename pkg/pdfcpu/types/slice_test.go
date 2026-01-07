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

func TestMemberOf(t *testing.T) {
	tests := []struct {
		name string
		s    string
		list []string
		want bool
	}{
		{"found first", "a", []string{"a", "b", "c"}, true},
		{"found last", "c", []string{"a", "b", "c"}, true},
		{"found middle", "b", []string{"a", "b", "c"}, true},
		{"not found", "d", []string{"a", "b", "c"}, false},
		{"empty list", "a", []string{}, false},
		{"nil list", "a", nil, false},
		{"empty string found", "", []string{"", "a"}, true},
		{"empty string not found", "", []string{"a", "b"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MemberOf(tt.s, tt.list)
			if got != tt.want {
				t.Errorf("MemberOf(%q, %v) = %v, want %v", tt.s, tt.list, got, tt.want)
			}
		})
	}
}

func TestIntMemberOf(t *testing.T) {
	tests := []struct {
		name string
		i    int
		list []int
		want bool
	}{
		{"found first", 1, []int{1, 2, 3}, true},
		{"found last", 3, []int{1, 2, 3}, true},
		{"not found", 4, []int{1, 2, 3}, false},
		{"empty list", 1, []int{}, false},
		{"nil list", 1, nil, false},
		{"zero found", 0, []int{0, 1, 2}, true},
		{"negative found", -1, []int{-1, 0, 1}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IntMemberOf(tt.i, tt.list)
			if got != tt.want {
				t.Errorf("IntMemberOf(%d, %v) = %v, want %v", tt.i, tt.list, got, tt.want)
			}
		})
	}
}

func TestIndRefMemberOf(t *testing.T) {
	ref1 := *NewIndirectRef(1, 0)
	ref2 := *NewIndirectRef(2, 0)
	ref3 := *NewIndirectRef(3, 0)

	tests := []struct {
		name string
		ref  IndirectRef
		arr  Array
		want bool
	}{
		{"found", ref1, Array{ref1, ref2, ref3}, true},
		{"not found", *NewIndirectRef(4, 0), Array{ref1, ref2, ref3}, false},
		{"empty array", ref1, Array{}, false},
		{"nil array", ref1, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IndRefMemberOf(tt.ref, tt.arr)
			if got != tt.want {
				t.Errorf("IndRefMemberOf(%v, %v) = %v, want %v", tt.ref, tt.arr, got, tt.want)
			}
		})
	}
}

func TestEqualSlices(t *testing.T) {
	tests := []struct {
		name string
		a    []string
		b    []string
		want bool
	}{
		{"both empty", []string{}, []string{}, true},
		{"both nil", nil, nil, true},
		{"equal single", []string{"a"}, []string{"a"}, true},
		{"equal multiple", []string{"a", "b", "c"}, []string{"a", "b", "c"}, true},
		{"different length", []string{"a", "b"}, []string{"a"}, false},
		{"different content", []string{"a", "b"}, []string{"a", "c"}, false},
		{"different order", []string{"a", "b"}, []string{"b", "a"}, false},
		{"nil vs empty", nil, []string{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EqualSlices(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("EqualSlices(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
