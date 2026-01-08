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

package model_test

import (
	"path/filepath"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

var testdataDir = filepath.Join("..", "..", "testdata")

// Helper to load a test PDF
func loadTestPDF(t *testing.T, filename string) *model.Context {
	t.Helper()
	inFile := filepath.Join(testdataDir, filename)
	ctx, err := pdfcpu.ReadFile(inFile, model.NewDefaultConfiguration())
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", filename, err)
	}
	if err := ctx.XRefTable.EnsurePageCount(); err != nil {
		t.Fatalf("EnsurePageCount() error = %v", err)
	}
	return ctx
}

func TestXRefTableVersion(t *testing.T) {
	tests := []struct {
		name    string
		pdfFile string
	}{
		{"simple PDF", "testImage.pdf"},
		{"programming book", "TheGoProgrammingLanguageCh1.pdf"},
		{"empty PDF", "empty.pdf"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := loadTestPDF(t, tt.pdfFile)

			// Test Version()
			version := ctx.XRefTable.Version()
			if version == model.V10 {
				t.Error("Version() returned V10, expected higher version")
			}

			// Test VersionString()
			versionStr := ctx.XRefTable.VersionString()
			if versionStr == "" {
				t.Error("VersionString() returned empty string")
			}
			if len(versionStr) < 3 {
				t.Errorf("VersionString() = %q, expected format like '1.4' or '1.7'", versionStr)
			}

			t.Logf("PDF %s: version = %s", tt.pdfFile, versionStr)
		})
	}
}

func TestXRefTableCatalog(t *testing.T) {
	tests := []struct {
		name    string
		pdfFile string
	}{
		{"simple PDF", "testImage.pdf"},
		{"book PDF", "TheGoProgrammingLanguageCh1.pdf"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := loadTestPDF(t, tt.pdfFile)

			catalog, err := ctx.XRefTable.Catalog()
			if err != nil {
				t.Fatalf("Catalog() error = %v", err)
			}

			if catalog == nil {
				t.Fatal("Catalog() returned nil")
			}

			// Verify Type is "Catalog"
			typeEntry := catalog.NameEntry("Type")
			if typeEntry == nil || *typeEntry != "Catalog" {
				t.Errorf("Catalog Type = %v, want 'Catalog'", typeEntry)
			}

			// Catalog should have Pages entry
			if _, found := catalog.Find("Pages"); !found {
				t.Error("Catalog missing 'Pages' entry")
			}
		})
	}
}

func TestXRefTablePages(t *testing.T) {
	tests := []struct {
		name    string
		pdfFile string
	}{
		{"simple PDF", "testImage.pdf"},
		{"book PDF", "TheGoProgrammingLanguageCh1.pdf"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := loadTestPDF(t, tt.pdfFile)

			pagesRef, err := ctx.XRefTable.Pages()
			if err != nil {
				t.Fatalf("Pages() error = %v", err)
			}

			if pagesRef == nil {
				t.Fatal("Pages() returned nil")
			}

			// Verify the reference is valid
			if pagesRef.ObjectNumber.Value() <= 0 {
				t.Errorf("Pages() returned invalid object number: %d", pagesRef.ObjectNumber.Value())
			}
		})
	}
}

func TestXRefTablePageDict(t *testing.T) {
	tests := []struct {
		name      string
		pdfFile   string
		pageNr    int
		expectErr bool
	}{
		{"first page", "testImage.pdf", 1, false},
		{"book first page", "TheGoProgrammingLanguageCh1.pdf", 1, false},
		{"invalid page zero", "testImage.pdf", 0, true},
		{"invalid negative page", "testImage.pdf", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := loadTestPDF(t, tt.pdfFile)

			pageDict, pageIndRef, inhPAttrs, err := ctx.XRefTable.PageDict(tt.pageNr, false)

			if tt.expectErr {
				if err == nil {
					t.Error("PageDict() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("PageDict() error = %v", err)
			}

			if pageDict == nil {
				t.Fatal("PageDict() returned nil dict")
			}

			if pageIndRef == nil {
				t.Fatal("PageDict() returned nil indirect reference")
			}

			if inhPAttrs == nil {
				t.Fatal("PageDict() returned nil inherited attributes")
			}

			// Verify Type is "Page"
			typeEntry := pageDict.NameEntry("Type")
			if typeEntry != nil && *typeEntry != "Page" {
				t.Errorf("Page Type = %v, want 'Page'", typeEntry)
			}
		})
	}
}

func TestXRefTablePageBoundaries(t *testing.T) {
	tests := []struct {
		name    string
		pdfFile string
	}{
		{"simple PDF", "testImage.pdf"},
		{"book PDF", "TheGoProgrammingLanguageCh1.pdf"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := loadTestPDF(t, tt.pdfFile)

			// Get all page boundaries (nil = all pages)
			boundaries, err := ctx.XRefTable.PageBoundaries(nil)
			if err != nil {
				t.Fatalf("PageBoundaries() error = %v", err)
			}

			if len(boundaries) == 0 {
				t.Error("PageBoundaries() returned empty slice")
			}

			if len(boundaries) != ctx.PageCount {
				t.Errorf("PageBoundaries() length = %d, want %d", len(boundaries), ctx.PageCount)
			}

			// Verify each boundary has at least a MediaBox
			for i, pb := range boundaries {
				if pb.Media == nil {
					t.Errorf("PageBoundaries[%d].Media is nil", i)
				}
			}
		})
	}
}

func TestXRefTablePageDims(t *testing.T) {
	tests := []struct {
		name    string
		pdfFile string
	}{
		{"simple PDF", "testImage.pdf"},
		{"book PDF", "TheGoProgrammingLanguageCh1.pdf"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := loadTestPDF(t, tt.pdfFile)

			dims, err := ctx.XRefTable.PageDims()
			if err != nil {
				t.Fatalf("PageDims() error = %v", err)
			}

			if len(dims) != ctx.PageCount {
				t.Errorf("PageDims() length = %d, want %d", len(dims), ctx.PageCount)
			}

			// Verify dimensions are positive
			for i, dim := range dims {
				if dim.Width <= 0 || dim.Height <= 0 {
					t.Errorf("PageDims[%d] = {Width: %f, Height: %f}, want positive values", i, dim.Width, dim.Height)
				}
			}
		})
	}
}

func TestXRefTableIsValid(t *testing.T) {
	ctx := loadTestPDF(t, "testImage.pdf")

	tests := []struct {
		name      string
		objNr     int
		genNr     int
		expectErr bool
	}{
		{"existing object 2", 2, 0, false},
		{"existing object 3", 3, 0, false},
		{"invalid object 0", 0, 0, true},
		{"invalid negative object", -1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ctx.XRefTable.IsObjValid(tt.objNr, tt.genNr)

			if tt.expectErr {
				if err == nil {
					t.Error("IsObjValid() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("IsObjValid() error = %v", err)
			}

			// Note: We don't check the returned value since Valid flag
			// is only set during validation, not during normal PDF reading
		})
	}
}

func TestXRefTableIsValidIndirectRef(t *testing.T) {
	ctx := loadTestPDF(t, "testImage.pdf")

	// Get the catalog to find a known indirect reference
	catalog, err := ctx.XRefTable.Catalog()
	if err != nil {
		t.Fatalf("Catalog() error = %v", err)
	}

	// Get Pages indirect reference from catalog
	pagesObj, found := catalog.Find("Pages")
	if !found {
		t.Fatal("Catalog missing Pages entry")
	}

	pagesRef, ok := pagesObj.(types.IndirectRef)
	if !ok {
		t.Fatalf("Pages entry is not an IndirectRef, got %T", pagesObj)
	}

	// Test IsValid with the pages reference - should not error
	_, err = ctx.XRefTable.IsValid(pagesRef)
	if err != nil {
		t.Fatalf("IsValid() error = %v", err)
	}

	// Note: We don't check the returned value since Valid flag
	// is only set during validation, not during normal PDF reading
}

func TestXRefTableDereference(t *testing.T) {
	ctx := loadTestPDF(t, "testImage.pdf")

	// Get catalog which contains indirect references
	catalog, err := ctx.XRefTable.Catalog()
	if err != nil {
		t.Fatalf("Catalog() error = %v", err)
	}

	// Get Pages entry (should be an indirect reference)
	pagesObj, found := catalog.Find("Pages")
	if !found {
		t.Fatal("Catalog missing Pages entry")
	}

	// Dereference it
	dereferenced, err := ctx.XRefTable.Dereference(pagesObj)
	if err != nil {
		t.Fatalf("Dereference() error = %v", err)
	}

	if dereferenced == nil {
		t.Fatal("Dereference() returned nil")
	}

	// The dereferenced object should be a Dict
	_, isDict := dereferenced.(types.Dict)
	if !isDict {
		t.Errorf("Dereferenced Pages = %T, expected types.Dict", dereferenced)
	}
}

func TestXRefTableDereferenceDict(t *testing.T) {
	ctx := loadTestPDF(t, "testImage.pdf")

	// Get catalog
	catalog, err := ctx.XRefTable.Catalog()
	if err != nil {
		t.Fatalf("Catalog() error = %v", err)
	}

	// Get Pages entry
	pagesObj, found := catalog.Find("Pages")
	if !found {
		t.Fatal("Catalog missing Pages entry")
	}

	// DereferenceDict should work since Pages is a Dict
	pagesDict, err := ctx.XRefTable.DereferenceDict(pagesObj)
	if err != nil {
		t.Fatalf("DereferenceDict() error = %v", err)
	}

	if pagesDict == nil {
		t.Fatal("DereferenceDict() returned nil")
	}

	// Verify it's a pages dictionary
	typeEntry := pagesDict.NameEntry("Type")
	if typeEntry == nil || *typeEntry != "Pages" {
		t.Errorf("Pages Type = %v, want 'Pages'", typeEntry)
	}

	// Should have Kids array
	if _, found := pagesDict.Find("Kids"); !found {
		t.Error("Pages dict missing 'Kids' entry")
	}

	// Should have Count
	if pagesDict.IntEntry("Count") == nil {
		t.Error("Pages dict missing 'Count' entry")
	}
}

func TestXRefTablePageNumber(t *testing.T) {
	ctx := loadTestPDF(t, "TheGoProgrammingLanguageCh1.pdf")

	// Get first page dict to get its object number
	pageDict, pageIndRef, _, err := ctx.XRefTable.PageDict(1, false)
	if err != nil {
		t.Fatalf("PageDict(1) error = %v", err)
	}

	if pageDict == nil || pageIndRef == nil {
		t.Fatal("PageDict returned nil")
	}

	pageObjNr := pageIndRef.ObjectNumber.Value()

	// Test PageNumber - should return 1 for the first page
	pageNr, err := ctx.XRefTable.PageNumber(pageObjNr)
	if err != nil {
		t.Fatalf("PageNumber(%d) error = %v", pageObjNr, err)
	}

	if pageNr != 1 {
		t.Errorf("PageNumber(%d) = %d, want 1", pageObjNr, pageNr)
	}
}
