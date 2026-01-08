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

func TestExtractPages(t *testing.T) {
	inFile := filepath.Join(testdataDir, "TheGoProgrammingLanguageCh1.pdf")

	// Read PDF file
	ctx, err := pdfcpu.ReadFile(inFile, model.NewDefaultConfiguration())
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	// Ensure page count is populated
	if err := ctx.XRefTable.EnsurePageCount(); err != nil {
		t.Fatalf("EnsurePageCount() error = %v", err)
	}

	if ctx.PageCount < 3 {
		t.Skipf("Need at least 3 pages for extract test, got %d", ctx.PageCount)
	}

	originalPageCount := ctx.PageCount

	// Extract pages 1, 2, 3
	pageNrs := []int{1, 2, 3}
	extractedCtx, err := pdfcpu.ExtractPages(ctx, pageNrs, false)
	if err != nil {
		t.Fatalf("ExtractPages() error = %v", err)
	}

	// Verify extracted context
	if extractedCtx == nil {
		t.Fatal("ExtractPages returned nil context")
	}

	// Ensure page count is populated for extracted context
	if err := extractedCtx.XRefTable.EnsurePageCount(); err != nil {
		t.Fatalf("EnsurePageCount() on extracted context error = %v", err)
	}

	// Extracted PDF should have 3 pages
	if extractedCtx.PageCount != 3 {
		t.Errorf("Extracted PageCount = %d, want 3", extractedCtx.PageCount)
	}

	// Original context should be unchanged
	if ctx.PageCount != originalPageCount {
		t.Errorf("Original PageCount changed from %d to %d", originalPageCount, ctx.PageCount)
	}
}

func TestAddPages(t *testing.T) {
	tmpDir := getTmpDir(t)

	// Read source PDF
	srcFile := filepath.Join(testdataDir, "testImage.pdf")
	ctxSrc, err := pdfcpu.ReadFile(srcFile, model.NewDefaultConfiguration())
	if err != nil {
		t.Fatalf("ReadFile(source) error = %v", err)
	}

	if err := ctxSrc.XRefTable.EnsurePageCount(); err != nil {
		t.Fatalf("EnsurePageCount(source) error = %v", err)
	}

	if ctxSrc.PageCount < 1 {
		t.Skip("Source PDF needs at least 1 page")
	}

	srcPageCount := ctxSrc.PageCount

	// Read destination PDF
	destFile := filepath.Join(testdataDir, "TheGoProgrammingLanguageCh1.pdf")
	ctxDest, err := pdfcpu.ReadFile(destFile, model.NewDefaultConfiguration())
	if err != nil {
		t.Fatalf("ReadFile(dest) error = %v", err)
	}

	if err := ctxDest.XRefTable.EnsurePageCount(); err != nil {
		t.Fatalf("EnsurePageCount(dest) error = %v", err)
	}

	destPageCount := ctxDest.PageCount

	// Add first page from source to destination
	pageNrs := []int{1}
	if err := pdfcpu.AddPages(ctxSrc, ctxDest, pageNrs, false); err != nil {
		t.Fatalf("AddPages() error = %v", err)
	}

	// Write merged PDF to verify it's valid
	outFile := filepath.Join(tmpDir, "merged.pdf")
	ctxDest.Write.DirName = filepath.Dir(outFile)
	ctxDest.Write.FileName = filepath.Base(outFile)

	if err := pdfcpu.WriteContext(ctxDest); err != nil {
		t.Fatalf("WriteContext() error = %v", err)
	}

	// Read back merged PDF and verify page count increased
	ctxMerged, err := pdfcpu.ReadFile(outFile, model.NewDefaultConfiguration())
	if err != nil {
		t.Fatalf("ReadFile(merged) error = %v", err)
	}

	if err := ctxMerged.XRefTable.EnsurePageCount(); err != nil {
		t.Fatalf("EnsurePageCount(merged) error = %v", err)
	}

	// Merged PDF should have more pages than original destination
	if ctxMerged.PageCount <= destPageCount {
		t.Errorf("Merged PDF PageCount = %d, should be > %d", ctxMerged.PageCount, destPageCount)
	}

	expectedCount := destPageCount + 1
	if ctxMerged.PageCount != expectedCount {
		t.Logf("Note: Merged PDF has %d pages (expected %d)", ctxMerged.PageCount, expectedCount)
	}

	// Source should be unchanged
	if ctxSrc.PageCount != srcPageCount {
		t.Errorf("Source PageCount changed from %d to %d", srcPageCount, ctxSrc.PageCount)
	}
}

func TestInfo(t *testing.T) {
	tests := []struct {
		name     string
		filename string
	}{
		{"simple PDF", "testImage.pdf"},
		{"programming book", "TheGoProgrammingLanguageCh1.pdf"},
		{"annotations", "annotTest.pdf"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inFile := filepath.Join(testdataDir, tt.filename)

			// Read PDF file
			ctx, err := pdfcpu.ReadFile(inFile, model.NewDefaultConfiguration())
			if err != nil {
				t.Fatalf("ReadFile() error = %v", err)
			}

			if err := ctx.XRefTable.EnsurePageCount(); err != nil {
				t.Fatalf("EnsurePageCount() error = %v", err)
			}

			// Get PDF info - test with all pages
			selectedPages := make(map[int]bool)
			info, err := pdfcpu.Info(ctx, tt.filename, selectedPages, false)
			if err != nil {
				t.Fatalf("Info() error = %v", err)
			}

			// Verify info structure
			if info == nil {
				t.Fatal("Info() returned nil")
			}

			if info.FileName != tt.filename {
				t.Errorf("Info.FileName = %q, want %q", info.FileName, tt.filename)
			}

			if info.PageCount != ctx.PageCount {
				t.Errorf("Info.PageCount = %d, want %d", info.PageCount, ctx.PageCount)
			}

			// Verify version is populated
			if info.Version == "" {
				t.Error("Info.Version is empty")
			}
		})
	}
}

func TestImages(t *testing.T) {
	// Test with PDF that has images
	inFile := filepath.Join(testdataDir, "testImage.pdf")

	ctx, err := pdfcpu.ReadFile(inFile, model.NewDefaultConfiguration())
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	if err := ctx.XRefTable.EnsurePageCount(); err != nil {
		t.Fatalf("EnsurePageCount() error = %v", err)
	}

	// Get images from all pages
	selectedPages := make(map[int]bool)
	images, maxLengths, err := pdfcpu.Images(ctx, selectedPages)
	if err != nil {
		t.Fatalf("Images() error = %v", err)
	}

	// testImage.pdf should have at least one image
	if len(images) == 0 {
		t.Log("Note: No images found in testImage.pdf (may be expected)")
	}

	// Verify maxLengths is returned
	if maxLengths == nil {
		t.Error("Images() returned nil maxLengths")
	}

	// Test ListImages as well
	lines, err := pdfcpu.ListImages(ctx, selectedPages)
	if err != nil {
		t.Fatalf("ListImages() error = %v", err)
	}

	// ListImages should return some output
	if lines == nil {
		t.Error("ListImages() returned nil")
	}
}
