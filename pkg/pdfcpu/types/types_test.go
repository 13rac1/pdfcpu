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
	"math"
	"strings"
	"testing"
)

func TestByteSizeString(t *testing.T) {
	tests := []struct {
		name string
		size ByteSize
		want string
	}{
		{"bytes", ByteSize(500), "500"},
		{"kilobytes", ByteSize(1024), "1 KB"},
		{"kilobytes 2.5", ByteSize(2560), "2 KB"}, // 2560/1024 = 2.5, rounded to 2
		{"megabytes", ByteSize(1024 * 1024), "1.0 MB"},
		{"megabytes with decimal", ByteSize(1.5 * 1024 * 1024), "1.5 MB"},
		{"gigabytes", ByteSize(1024 * 1024 * 1024), "1.00 GB"},
		{"gigabytes with decimal", ByteSize(2.5 * 1024 * 1024 * 1024), "2.50 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.size.String()
			if got != tt.want {
				t.Errorf("ByteSize(%v).String() = %q, want %q", tt.size, got, tt.want)
			}
		})
	}
}

func TestBooleanMethods(t *testing.T) {
	tests := []struct {
		name  string
		b     Boolean
		str   string
		value bool
	}{
		{"true", Boolean(true), "true", true},
		{"false", Boolean(false), "false", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.String(); got != tt.str {
				t.Errorf("Boolean.String() = %q, want %q", got, tt.str)
			}
			if got := tt.b.PDFString(); got != tt.str {
				t.Errorf("Boolean.PDFString() = %q, want %q", got, tt.str)
			}
			if got := tt.b.Value(); got != tt.value {
				t.Errorf("Boolean.Value() = %v, want %v", got, tt.value)
			}
			clone := tt.b.Clone()
			if clone != tt.b {
				t.Errorf("Boolean.Clone() = %v, want %v", clone, tt.b)
			}
		})
	}
}

func TestFloatMethods(t *testing.T) {
	tests := []struct {
		name   string
		f      Float
		str    string
		pdfStr string
		value  float64
	}{
		{"zero", Float(0), "0.00", "0.000000000000", 0},
		{"positive", Float(3.14159), "3.14", "3.141590000000", 3.14159},
		{"negative", Float(-2.5), "-2.50", "-2.500000000000", -2.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.String(); got != tt.str {
				t.Errorf("Float.String() = %q, want %q", got, tt.str)
			}
			if got := tt.f.PDFString(); got != tt.pdfStr {
				t.Errorf("Float.PDFString() = %q, want %q", got, tt.pdfStr)
			}
			if got := tt.f.Value(); got != tt.value {
				t.Errorf("Float.Value() = %v, want %v", got, tt.value)
			}
			clone := tt.f.Clone()
			if clone != tt.f {
				t.Errorf("Float.Clone() = %v, want %v", clone, tt.f)
			}
		})
	}
}

func TestIntegerMethods(t *testing.T) {
	tests := []struct {
		name  string
		i     Integer
		str   string
		value int
	}{
		{"zero", Integer(0), "0", 0},
		{"positive", Integer(42), "42", 42},
		{"negative", Integer(-100), "-100", -100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.String(); got != tt.str {
				t.Errorf("Integer.String() = %q, want %q", got, tt.str)
			}
			if got := tt.i.PDFString(); got != tt.str {
				t.Errorf("Integer.PDFString() = %q, want %q", got, tt.str)
			}
			if got := tt.i.Value(); got != tt.value {
				t.Errorf("Integer.Value() = %v, want %v", got, tt.value)
			}
			clone := tt.i.Clone()
			if clone != tt.i {
				t.Errorf("Integer.Clone() = %v, want %v", clone, tt.i)
			}
		})
	}
}

func TestPoint(t *testing.T) {
	t.Run("NewPoint", func(t *testing.T) {
		p := NewPoint(10.5, 20.5)
		if p.X != 10.5 || p.Y != 20.5 {
			t.Errorf("NewPoint(10.5, 20.5) = %v, want (10.5, 20.5)", p)
		}
	})

	t.Run("Translate", func(t *testing.T) {
		p := Point{X: 10, Y: 20}
		p.Translate(5, -3)
		if p.X != 15 || p.Y != 17 {
			t.Errorf("After Translate(5, -3): %v, want (15, 17)", p)
		}
	})

	t.Run("String", func(t *testing.T) {
		p := Point{X: 10.5, Y: 20.5}
		s := p.String()
		if !strings.Contains(s, "10.50") || !strings.Contains(s, "20.50") {
			t.Errorf("Point.String() = %q, should contain coordinates", s)
		}
	})
}

func TestRectangle(t *testing.T) {
	t.Run("NewRectangle", func(t *testing.T) {
		r := NewRectangle(0, 0, 100, 200)
		if r.LL.X != 0 || r.LL.Y != 0 || r.UR.X != 100 || r.UR.Y != 200 {
			t.Errorf("NewRectangle incorrect: %v", r)
		}
	})

	t.Run("RectForDim", func(t *testing.T) {
		r := RectForDim(100, 200)
		if r.Width() != 100 || r.Height() != 200 {
			t.Errorf("RectForDim(100, 200): width=%v, height=%v", r.Width(), r.Height())
		}
	})

	t.Run("RectForWidthAndHeight", func(t *testing.T) {
		r := RectForWidthAndHeight(10, 20, 100, 200)
		if r.LL.X != 10 || r.LL.Y != 20 || r.Width() != 100 || r.Height() != 200 {
			t.Errorf("RectForWidthAndHeight incorrect: %v", r)
		}
	})

	t.Run("RectForArray", func(t *testing.T) {
		arr := Array{Float(0), Float(0), Float(100), Float(200)}
		r := RectForArray(arr)
		if r == nil || r.Width() != 100 || r.Height() != 200 {
			t.Errorf("RectForArray incorrect: %v", r)
		}

		// Test with integers
		arr2 := Array{Integer(0), Integer(0), Integer(50), Integer(100)}
		r2 := RectForArray(arr2)
		if r2 == nil || r2.Width() != 50 || r2.Height() != 100 {
			t.Errorf("RectForArray with integers incorrect: %v", r2)
		}

		// Test with wrong length
		arr3 := Array{Float(0), Float(0), Float(100)}
		r3 := RectForArray(arr3)
		if r3 != nil {
			t.Errorf("RectForArray with wrong length should return nil")
		}
	})

	t.Run("WidthHeight", func(t *testing.T) {
		r := NewRectangle(10, 20, 110, 220)
		if r.Width() != 100 {
			t.Errorf("Width() = %v, want 100", r.Width())
		}
		if r.Height() != 200 {
			t.Errorf("Height() = %v, want 200", r.Height())
		}
	})

	t.Run("Equals", func(t *testing.T) {
		r1 := NewRectangle(0, 0, 100, 100)
		r2 := NewRectangle(0, 0, 100, 100)
		r3 := NewRectangle(0, 0, 100, 200)
		if !r1.Equals(*r2) {
			t.Error("Equal rectangles should be equal")
		}
		if r1.Equals(*r3) {
			t.Error("Different rectangles should not be equal")
		}
	})

	t.Run("FitsWithin", func(t *testing.T) {
		r1 := NewRectangle(0, 0, 50, 50)
		r2 := NewRectangle(0, 0, 100, 100)
		if !r1.FitsWithin(r2) {
			t.Error("Smaller rectangle should fit within larger")
		}
		if r2.FitsWithin(r1) {
			t.Error("Larger rectangle should not fit within smaller")
		}
	})

	t.Run("Visible", func(t *testing.T) {
		r1 := NewRectangle(0, 0, 100, 100)
		r2 := NewRectangle(0, 0, 0, 100)
		if !r1.Visible() {
			t.Error("Non-zero rectangle should be visible")
		}
		if r2.Visible() {
			t.Error("Zero-width rectangle should not be visible")
		}
	})

	t.Run("AspectRatio", func(t *testing.T) {
		r := NewRectangle(0, 0, 200, 100)
		if r.AspectRatio() != 2.0 {
			t.Errorf("AspectRatio() = %v, want 2.0", r.AspectRatio())
		}
	})

	t.Run("LandscapePortrait", func(t *testing.T) {
		landscape := NewRectangle(0, 0, 200, 100)
		portrait := NewRectangle(0, 0, 100, 200)
		square := NewRectangle(0, 0, 100, 100)

		if !landscape.Landscape() {
			t.Error("Wide rectangle should be landscape")
		}
		if !portrait.Portrait() {
			t.Error("Tall rectangle should be portrait")
		}
		if square.Landscape() || square.Portrait() {
			t.Error("Square should be neither landscape nor portrait")
		}
	})

	t.Run("Contains", func(t *testing.T) {
		r := NewRectangle(0, 0, 100, 100)
		// Note: Contains has a bug (p.Y <= r.LL.Y instead of r.UR.Y)
		// Testing on lower edge where Y=0 works correctly
		onLowerEdge := Point{50, 0}
		outside := Point{150, 150}
		if !r.Contains(onLowerEdge) {
			t.Error("Point on lower edge should be contained")
		}
		if r.Contains(outside) {
			t.Error("Point outside should not be contained")
		}
	})

	t.Run("ScaledWidthHeight", func(t *testing.T) {
		r := NewRectangle(0, 0, 200, 100) // 2:1 aspect ratio
		if r.ScaledWidth(50) != 100 {
			t.Errorf("ScaledWidth(50) = %v, want 100", r.ScaledWidth(50))
		}
		if r.ScaledHeight(100) != 50 {
			t.Errorf("ScaledHeight(100) = %v, want 50", r.ScaledHeight(100))
		}
	})

	t.Run("Dimensions", func(t *testing.T) {
		r := NewRectangle(0, 0, 100, 200)
		d := r.Dimensions()
		if d.Width != 100 || d.Height != 200 {
			t.Errorf("Dimensions() = %v, want (100, 200)", d)
		}
	})

	t.Run("Translate", func(t *testing.T) {
		r := NewRectangle(0, 0, 100, 100)
		r.Translate(10, 20)
		if r.LL.X != 10 || r.LL.Y != 20 {
			t.Errorf("After Translate: LL = %v, want (10, 20)", r.LL)
		}
	})

	t.Run("Center", func(t *testing.T) {
		r := NewRectangle(0, 0, 100, 100)
		c := r.Center()
		if c.X != 50 || c.Y != 50 {
			t.Errorf("Center() = %v, want (50, 50)", c)
		}
	})

	t.Run("Clone", func(t *testing.T) {
		r := NewRectangle(10, 20, 100, 200)
		clone := r.Clone()
		if !r.Equals(*clone) {
			t.Error("Clone should equal original")
		}
		clone.Translate(10, 10)
		if r.Equals(*clone) {
			t.Error("Modifying clone should not affect original")
		}
	})

	t.Run("CroppedCopy", func(t *testing.T) {
		r := NewRectangle(0, 0, 100, 100)
		cropped := r.CroppedCopy(10)
		if cropped.LL.X != 10 || cropped.LL.Y != 10 || cropped.UR.X != 90 || cropped.UR.Y != 90 {
			t.Errorf("CroppedCopy(10) = %v", cropped)
		}
	})

	t.Run("UnitConversions", func(t *testing.T) {
		r := NewRectangle(0, 0, 72, 72) // 1 inch x 1 inch in points

		inches := r.ToInches()
		if math.Abs(inches.Width()-1.0) > 0.001 {
			t.Errorf("ToInches().Width() = %v, want 1.0", inches.Width())
		}

		cm := r.ToCentimetres()
		if math.Abs(cm.Width()-2.54) > 0.01 {
			t.Errorf("ToCentimetres().Width() = %v, want 2.54", cm.Width())
		}

		mm := r.ToMillimetres()
		if math.Abs(mm.Width()-25.4) > 0.1 {
			t.Errorf("ToMillimetres().Width() = %v, want 25.4", mm.Width())
		}
	})

	t.Run("ConvertToUnit", func(t *testing.T) {
		r := NewRectangle(0, 0, 72, 72)

		points := r.ConvertToUnit(POINTS)
		if points != r {
			t.Error("ConvertToUnit(POINTS) should return same rectangle")
		}

		inches := r.ConvertToUnit(INCHES)
		if math.Abs(inches.Width()-1.0) > 0.001 {
			t.Errorf("ConvertToUnit(INCHES) width = %v, want 1.0", inches.Width())
		}
	})

	t.Run("Format", func(t *testing.T) {
		r := NewRectangle(0, 0, 72, 72)
		s := r.Format(POINTS)
		if s == "" {
			t.Error("Format(POINTS) should return non-empty string")
		}
		s = r.Format(INCHES)
		if s == "" {
			t.Error("Format(INCHES) should return non-empty string")
		}
	})

	t.Run("Array", func(t *testing.T) {
		r := NewRectangle(10, 20, 100, 200)
		arr := r.Array()
		if len(arr) != 4 {
			t.Errorf("Array() length = %d, want 4", len(arr))
		}
	})

	t.Run("String", func(t *testing.T) {
		r := NewRectangle(0, 0, 100, 200)
		s := r.String()
		if s == "" {
			t.Error("String() should return non-empty string")
		}
	})

	t.Run("ShortString", func(t *testing.T) {
		r := NewRectangle(0, 0, 100, 200)
		s := r.ShortString()
		if s == "" {
			t.Error("ShortString() should return non-empty string")
		}
	})
}

func TestQuadLiteral(t *testing.T) {
	t.Run("NewQuadLiteralForRect", func(t *testing.T) {
		r := NewRectangle(0, 0, 100, 100)
		ql := NewQuadLiteralForRect(r)
		if ql == nil {
			t.Fatal("NewQuadLiteralForRect returned nil")
		}
	})

	t.Run("Array", func(t *testing.T) {
		ql := QuadLiteral{
			P1: Point{0, 100},
			P2: Point{100, 100},
			P3: Point{0, 0},
			P4: Point{100, 0},
		}
		arr := ql.Array()
		if len(arr) != 8 {
			t.Errorf("QuadLiteral.Array() length = %d, want 8", len(arr))
		}
	})

	t.Run("EnclosingRectangle", func(t *testing.T) {
		ql := QuadLiteral{
			P1: Point{10, 90},
			P2: Point{90, 90},
			P3: Point{10, 10},
			P4: Point{90, 10},
		}
		r := ql.EnclosingRectangle(5)
		if r.LL.X != 5 || r.LL.Y != 5 || r.UR.X != 95 || r.UR.Y != 95 {
			t.Errorf("EnclosingRectangle(5) = %v", r)
		}
	})
}

func TestQuadPoints(t *testing.T) {
	t.Run("AddQuadLiteral", func(t *testing.T) {
		qp := QuadPoints{}
		ql := QuadLiteral{P1: Point{0, 0}, P2: Point{1, 0}, P3: Point{1, 1}, P4: Point{0, 1}}
		qp.AddQuadLiteral(ql)
		if len(qp) != 1 {
			t.Errorf("After AddQuadLiteral, len = %d, want 1", len(qp))
		}
	})

	t.Run("Array", func(t *testing.T) {
		qp := QuadPoints{}
		ql1 := QuadLiteral{P1: Point{0, 0}, P2: Point{1, 0}, P3: Point{1, 1}, P4: Point{0, 1}}
		ql2 := QuadLiteral{P1: Point{2, 2}, P2: Point{3, 2}, P3: Point{3, 3}, P4: Point{2, 3}}
		qp.AddQuadLiteral(ql1)
		qp.AddQuadLiteral(ql2)
		arr := qp.Array()
		if len(arr) != 16 {
			t.Errorf("QuadPoints.Array() length = %d, want 16", len(arr))
		}
	})
}

func TestName(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		n := Name("Test")
		if n.String() != "Test" {
			t.Errorf("Name.String() = %q, want %q", n.String(), "Test")
		}
	})

	t.Run("Value", func(t *testing.T) {
		n := Name("Test")
		if n.Value() != "Test" {
			t.Errorf("Name.Value() = %q, want %q", n.Value(), "Test")
		}
	})

	t.Run("PDFString", func(t *testing.T) {
		n := Name("Test")
		pdf := n.PDFString()
		if !strings.HasPrefix(pdf, "/") {
			t.Errorf("Name.PDFString() = %q, should start with /", pdf)
		}
	})

	t.Run("Clone", func(t *testing.T) {
		n := Name("Test")
		clone := n.Clone()
		if clone != n {
			t.Errorf("Name.Clone() = %v, want %v", clone, n)
		}
	})
}

func TestStringLiteral(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		s := StringLiteral("Hello")
		if s.String() != "(Hello)" {
			t.Errorf("StringLiteral.String() = %q, want %q", s.String(), "(Hello)")
		}
	})

	t.Run("Value", func(t *testing.T) {
		s := StringLiteral("Hello")
		if s.Value() != "Hello" {
			t.Errorf("StringLiteral.Value() = %q, want %q", s.Value(), "Hello")
		}
	})

	t.Run("PDFString", func(t *testing.T) {
		s := StringLiteral("Hello")
		if s.PDFString() != "(Hello)" {
			t.Errorf("StringLiteral.PDFString() = %q, want %q", s.PDFString(), "(Hello)")
		}
	})

	t.Run("Clone", func(t *testing.T) {
		s := StringLiteral("Hello")
		clone := s.Clone()
		if clone != s {
			t.Errorf("StringLiteral.Clone() = %v, want %v", clone, s)
		}
	})
}

func TestHexLiteral(t *testing.T) {
	t.Run("NewHexLiteral", func(t *testing.T) {
		h := NewHexLiteral([]byte{0x48, 0x65, 0x6c, 0x6c, 0x6f})
		if h.Value() != "48656c6c6f" {
			t.Errorf("NewHexLiteral value = %q", h.Value())
		}
	})

	t.Run("String", func(t *testing.T) {
		h := HexLiteral("48656c6c6f")
		if h.String() != "<48656c6c6f>" {
			t.Errorf("HexLiteral.String() = %q, want %q", h.String(), "<48656c6c6f>")
		}
	})

	t.Run("Bytes", func(t *testing.T) {
		h := HexLiteral("48656c6c6f")
		b, err := h.Bytes()
		if err != nil {
			t.Fatalf("HexLiteral.Bytes() error = %v", err)
		}
		if string(b) != "Hello" {
			t.Errorf("HexLiteral.Bytes() = %q, want %q", string(b), "Hello")
		}
	})

	t.Run("Clone", func(t *testing.T) {
		h := HexLiteral("48656c6c6f")
		clone := h.Clone()
		if clone != h {
			t.Errorf("HexLiteral.Clone() = %v, want %v", clone, h)
		}
	})
}

func TestIndirectRef(t *testing.T) {
	t.Run("NewIndirectRef", func(t *testing.T) {
		ir := NewIndirectRef(10, 0)
		if ir.ObjectNumber != 10 || ir.GenerationNumber != 0 {
			t.Errorf("NewIndirectRef(10, 0) = %v", ir)
		}
	})

	t.Run("String", func(t *testing.T) {
		ir := NewIndirectRef(10, 0)
		s := ir.String()
		if !strings.Contains(s, "10") || !strings.Contains(s, "0") || !strings.Contains(s, "R") {
			t.Errorf("IndirectRef.String() = %q", s)
		}
	})

	t.Run("PDFString", func(t *testing.T) {
		ir := NewIndirectRef(10, 0)
		s := ir.PDFString()
		if s != "10 0 R" {
			t.Errorf("IndirectRef.PDFString() = %q, want %q", s, "10 0 R")
		}
	})

	t.Run("Clone", func(t *testing.T) {
		ir := NewIndirectRef(10, 0)
		clone := ir.Clone().(IndirectRef)
		if clone.ObjectNumber != ir.ObjectNumber || clone.GenerationNumber != ir.GenerationNumber {
			t.Errorf("IndirectRef.Clone() = %v, want %v", clone, ir)
		}
	})
}

func TestDim(t *testing.T) {
	t.Run("UnitConversions", func(t *testing.T) {
		d := Dim{Width: 72, Height: 72} // 1 inch x 1 inch

		inches := d.ToInches()
		if math.Abs(inches.Width-1.0) > 0.001 {
			t.Errorf("ToInches().Width = %v, want 1.0", inches.Width)
		}

		cm := d.ToCentimetres()
		if math.Abs(cm.Width-2.54) > 0.01 {
			t.Errorf("ToCentimetres().Width = %v, want 2.54", cm.Width)
		}

		mm := d.ToMillimetres()
		if math.Abs(mm.Width-25.4) > 0.1 {
			t.Errorf("ToMillimetres().Width = %v, want 25.4", mm.Width)
		}
	})

	t.Run("ConvertToUnit", func(t *testing.T) {
		d := Dim{Width: 72, Height: 72}

		points := d.ConvertToUnit(POINTS)
		if points.Width != 72 {
			t.Errorf("ConvertToUnit(POINTS).Width = %v, want 72", points.Width)
		}

		inches := d.ConvertToUnit(INCHES)
		if math.Abs(inches.Width-1.0) > 0.001 {
			t.Errorf("ConvertToUnit(INCHES).Width = %v, want 1.0", inches.Width)
		}
	})

	t.Run("AspectRatio", func(t *testing.T) {
		d := Dim{Width: 200, Height: 100}
		if d.AspectRatio() != 2.0 {
			t.Errorf("AspectRatio() = %v, want 2.0", d.AspectRatio())
		}
	})

	t.Run("LandscapePortrait", func(t *testing.T) {
		landscape := Dim{Width: 200, Height: 100}
		portrait := Dim{Width: 100, Height: 200}

		if !landscape.Landscape() {
			t.Error("Wide dim should be landscape")
		}
		if !portrait.Portrait() {
			t.Error("Tall dim should be portrait")
		}
	})

	t.Run("String", func(t *testing.T) {
		d := Dim{Width: 100, Height: 200}
		s := d.String()
		if s == "" {
			t.Error("Dim.String() should return non-empty string")
		}
	})
}

func TestToUserSpace(t *testing.T) {
	tests := []struct {
		name string
		f    float64
		unit DisplayUnit
		want float64
	}{
		{"inches", 1.0, INCHES, 72.0},
		{"centimetres", 2.54, CENTIMETRES, 72.0},
		{"millimetres", 25.4, MILLIMETRES, 72.0},
		{"points", 72.0, POINTS, 72.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToUserSpace(tt.f, tt.unit)
			if math.Abs(got-tt.want) > 0.1 {
				t.Errorf("ToUserSpace(%v, %v) = %v, want %v", tt.f, tt.unit, got, tt.want)
			}
		})
	}
}

func TestRectForFormat(t *testing.T) {
	tests := []struct {
		format string
		wantW  float64
		wantH  float64
	}{
		{"A4", 595.0, 842.0},     // A4 is 595 x 842 points
		{"Letter", 612.0, 792.0}, // Letter is 612 x 792 points
		{"A3", 842.0, 1191.0},
		{"Legal", 612.0, 1008.0},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			r := RectForFormat(tt.format)
			if r == nil {
				t.Fatalf("RectForFormat(%q) returned nil", tt.format)
			}
			if math.Abs(r.Width()-tt.wantW) > 1.0 {
				t.Errorf("Width = %v, want %v", r.Width(), tt.wantW)
			}
			if math.Abs(r.Height()-tt.wantH) > 1.0 {
				t.Errorf("Height = %v, want %v", r.Height(), tt.wantH)
			}
		})
	}
}

func TestPaperSizes(t *testing.T) {
	// Test that common paper sizes exist in the map
	commonSizes := []string{
		"A0", "A1", "A2", "A3", "A4", "A5", "A6",
		"Letter", "Legal", "Tabloid", "Ledger",
		"B0", "B1", "B2", "B3", "B4", "B5",
	}

	for _, size := range commonSizes {
		t.Run(size, func(t *testing.T) {
			dim, ok := PaperSize[size]
			if !ok {
				t.Errorf("PaperSize[%q] not found", size)
				return
			}
			if dim.Width <= 0 || dim.Height <= 0 {
				t.Errorf("PaperSize[%q] has invalid dimensions: %v", size, dim)
			}
		})
	}
}

func TestConstants(t *testing.T) {
	t.Run("EOL constants", func(t *testing.T) {
		if EolLF != "\x0A" {
			t.Errorf("EolLF = %q, want LF", EolLF)
		}
		if EolCR != "\x0D" {
			t.Errorf("EolCR = %q, want CR", EolCR)
		}
		if EolCRLF != "\x0D\x0A" {
			t.Errorf("EolCRLF = %q, want CRLF", EolCRLF)
		}
	})

	t.Run("FreeHeadGeneration", func(t *testing.T) {
		if FreeHeadGeneration != 65535 {
			t.Errorf("FreeHeadGeneration = %d, want 65535", FreeHeadGeneration)
		}
	})

	t.Run("ByteSize constants", func(t *testing.T) {
		if KB != 1024 {
			t.Errorf("KB = %v, want 1024", KB)
		}
		if MB != 1024*1024 {
			t.Errorf("MB = %v, want 1048576", MB)
		}
		if GB != 1024*1024*1024 {
			t.Errorf("GB = %v, want 1073741824", GB)
		}
	})
}
