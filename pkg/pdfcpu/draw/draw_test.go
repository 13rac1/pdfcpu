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

package draw

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/color"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

func TestRenderModeConstants(t *testing.T) {
	if RMFill != 0 {
		t.Errorf("RMFill = %d, want 0", RMFill)
	}
	if RMStroke != 1 {
		t.Errorf("RMStroke = %d, want 1", RMStroke)
	}
	if RMFillAndStroke != 2 {
		t.Errorf("RMFillAndStroke = %d, want 2", RMFillAndStroke)
	}
}

func TestSetLineJoinStyle(t *testing.T) {
	var buf bytes.Buffer
	SetLineJoinStyle(&buf, types.LJMiter)
	got := buf.String()
	if !strings.Contains(got, "j") {
		t.Errorf("SetLineJoinStyle output = %q, should contain 'j'", got)
	}
}

func TestSetLineWidth(t *testing.T) {
	var buf bytes.Buffer
	SetLineWidth(&buf, 2.5)
	got := buf.String()
	if !strings.Contains(got, "2.50") || !strings.Contains(got, "w") {
		t.Errorf("SetLineWidth output = %q, should contain '2.50 w'", got)
	}
}

func TestSetStrokeColor(t *testing.T) {
	var buf bytes.Buffer
	SetStrokeColor(&buf, color.Red)
	got := buf.String()
	if !strings.Contains(got, "RG") {
		t.Errorf("SetStrokeColor output = %q, should contain 'RG'", got)
	}
}

func TestSetFillColor(t *testing.T) {
	var buf bytes.Buffer
	SetFillColor(&buf, color.Blue)
	got := buf.String()
	if !strings.Contains(got, "rg") {
		t.Errorf("SetFillColor output = %q, should contain 'rg'", got)
	}
}

func TestDrawLineSimple(t *testing.T) {
	var buf bytes.Buffer
	DrawLineSimple(&buf, 0, 0, 100, 100)
	got := buf.String()
	if !strings.Contains(got, "m") || !strings.Contains(got, "l") || !strings.Contains(got, "s") {
		t.Errorf("DrawLineSimple output = %q, should contain 'm', 'l', 's'", got)
	}
}

func TestDrawLine(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		var buf bytes.Buffer
		DrawLine(&buf, 0, 0, 100, 100, 1.0, nil, nil)
		got := buf.String()
		if !strings.HasPrefix(got, "q") || !strings.HasSuffix(strings.TrimSpace(got), "Q") {
			t.Errorf("DrawLine output = %q, should be wrapped in q/Q", got)
		}
	})

	t.Run("with stroke color", func(t *testing.T) {
		var buf bytes.Buffer
		red := color.Red
		DrawLine(&buf, 0, 0, 100, 100, 1.0, &red, nil)
		got := buf.String()
		if !strings.Contains(got, "RG") {
			t.Errorf("DrawLine with color output = %q, should contain 'RG'", got)
		}
	})

	t.Run("with style", func(t *testing.T) {
		var buf bytes.Buffer
		style := types.LJRound
		DrawLine(&buf, 0, 0, 100, 100, 1.0, nil, &style)
		got := buf.String()
		if !strings.Contains(got, "j") {
			t.Errorf("DrawLine with style output = %q, should contain 'j'", got)
		}
	})
}

func TestDrawRectSimple(t *testing.T) {
	var buf bytes.Buffer
	r := types.NewRectangle(0, 0, 100, 100)
	DrawRectSimple(&buf, r)
	got := buf.String()
	if !strings.Contains(got, "re") || !strings.Contains(got, "s") {
		t.Errorf("DrawRectSimple output = %q, should contain 're' and 's'", got)
	}
}

func TestDrawRect(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		var buf bytes.Buffer
		r := types.NewRectangle(0, 0, 100, 100)
		DrawRect(&buf, r, 1.0, nil, nil)
		got := buf.String()
		if !strings.HasPrefix(got, "q") {
			t.Errorf("DrawRect output = %q, should start with 'q'", got)
		}
	})

	t.Run("with all options", func(t *testing.T) {
		var buf bytes.Buffer
		r := types.NewRectangle(0, 0, 100, 100)
		red := color.Red
		style := types.LJBevel
		DrawRect(&buf, r, 2.0, &red, &style)
		got := buf.String()
		if !strings.Contains(got, "RG") || !strings.Contains(got, "j") {
			t.Errorf("DrawRect with options output = %q", got)
		}
	})
}

func TestFillRect(t *testing.T) {
	t.Run("with stroke color", func(t *testing.T) {
		var buf bytes.Buffer
		r := types.NewRectangle(0, 0, 100, 100)
		stroke := color.Black
		FillRect(&buf, r, 1.0, &stroke, color.White, nil)
		got := buf.String()
		if !strings.Contains(got, "B") {
			t.Errorf("FillRect output = %q, should contain 'B'", got)
		}
	})

	t.Run("without stroke color", func(t *testing.T) {
		var buf bytes.Buffer
		r := types.NewRectangle(0, 0, 100, 100)
		FillRect(&buf, r, 1.0, nil, color.Gray, nil)
		got := buf.String()
		if !strings.Contains(got, "rg") {
			t.Errorf("FillRect output = %q, should contain 'rg'", got)
		}
	})

	t.Run("with style", func(t *testing.T) {
		var buf bytes.Buffer
		r := types.NewRectangle(0, 0, 100, 100)
		style := types.LJMiter
		FillRect(&buf, r, 1.0, nil, color.Gray, &style)
		got := buf.String()
		if !strings.Contains(got, "j") {
			t.Errorf("FillRect with style output = %q, should contain 'j'", got)
		}
	})
}

func TestDrawCircle(t *testing.T) {
	t.Run("stroke only", func(t *testing.T) {
		var buf bytes.Buffer
		DrawCircle(&buf, 50, 50, 25, color.Black, nil)
		got := buf.String()
		if !strings.Contains(got, "c") || !strings.Contains(got, "s") {
			t.Errorf("DrawCircle output = %q, should contain curves", got)
		}
	})

	t.Run("with fill", func(t *testing.T) {
		var buf bytes.Buffer
		fill := color.Red
		DrawCircle(&buf, 50, 50, 25, color.Black, &fill)
		got := buf.String()
		if !strings.Contains(got, "f") {
			t.Errorf("DrawCircle with fill output = %q, should contain 'f'", got)
		}
	})
}

func TestFillRectNoBorder(t *testing.T) {
	var buf bytes.Buffer
	r := types.NewRectangle(0, 0, 100, 100)
	FillRectNoBorder(&buf, r, color.Blue)
	got := buf.String()
	if !strings.Contains(got, "B") {
		t.Errorf("FillRectNoBorder output = %q, should contain 'B'", got)
	}
}

func TestDrawGrid(t *testing.T) {
	t.Run("without fill", func(t *testing.T) {
		var buf bytes.Buffer
		r := types.NewRectangle(0, 0, 100, 100)
		DrawGrid(&buf, 2, 2, r, color.Black, nil)
		got := buf.String()
		// Should have lines
		if !strings.Contains(got, "l") {
			t.Errorf("DrawGrid output = %q, should contain lines", got)
		}
	})

	t.Run("with fill", func(t *testing.T) {
		var buf bytes.Buffer
		r := types.NewRectangle(0, 0, 100, 100)
		fill := color.LightGray
		DrawGrid(&buf, 3, 3, r, color.Black, &fill)
		got := buf.String()
		if !strings.Contains(got, "rg") {
			t.Errorf("DrawGrid with fill output = %q, should contain fill color", got)
		}
	})
}

func TestDrawHairCross(t *testing.T) {
	t.Run("at origin", func(t *testing.T) {
		var buf bytes.Buffer
		r := types.NewRectangle(0, 0, 100, 100)
		DrawHairCross(&buf, 0, 0, r)
		got := buf.String()
		// Should have two lines
		if strings.Count(got, "q") < 2 {
			t.Errorf("DrawHairCross output = %q, should have 2 lines", got)
		}
	})

	t.Run("at specific point", func(t *testing.T) {
		var buf bytes.Buffer
		r := types.NewRectangle(0, 0, 100, 100)
		DrawHairCross(&buf, 25, 75, r)
		got := buf.String()
		if got == "" {
			t.Error("DrawHairCross should produce output")
		}
	})
}

func TestHorSepLine(t *testing.T) {
	t.Run("single column", func(t *testing.T) {
		got := HorSepLine([]int{10})
		// HBar is a multi-byte Unicode character, so check rune count instead
		runeCount := len([]rune(got))
		if runeCount != 10 {
			t.Errorf("HorSepLine([10]) rune count = %d, want 10", runeCount)
		}
	})

	t.Run("multiple columns", func(t *testing.T) {
		got := HorSepLine([]int{5, 10, 5})
		if !strings.Contains(got, CrossBar) {
			t.Errorf("HorSepLine with multiple columns should contain CrossBar")
		}
	})

	t.Run("empty", func(t *testing.T) {
		got := HorSepLine([]int{})
		if got != "" {
			t.Errorf("HorSepLine([]) = %q, want empty", got)
		}
	})
}

func TestBarConstants(t *testing.T) {
	if HBar == "" {
		t.Error("HBar should not be empty")
	}
	if VBar == "" {
		t.Error("VBar should not be empty")
	}
	if CrossBar == "" {
		t.Error("CrossBar should not be empty")
	}
}
