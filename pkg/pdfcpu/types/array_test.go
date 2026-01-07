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

package types

import (
	"strings"
	"testing"
)

func TestNewStringLiteralArray(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		arr := NewStringLiteralArray()
		if len(arr) != 0 {
			t.Errorf("len = %d, want 0", len(arr))
		}
	})

	t.Run("single", func(t *testing.T) {
		arr := NewStringLiteralArray("hello")
		if len(arr) != 1 {
			t.Errorf("len = %d, want 1", len(arr))
		}
		if s, ok := arr[0].(StringLiteral); !ok || s.Value() != "hello" {
			t.Errorf("arr[0] = %v, want StringLiteral(hello)", arr[0])
		}
	})

	t.Run("multiple", func(t *testing.T) {
		arr := NewStringLiteralArray("a", "b", "c")
		if len(arr) != 3 {
			t.Errorf("len = %d, want 3", len(arr))
		}
	})
}

func TestNewHexLiteralArray(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		arr := NewHexLiteralArray()
		if len(arr) != 0 {
			t.Errorf("len = %d, want 0", len(arr))
		}
	})

	t.Run("single", func(t *testing.T) {
		arr := NewHexLiteralArray("Hi")
		if len(arr) != 1 {
			t.Errorf("len = %d, want 1", len(arr))
		}
		if h, ok := arr[0].(HexLiteral); !ok {
			t.Errorf("arr[0] = %T, want HexLiteral", arr[0])
		} else {
			b, _ := h.Bytes()
			if string(b) != "Hi" {
				t.Errorf("arr[0].Bytes() = %q, want %q", string(b), "Hi")
			}
		}
	})
}

func TestNewNameArray(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		arr := NewNameArray()
		if len(arr) != 0 {
			t.Errorf("len = %d, want 0", len(arr))
		}
	})

	t.Run("multiple", func(t *testing.T) {
		arr := NewNameArray("Type", "Page", "Font")
		if len(arr) != 3 {
			t.Errorf("len = %d, want 3", len(arr))
		}
		if n, ok := arr[0].(Name); !ok || n.Value() != "Type" {
			t.Errorf("arr[0] = %v, want Name(Type)", arr[0])
		}
	})
}

func TestNewNumberArray(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		arr := NewNumberArray()
		if len(arr) != 0 {
			t.Errorf("len = %d, want 0", len(arr))
		}
	})

	t.Run("multiple", func(t *testing.T) {
		arr := NewNumberArray(1.5, 2.5, 3.5)
		if len(arr) != 3 {
			t.Errorf("len = %d, want 3", len(arr))
		}
		if f, ok := arr[0].(Float); !ok || f.Value() != 1.5 {
			t.Errorf("arr[0] = %v, want Float(1.5)", arr[0])
		}
	})
}

func TestNewIntegerArray(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		arr := NewIntegerArray()
		if len(arr) != 0 {
			t.Errorf("len = %d, want 0", len(arr))
		}
	})

	t.Run("multiple", func(t *testing.T) {
		arr := NewIntegerArray(1, 2, 3)
		if len(arr) != 3 {
			t.Errorf("len = %d, want 3", len(arr))
		}
		if i, ok := arr[0].(Integer); !ok || i.Value() != 1 {
			t.Errorf("arr[0] = %v, want Integer(1)", arr[0])
		}
	})
}

func TestArrayClone(t *testing.T) {
	t.Run("empty array", func(t *testing.T) {
		arr := Array{}
		clone := arr.Clone().(Array)
		if len(clone) != 0 {
			t.Errorf("clone len = %d, want 0", len(clone))
		}
	})

	t.Run("with nil element", func(t *testing.T) {
		arr := Array{Integer(1), nil, Integer(3)}
		clone := arr.Clone().(Array)
		if len(clone) != 3 {
			t.Errorf("clone len = %d, want 3", len(clone))
		}
		if clone[1] != nil {
			t.Errorf("clone[1] = %v, want nil", clone[1])
		}
	})

	t.Run("mixed types", func(t *testing.T) {
		arr := Array{Integer(1), Float(2.5), StringLiteral("hello")}
		clone := arr.Clone().(Array)
		if len(clone) != 3 {
			t.Errorf("clone len = %d, want 3", len(clone))
		}
		// Verify independence - modifying original doesn't affect clone
		if clone[0].(Integer).Value() != 1 {
			t.Errorf("clone[0] = %v, want 1", clone[0])
		}
	})
}

func TestArrayString(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		arr := Array{}
		s := arr.String()
		if !strings.HasPrefix(s, "[") || !strings.HasSuffix(s, "]") {
			t.Errorf("String() = %q, should be bracketed", s)
		}
	})

	t.Run("with elements", func(t *testing.T) {
		arr := Array{Integer(1), Integer(2)}
		s := arr.String()
		if !strings.Contains(s, "1") || !strings.Contains(s, "2") {
			t.Errorf("String() = %q, should contain elements", s)
		}
	})

	t.Run("with null", func(t *testing.T) {
		arr := Array{nil, Integer(1)}
		s := arr.String()
		if !strings.Contains(s, "null") {
			t.Errorf("String() = %q, should contain null", s)
		}
	})

	t.Run("with nested array", func(t *testing.T) {
		inner := Array{Integer(1)}
		arr := Array{inner}
		s := arr.String()
		if s == "" {
			t.Error("String() should not be empty")
		}
	})
}

func TestArrayPDFString(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		arr := Array{}
		s := arr.PDFString()
		if s != "[]" {
			t.Errorf("PDFString() = %q, want []", s)
		}
	})

	t.Run("with integers", func(t *testing.T) {
		arr := Array{Integer(1), Integer(2)}
		s := arr.PDFString()
		if !strings.HasPrefix(s, "[") || !strings.HasSuffix(s, "]") {
			t.Errorf("PDFString() = %q, should be bracketed", s)
		}
	})

	t.Run("with null", func(t *testing.T) {
		arr := Array{nil}
		s := arr.PDFString()
		if !strings.Contains(s, "null") {
			t.Errorf("PDFString() = %q, should contain null", s)
		}
	})

	t.Run("with float", func(t *testing.T) {
		arr := Array{Float(1.5)}
		s := arr.PDFString()
		if !strings.Contains(s, "1.5") {
			t.Errorf("PDFString() = %q, should contain 1.5", s)
		}
	})

	t.Run("with boolean", func(t *testing.T) {
		arr := Array{Boolean(true)}
		s := arr.PDFString()
		if !strings.Contains(s, "true") {
			t.Errorf("PDFString() = %q, should contain true", s)
		}
	})

	t.Run("with string literal", func(t *testing.T) {
		arr := Array{StringLiteral("test")}
		s := arr.PDFString()
		if !strings.Contains(s, "(test)") {
			t.Errorf("PDFString() = %q, should contain (test)", s)
		}
	})

	t.Run("with hex literal", func(t *testing.T) {
		arr := Array{HexLiteral("48656c6c6f")}
		s := arr.PDFString()
		if !strings.Contains(s, "<48656c6c6f>") {
			t.Errorf("PDFString() = %q, should contain hex", s)
		}
	})

	t.Run("with name", func(t *testing.T) {
		arr := Array{Name("Test")}
		s := arr.PDFString()
		if !strings.Contains(s, "/") {
			t.Errorf("PDFString() = %q, should contain /", s)
		}
	})

	t.Run("with indirect ref", func(t *testing.T) {
		arr := Array{*NewIndirectRef(10, 0)}
		s := arr.PDFString()
		if !strings.Contains(s, "10 0 R") {
			t.Errorf("PDFString() = %q, should contain indirect ref", s)
		}
	})

	t.Run("with nested array", func(t *testing.T) {
		inner := Array{Integer(1)}
		arr := Array{inner}
		s := arr.PDFString()
		if !strings.Contains(s, "[[") {
			t.Errorf("PDFString() = %q, should contain nested brackets", s)
		}
	})
}

func TestArrayRemoveNulls(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		arr := Array{}
		result := arr.RemoveNulls()
		if len(result) != 0 {
			t.Errorf("RemoveNulls() len = %d, want 0", len(result))
		}
	})

	t.Run("no nulls", func(t *testing.T) {
		arr := Array{Integer(1), Integer(2)}
		result := arr.RemoveNulls()
		if len(result) != 2 {
			t.Errorf("RemoveNulls() len = %d, want 2", len(result))
		}
	})

	t.Run("with nulls", func(t *testing.T) {
		arr := Array{nil, Integer(1), nil, Integer(2), nil}
		result := arr.RemoveNulls()
		if len(result) != 2 {
			t.Errorf("RemoveNulls() len = %d, want 2", len(result))
		}
		if result[0].(Integer).Value() != 1 {
			t.Errorf("result[0] = %v, want 1", result[0])
		}
	})

	t.Run("all nulls", func(t *testing.T) {
		arr := Array{nil, nil, nil}
		result := arr.RemoveNulls()
		if len(result) != 0 {
			t.Errorf("RemoveNulls() len = %d, want 0", len(result))
		}
	})
}
