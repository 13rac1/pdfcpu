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

package scan

import (
	"bufio"
	"strings"
	"testing"
)

func TestLines(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "empty",
			input: "",
			want:  nil,
		},
		{
			name:  "single line no EOL",
			input: "hello",
			want:  []string{"hello"},
		},
		{
			name:  "single line LF",
			input: "hello\n",
			want:  []string{"hello"},
		},
		{
			name:  "single line CR",
			input: "hello\r",
			want:  []string{"hello"},
		},
		{
			name:  "single line CRLF",
			input: "hello\r\n",
			want:  []string{"hello"},
		},
		{
			name:  "two lines LF",
			input: "hello\nworld",
			want:  []string{"hello", "world"},
		},
		{
			name:  "two lines CR",
			input: "hello\rworld",
			want:  []string{"hello", "world"},
		},
		{
			name:  "two lines CRLF",
			input: "hello\r\nworld",
			want:  []string{"hello", "world"},
		},
		{
			name:  "mixed EOL",
			input: "one\ntwo\rthree\r\nfour",
			want:  []string{"one", "two", "three", "four"},
		},
		{
			name:  "empty lines",
			input: "a\n\nb",
			want:  []string{"a", "", "b"},
		},
		{
			name:  "CR before LF but not adjacent",
			input: "a\rb\nc",
			want:  []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := bufio.NewScanner(strings.NewReader(tt.input))
			scanner.Split(Lines)
			var got []string
			for scanner.Scan() {
				got = append(got, scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				t.Errorf("scanner error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Errorf("Lines got %d lines, want %d", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("line %d = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestLinesSingleEOL(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "empty",
			input: "",
			want:  nil,
		},
		{
			name:  "single line no EOL",
			input: "hello",
			want:  []string{"hello"},
		},
		{
			name:  "single line LF",
			input: "hello\n",
			want:  []string{"hello"},
		},
		{
			name:  "single line CR",
			input: "hello\r",
			want:  []string{"hello"},
		},
		{
			name:  "CRLF treated as CR then LF",
			input: "hello\r\n",
			want:  []string{"hello", ""},
		},
		{
			name:  "two lines LF",
			input: "hello\nworld",
			want:  []string{"hello", "world"},
		},
		{
			name:  "two lines CR",
			input: "hello\rworld",
			want:  []string{"hello", "world"},
		},
		{
			name:  "mixed EOL",
			input: "one\ntwo\rthree",
			want:  []string{"one", "two", "three"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := bufio.NewScanner(strings.NewReader(tt.input))
			scanner.Split(LinesSingleEOL)
			var got []string
			for scanner.Scan() {
				got = append(got, scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				t.Errorf("scanner error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Errorf("LinesSingleEOL got %d lines, want %d", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("line %d = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestLinesDirectCall(t *testing.T) {
	t.Run("atEOF with empty data", func(t *testing.T) {
		advance, token, err := Lines([]byte{}, true)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if advance != 0 {
			t.Errorf("advance = %d, want 0", advance)
		}
		if token != nil {
			t.Errorf("token = %v, want nil", token)
		}
	})

	t.Run("not atEOF with no EOL", func(t *testing.T) {
		advance, token, err := Lines([]byte("hello"), false)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if advance != 0 {
			t.Errorf("advance = %d, want 0", advance)
		}
		if token != nil {
			t.Errorf("token = %v, want nil", token)
		}
	})

	t.Run("atEOF with data no EOL", func(t *testing.T) {
		advance, token, err := Lines([]byte("hello"), true)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if advance != 5 {
			t.Errorf("advance = %d, want 5", advance)
		}
		if string(token) != "hello" {
			t.Errorf("token = %q, want %q", string(token), "hello")
		}
	})
}

func TestLinesSingleEOLDirectCall(t *testing.T) {
	t.Run("atEOF with empty data", func(t *testing.T) {
		advance, token, err := LinesSingleEOL([]byte{}, true)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if advance != 0 {
			t.Errorf("advance = %d, want 0", advance)
		}
		if token != nil {
			t.Errorf("token = %v, want nil", token)
		}
	})

	t.Run("not atEOF with no EOL", func(t *testing.T) {
		advance, token, err := LinesSingleEOL([]byte("hello"), false)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if advance != 0 {
			t.Errorf("advance = %d, want 0", advance)
		}
		if token != nil {
			t.Errorf("token = %v, want nil", token)
		}
	})

	t.Run("atEOF with data no EOL", func(t *testing.T) {
		advance, token, err := LinesSingleEOL([]byte("hello"), true)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if advance != 5 {
			t.Errorf("advance = %d, want 5", advance)
		}
		if string(token) != "hello" {
			t.Errorf("token = %q, want %q", string(token), "hello")
		}
	})
}
