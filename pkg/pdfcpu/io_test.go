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

package pdfcpu

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestWrite(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("write new file", func(t *testing.T) {
		content := []byte("test content")
		path := filepath.Join(tmpDir, "new_file.txt")

		written, err := Write(bytes.NewReader(content), path, false)
		if err != nil {
			t.Fatalf("Write() error = %v", err)
		}
		if !written {
			t.Error("Write() returned false for new file")
		}

		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("ReadFile() error = %v", err)
		}
		if string(data) != "test content" {
			t.Errorf("File content = %q, want %q", data, "test content")
		}
	})

	t.Run("overwrite existing file", func(t *testing.T) {
		path := filepath.Join(tmpDir, "existing.txt")
		os.WriteFile(path, []byte("old content"), 0644)

		written, err := Write(bytes.NewReader([]byte("new content")), path, true)
		if err != nil {
			t.Fatalf("Write() error = %v", err)
		}
		if !written {
			t.Error("Write() returned false with overwrite=true")
		}

		data, _ := os.ReadFile(path)
		if string(data) != "new content" {
			t.Errorf("File content = %q, want %q", data, "new content")
		}
	})

	t.Run("skip existing file when overwrite false", func(t *testing.T) {
		path := filepath.Join(tmpDir, "skip.txt")
		os.WriteFile(path, []byte("original"), 0644)

		written, err := Write(bytes.NewReader([]byte("ignored")), path, false)
		if err != nil {
			t.Fatalf("Write() error = %v", err)
		}
		if written {
			t.Error("Write() returned true when should skip")
		}

		data, _ := os.ReadFile(path)
		if string(data) != "original" {
			t.Errorf("File content = %q, want %q", data, "original")
		}
	})

	t.Run("error on invalid path", func(t *testing.T) {
		path := filepath.Join(tmpDir, "nonexistent", "subdir", "file.txt")
		_, err := Write(bytes.NewReader([]byte("content")), path, false)
		if err == nil {
			t.Error("Write() expected error for invalid path")
		}
	})
}

func TestCopyFile(t *testing.T) {
	tmpDir := t.TempDir()

	srcPath := filepath.Join(tmpDir, "source.txt")
	os.WriteFile(srcPath, []byte("source content"), 0644)

	t.Run("copy to new file", func(t *testing.T) {
		destPath := filepath.Join(tmpDir, "dest_new.txt")

		copied, err := CopyFile(srcPath, destPath, false)
		if err != nil {
			t.Fatalf("CopyFile() error = %v", err)
		}
		if !copied {
			t.Error("CopyFile() returned false for new dest")
		}

		data, _ := os.ReadFile(destPath)
		if string(data) != "source content" {
			t.Errorf("Dest content = %q, want %q", data, "source content")
		}
	})

	t.Run("overwrite existing dest", func(t *testing.T) {
		destPath := filepath.Join(tmpDir, "dest_overwrite.txt")
		os.WriteFile(destPath, []byte("old dest"), 0644)

		copied, err := CopyFile(srcPath, destPath, true)
		if err != nil {
			t.Fatalf("CopyFile() error = %v", err)
		}
		if !copied {
			t.Error("CopyFile() returned false with overwrite=true")
		}

		data, _ := os.ReadFile(destPath)
		if string(data) != "source content" {
			t.Errorf("Dest content = %q, want %q", data, "source content")
		}
	})

	t.Run("skip existing dest when overwrite false", func(t *testing.T) {
		destPath := filepath.Join(tmpDir, "dest_skip.txt")
		os.WriteFile(destPath, []byte("keep this"), 0644)

		copied, err := CopyFile(srcPath, destPath, false)
		if err != nil {
			t.Fatalf("CopyFile() error = %v", err)
		}
		if copied {
			t.Error("CopyFile() returned true when should skip")
		}

		data, _ := os.ReadFile(destPath)
		if string(data) != "keep this" {
			t.Errorf("Dest content = %q, want %q", data, "keep this")
		}
	})

	t.Run("error on non-existent source", func(t *testing.T) {
		_, err := CopyFile(filepath.Join(tmpDir, "nonexistent.txt"), filepath.Join(tmpDir, "out.txt"), false)
		if err == nil {
			t.Error("CopyFile() expected error for non-existent source")
		}
	})

	t.Run("error on invalid dest path", func(t *testing.T) {
		destPath := filepath.Join(tmpDir, "nonexistent", "subdir", "file.txt")
		_, err := CopyFile(srcPath, destPath, false)
		if err == nil {
			t.Error("CopyFile() expected error for invalid dest path")
		}
	})
}
