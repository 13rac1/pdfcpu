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

func TestXRefTableExists(t *testing.T) {
	ctx := loadTestPDF(t, "testImage.pdf")

	tests := []struct {
		name   string
		objNr  int
		exists bool
	}{
		{"object 0 (free list head)", 0, true}, // Object 0 always exists as free list head
		{"object 1", 1, true},
		{"object 2", 2, true},
		{"object 3", 3, true},
		{"object 9999", 9999, false},
		{"negative object", -1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists := ctx.XRefTable.Exists(tt.objNr)
			if exists != tt.exists {
				t.Errorf("Exists(%d) = %v, want %v", tt.objNr, exists, tt.exists)
			}
		})
	}
}

func TestXRefTableFind(t *testing.T) {
	ctx := loadTestPDF(t, "testImage.pdf")

	tests := []struct {
		name      string
		objNr     int
		wantFound bool
		wantFree  bool
	}{
		{"free list head object 0", 0, true, true},  // Object 0 is free
		{"existing object 1", 1, true, false},
		{"existing object 2", 2, true, false},
		{"existing object 3", 3, true, false},
		{"non-existing object 9999", 9999, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry, found := ctx.XRefTable.Find(tt.objNr)

			if found != tt.wantFound {
				t.Errorf("Find(%d) found = %v, want %v", tt.objNr, found, tt.wantFound)
			}

			if found {
				if entry == nil {
					t.Errorf("Find(%d) returned nil entry when found=true", tt.objNr)
				}
				if entry.Free != tt.wantFree {
					t.Errorf("Find(%d).Free = %v, want %v", tt.objNr, entry.Free, tt.wantFree)
				}
			}
		})
	}
}

func TestXRefTableFindObject(t *testing.T) {
	ctx := loadTestPDF(t, "testImage.pdf")

	// Get catalog object number
	catalog, err := ctx.XRefTable.Catalog()
	if err != nil {
		t.Fatalf("Catalog() error = %v", err)
	}

	// Get Pages reference
	pagesObj, found := catalog.Find("Pages")
	if !found {
		t.Fatal("Catalog missing Pages entry")
	}
	pagesRef := pagesObj.(types.IndirectRef)

	tests := []struct {
		name      string
		objNr     int
		wantErr   bool
		checkType string
	}{
		{"pages object", pagesRef.ObjectNumber.Value(), false, "dict"},
		{"catalog object", 1, false, "dict"},
		{"non-existing object", 9999, true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj, err := ctx.XRefTable.FindObject(tt.objNr)

			if tt.wantErr {
				if err == nil {
					t.Errorf("FindObject(%d) expected error, got nil", tt.objNr)
				}
				return
			}

			if err != nil {
				t.Fatalf("FindObject(%d) error = %v", tt.objNr, err)
			}

			if obj == nil {
				t.Errorf("FindObject(%d) returned nil object", tt.objNr)
			}

			// Verify object type
			switch tt.checkType {
			case "dict":
				if _, ok := obj.(types.Dict); !ok {
					t.Errorf("FindObject(%d) = %T, want types.Dict", tt.objNr, obj)
				}
			}
		})
	}
}

func TestXRefTableFindTableEntry(t *testing.T) {
	ctx := loadTestPDF(t, "testImage.pdf")

	tests := []struct {
		name      string
		objNr     int
		genNr     int
		wantFound bool
	}{
		{"object 0 gen 65535 (free)", 0, 65535, true}, // Object 0 with special generation
		{"object 1 gen 0", 1, 0, true},
		{"object 2 gen 0", 2, 0, true},
		{"object 3 gen 0", 3, 0, true},
		{"non-existing object 9999", 9999, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry, found := ctx.XRefTable.FindTableEntry(tt.objNr, tt.genNr)

			if found != tt.wantFound {
				t.Errorf("FindTableEntry(%d, %d) found = %v, want %v", tt.objNr, tt.genNr, found, tt.wantFound)
			}

			if found && entry == nil {
				t.Errorf("FindTableEntry(%d, %d) returned nil entry when found=true", tt.objNr, tt.genNr)
			}
		})
	}
}

func TestXRefTableFindTableEntryForIndRef(t *testing.T) {
	ctx := loadTestPDF(t, "testImage.pdf")

	// Get catalog to find an indirect reference
	catalog, err := ctx.XRefTable.Catalog()
	if err != nil {
		t.Fatalf("Catalog() error = %v", err)
	}

	pagesObj, found := catalog.Find("Pages")
	if !found {
		t.Fatal("Catalog missing Pages entry")
	}

	pagesRef := pagesObj.(types.IndirectRef)

	// Test with valid indirect reference
	entry, found := ctx.XRefTable.FindTableEntryForIndRef(&pagesRef)
	if !found {
		t.Error("FindTableEntryForIndRef(Pages) found = false, want true")
	}
	if entry == nil {
		t.Error("FindTableEntryForIndRef(Pages) returned nil entry")
	}

	// Test with invalid indirect reference
	invalidRef := types.IndirectRef{
		ObjectNumber:     types.Integer(9999),
		GenerationNumber: types.Integer(0),
	}
	_, found = ctx.XRefTable.FindTableEntryForIndRef(&invalidRef)
	if found {
		t.Error("FindTableEntryForIndRef(invalid) found = true, want false")
	}
}

func TestXRefTableCatalogHasPieceInfo(t *testing.T) {
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

			hasPieceInfo, err := ctx.XRefTable.CatalogHasPieceInfo()
			if err != nil {
				t.Fatalf("CatalogHasPieceInfo() error = %v", err)
			}

			// We don't assert the value, just verify no error
			// (most simple PDFs won't have PieceInfo)
			t.Logf("PDF %s has PieceInfo: %v", tt.pdfFile, hasPieceInfo)
		})
	}
}

func TestXRefTableNamesDict(t *testing.T) {
	ctx := loadTestPDF(t, "testImage.pdf")

	namesDict, err := ctx.XRefTable.NamesDict()
	if err != nil {
		t.Fatalf("NamesDict() error = %v", err)
	}

	// namesDict may be nil if PDF has no Names dictionary
	// Just verify no error occurred
	t.Logf("NamesDict exists: %v", namesDict != nil)
}

func TestXRefTableIDFirstElement(t *testing.T) {
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

			id, err := ctx.XRefTable.IDFirstElement()
			if err != nil {
				// Some PDFs may not have ID array - that's OK
				t.Logf("IDFirstElement() error (expected for some PDFs): %v", err)
				return
			}

			if id != nil && len(id) == 0 {
				t.Error("IDFirstElement() returned empty byte slice")
			}

			t.Logf("PDF %s ID first element length: %d", tt.pdfFile, len(id))
		})
	}
}

func TestXRefTablePDF20(t *testing.T) {
	tests := []struct {
		name    string
		pdfFile string
	}{
		{"simple PDF", "testImage.pdf"},
		{"book PDF", "TheGoProgrammingLanguageCh1.pdf"},
		{"empty PDF", "empty.pdf"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := loadTestPDF(t, tt.pdfFile)

			isPDF20 := ctx.XRefTable.PDF20()

			// Most test PDFs are 1.x, not 2.0
			t.Logf("PDF %s is PDF 2.0: %v", tt.pdfFile, isPDF20)
		})
	}
}

func TestXRefTableMissingObjects(t *testing.T) {
	ctx := loadTestPDF(t, "testImage.pdf")

	count, details := ctx.XRefTable.MissingObjects()

	// Well-formed PDFs should have no missing objects
	if count > 0 {
		t.Logf("Warning: PDF has %d missing objects: %s", count, *details)
	}
}

func TestXRefTableDereferenceInteger(t *testing.T) {
	ctx := loadTestPDF(t, "testImage.pdf")

	// Get catalog to find integer values
	catalog, err := ctx.XRefTable.Catalog()
	if err != nil {
		t.Fatalf("Catalog() error = %v", err)
	}

	// Get Pages dict
	pagesObj, _ := catalog.Find("Pages")
	pagesDict, _ := ctx.XRefTable.DereferenceDict(pagesObj)

	// Count should be an integer
	countObj, found := pagesDict.Find("Count")
	if !found {
		t.Fatal("Pages dict missing Count entry")
	}

	count, err := ctx.XRefTable.DereferenceInteger(countObj)
	if err != nil {
		t.Fatalf("DereferenceInteger(Count) error = %v", err)
	}

	if count == nil {
		t.Error("DereferenceInteger(Count) returned nil")
	}

	if *count <= 0 {
		t.Errorf("DereferenceInteger(Count) = %d, want > 0", *count)
	}

	t.Logf("Page count: %d", *count)
}

func TestXRefTableDereferenceNumber(t *testing.T) {
	ctx := loadTestPDF(t, "testImage.pdf")

	// Get a page dict to find number values
	pageDict, _, _, err := ctx.XRefTable.PageDict(1, false)
	if err != nil {
		t.Fatalf("PageDict(1) error = %v", err)
	}

	// MediaBox contains numbers
	mediaBoxObj, found := pageDict.Find("MediaBox")
	if !found {
		t.Fatal("Page dict missing MediaBox")
	}

	mediaBoxArray, err := ctx.XRefTable.DereferenceArray(mediaBoxObj)
	if err != nil {
		t.Fatalf("DereferenceArray(MediaBox) error = %v", err)
	}

	if len(mediaBoxArray) < 4 {
		t.Fatalf("MediaBox array length = %d, want >= 4", len(mediaBoxArray))
	}

	// Test dereferencing a number from the array
	width, err := ctx.XRefTable.DereferenceNumber(mediaBoxArray[2])
	if err != nil {
		t.Fatalf("DereferenceNumber(width) error = %v", err)
	}

	if width <= 0 {
		t.Errorf("DereferenceNumber(width) = %f, want > 0", width)
	}

	t.Logf("MediaBox width: %f", width)
}

func TestXRefTableDereferenceName(t *testing.T) {
	ctx := loadTestPDF(t, "testImage.pdf")

	// Get catalog
	catalog, err := ctx.XRefTable.Catalog()
	if err != nil {
		t.Fatalf("Catalog() error = %v", err)
	}

	// Type should be a name
	typeObj, found := catalog.Find("Type")
	if !found {
		t.Fatal("Catalog missing Type entry")
	}

	typeName, err := ctx.XRefTable.DereferenceName(typeObj, model.V10, nil)
	if err != nil {
		t.Fatalf("DereferenceName(Type) error = %v", err)
	}

	if typeName != "Catalog" {
		t.Errorf("DereferenceName(Type) = %q, want 'Catalog'", typeName)
	}
}

func TestXRefTableDereferenceArray(t *testing.T) {
	ctx := loadTestPDF(t, "testImage.pdf")

	// Get a page dict
	pageDict, _, _, err := ctx.XRefTable.PageDict(1, false)
	if err != nil {
		t.Fatalf("PageDict(1) error = %v", err)
	}

	// MediaBox should be an array
	mediaBoxObj, found := pageDict.Find("MediaBox")
	if !found {
		t.Fatal("Page dict missing MediaBox")
	}

	mediaBox, err := ctx.XRefTable.DereferenceArray(mediaBoxObj)
	if err != nil {
		t.Fatalf("DereferenceArray(MediaBox) error = %v", err)
	}

	if len(mediaBox) != 4 {
		t.Errorf("DereferenceArray(MediaBox) length = %d, want 4", len(mediaBox))
	}
}

func TestXRefTableDereferenceText(t *testing.T) {
	ctx := loadTestPDF(t, "testImage.pdf")

	// DereferenceText works with string literals and hex strings
	// For this test, we just verify the method exists and doesn't panic on valid types
	// Creating a simple test with a name converted to text
	testName := types.Name("TestName")

	// Convert name to text - note: names aren't valid for DereferenceText
	// but we can test with a string literal
	testString := types.StringLiteral("Test String")
	text, err := ctx.XRefTable.DereferenceText(testString)
	if err != nil {
		t.Fatalf("DereferenceText(StringLiteral) error = %v", err)
	}

	if text != "Test String" {
		t.Errorf("DereferenceText(StringLiteral) = %q, want 'Test String'", text)
	}

	// Test with indirect reference to a string (if we can find one)
	t.Logf("DereferenceText works with string literals: testName=%v", testName)
}

func TestXRefTablePageDictIndRef(t *testing.T) {
	ctx := loadTestPDF(t, "testImage.pdf")

	// Get indirect reference for first page
	indRef, err := ctx.XRefTable.PageDictIndRef(1)
	if err != nil {
		t.Fatalf("PageDictIndRef(1) error = %v", err)
	}

	if indRef == nil {
		t.Fatal("PageDictIndRef(1) returned nil")
	}

	if indRef.ObjectNumber.Value() <= 0 {
		t.Errorf("PageDictIndRef(1) object number = %d, want > 0", indRef.ObjectNumber.Value())
	}

	t.Logf("Page 1 indirect reference: %d %d R", indRef.ObjectNumber.Value(), indRef.GenerationNumber.Value())

	// Test page number beyond page count
	invalidPageNr := ctx.PageCount + 100
	indRef2, err := ctx.XRefTable.PageDictIndRef(invalidPageNr)
	if err == nil && indRef2 != nil {
		t.Errorf("PageDictIndRef(%d) should fail for page beyond count", invalidPageNr)
	}
}

func TestXRefTableRectForArray(t *testing.T) {
	ctx := loadTestPDF(t, "testImage.pdf")

	// Get a page dict
	pageDict, _, _, err := ctx.XRefTable.PageDict(1, false)
	if err != nil {
		t.Fatalf("PageDict(1) error = %v", err)
	}

	// Get MediaBox array
	mediaBoxObj, found := pageDict.Find("MediaBox")
	if !found {
		t.Fatal("Page dict missing MediaBox")
	}

	mediaBoxArray, err := ctx.XRefTable.DereferenceArray(mediaBoxObj)
	if err != nil {
		t.Fatalf("DereferenceArray(MediaBox) error = %v", err)
	}

	// Convert to Rectangle
	rect, err := ctx.XRefTable.RectForArray(mediaBoxArray)
	if err != nil {
		t.Fatalf("RectForArray(MediaBox) error = %v", err)
	}

	if rect == nil {
		t.Fatal("RectForArray(MediaBox) returned nil")
	}

	// Verify rectangle has positive dimensions
	if rect.Width() <= 0 || rect.Height() <= 0 {
		t.Errorf("RectForArray(MediaBox) = {width: %f, height: %f}, want positive values", rect.Width(), rect.Height())
	}

	t.Logf("MediaBox rectangle: %v", rect)
}

func TestContextString(t *testing.T) {
	ctx := loadTestPDF(t, "testImage.pdf")

	str := ctx.String()
	if str == "" {
		t.Error("Context.String() returned empty string")
	}

	t.Logf("Context string: %s", str)
}

func TestContextUnitString(t *testing.T) {
	ctx := loadTestPDF(t, "testImage.pdf")

	unitStr := ctx.UnitString()
	if unitStr == "" {
		t.Error("Context.UnitString() returned empty string")
	}

	t.Logf("Context unit: %s", unitStr)
}
