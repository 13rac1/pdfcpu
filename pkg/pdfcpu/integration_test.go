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

package pdfcpu_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/validate"
)

var (
	testdataDir = filepath.Join("..", "testdata")
)

func getTmpDir(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "pdfcpu_integration")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(tmpDir) })
	return tmpDir
}

func TestReadValidatePDFs(t *testing.T) {
	tests := []struct {
		name     string
		filename string
	}{
		{"image PDF", "testImage.pdf"},
		{"programming book", "TheGoProgrammingLanguageCh1.pdf"},
		{"FOSDEM presentation", "FOSDEM14_HPC_devroom_14_GoCUDA.pdf"},
		{"annotations", "annotTest.pdf"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inFile := filepath.Join(testdataDir, tt.filename)

			// Read PDF file
			ctx, err := pdfcpu.ReadFile(inFile, model.NewDefaultConfiguration())
			if err != nil {
				t.Fatalf("ReadFile(%q) error = %v", tt.filename, err)
			}

			// Verify context is valid
			if ctx == nil {
				t.Fatal("ReadFile returned nil context")
			}

			if ctx.XRefTable == nil {
				t.Fatal("Context has nil XRefTable")
			}

			// Ensure page count is populated
			if err := ctx.XRefTable.EnsurePageCount(); err != nil {
				t.Fatalf("EnsurePageCount() error = %v", err)
			}

			// Validate XRefTable
			if err := validate.XRefTable(ctx); err != nil {
				t.Errorf("validate.XRefTable() error = %v", err)
			}

			// PageCount should be positive for real PDFs
			if ctx.PageCount <= 0 {
				t.Errorf("ReadFile(%q) PageCount = %d, want > 0", tt.filename, ctx.PageCount)
			}
		})
	}
}

func TestReadWriteRoundtrip(t *testing.T) {
	tests := []struct {
		name     string
		filename string
	}{
		{"image PDF", "testImage.pdf"},
		{"programming book", "TheGoProgrammingLanguageCh1.pdf"},
	}

	tmpDir := getTmpDir(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inFile := filepath.Join(testdataDir, tt.filename)
			outFile := filepath.Join(tmpDir, "roundtrip_"+tt.filename)

			// Read original PDF
			ctx1, err := pdfcpu.ReadFile(inFile, model.NewDefaultConfiguration())
			if err != nil {
				t.Fatalf("ReadFile(%q) error = %v", tt.filename, err)
			}
			originalPageCount := ctx1.PageCount

			// Set output file
			ctx1.Write.DirName = filepath.Dir(outFile)
			ctx1.Write.FileName = filepath.Base(outFile)

			// Write to temp file
			if err := pdfcpu.WriteContext(ctx1); err != nil {
				t.Fatalf("WriteContext() error = %v", err)
			}

			// Read back the written PDF
			ctx2, err := pdfcpu.ReadFile(outFile, model.NewDefaultConfiguration())
			if err != nil {
				t.Fatalf("ReadFile(written file) error = %v", err)
			}

			// Verify page counts match
			if ctx2.PageCount != originalPageCount {
				t.Errorf("Roundtrip PageCount = %d, want %d", ctx2.PageCount, originalPageCount)
			}
		})
	}
}

func TestOptimizeContext(t *testing.T) {
	tests := []struct {
		name     string
		filename string
	}{
		{"image PDF", "testImage.pdf"},
		{"programming book", "TheGoProgrammingLanguageCh1.pdf"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inFile := filepath.Join(testdataDir, tt.filename)

			// Read PDF file
			ctx, err := pdfcpu.ReadFile(inFile, model.NewDefaultConfiguration())
			if err != nil {
				t.Fatalf("ReadFile(%q) error = %v", tt.filename, err)
			}

			// Ensure page count is populated
			if err := ctx.XRefTable.EnsurePageCount(); err != nil {
				t.Fatalf("EnsurePageCount() error = %v", err)
			}

			if ctx.PageCount == 0 {
				t.Skipf("Skipping optimization - PDF has no pages")
			}

			// Optimize XRefTable - verify it doesn't error
			if err := pdfcpu.OptimizeXRefTable(ctx); err != nil {
				t.Fatalf("OptimizeXRefTable() error = %v", err)
			}

			// Optimization completed successfully (ctx.Optimized flag may or may not be set
			// depending on whether any optimizations were actually performed)
		})
	}
}

func TestRotatePages(t *testing.T) {
	tmpDir := getTmpDir(t)
	inFile := filepath.Join(testdataDir, "TheGoProgrammingLanguageCh1.pdf")
	outFile := filepath.Join(tmpDir, "rotated.pdf")

	// Read PDF file
	ctx, err := pdfcpu.ReadFile(inFile, model.NewDefaultConfiguration())
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	// Ensure page count is populated
	if err := ctx.XRefTable.EnsurePageCount(); err != nil {
		t.Fatalf("EnsurePageCount() error = %v", err)
	}

	if ctx.PageCount == 0 {
		t.Skip("Skipping test - PDF has no pages")
	}

	// Rotate first page 90 degrees
	selectedPages := make(map[int]bool)
	selectedPages[1] = true
	if err := pdfcpu.RotatePages(ctx, selectedPages, 90); err != nil {
		t.Fatalf("RotatePages() error = %v", err)
	}

	// Set output file
	ctx.Write.DirName = filepath.Dir(outFile)
	ctx.Write.FileName = filepath.Base(outFile)

	// Write rotated PDF
	if err := pdfcpu.WriteContext(ctx); err != nil {
		t.Fatalf("WriteContext() error = %v", err)
	}

	// Read back the rotated PDF to verify it's valid
	ctx2, err := pdfcpu.ReadFile(outFile, model.NewDefaultConfiguration())
	if err != nil {
		t.Fatalf("ReadFile(rotated) error = %v", err)
	}

	// Ensure page count is the same
	if err := ctx2.XRefTable.EnsurePageCount(); err != nil {
		t.Fatalf("EnsurePageCount() on rotated PDF error = %v", err)
	}

	if ctx2.PageCount != ctx.PageCount {
		t.Errorf("Rotated PDF PageCount = %d, want %d", ctx2.PageCount, ctx.PageCount)
	}
}

func TestEmptyPDF(t *testing.T) {
	inFile := filepath.Join(testdataDir, "empty.pdf")

	// Read empty PDF file
	ctx, err := pdfcpu.ReadFile(inFile, model.NewDefaultConfiguration())
	if err != nil {
		t.Fatalf("ReadFile(empty.pdf) error = %v", err)
	}

	// Verify it reads without error (PageCount can be 0 for empty PDFs)
	if ctx == nil {
		t.Error("ReadFile(empty.pdf) returned nil context")
	}

	// Empty PDF should have PageCount >= 0
	if ctx.PageCount < 0 {
		t.Errorf("empty.pdf PageCount = %d, want >= 0", ctx.PageCount)
	}
}
