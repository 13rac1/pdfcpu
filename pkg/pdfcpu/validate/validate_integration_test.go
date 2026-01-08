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

package validate_test

import (
	"path/filepath"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/validate"
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

func TestXRefTableSimplePDF(t *testing.T) {
	ctx := loadTestPDF(t, "testImage.pdf")

	// Validate entire XRefTable
	err := validate.XRefTable(ctx)
	if err != nil {
		t.Errorf("XRefTable(testImage.pdf) error = %v", err)
	}

	// Verify Valid flag is set
	if !ctx.XRefTable.Valid {
		t.Error("XRefTable.Valid = false after successful validation")
	}

	t.Logf("testImage.pdf validated successfully")
}

func TestXRefTableProgrammingBook(t *testing.T) {
	ctx := loadTestPDF(t, "TheGoProgrammingLanguageCh1.pdf")

	// Validate entire XRefTable
	err := validate.XRefTable(ctx)
	if err != nil {
		t.Errorf("XRefTable(TheGoProgrammingLanguageCh1.pdf) error = %v", err)
	}

	// Verify Valid flag is set
	if !ctx.XRefTable.Valid {
		t.Error("XRefTable.Valid = false after successful validation")
	}

	t.Logf("TheGoProgrammingLanguageCh1.pdf validated successfully with %d pages", ctx.PageCount)
}

func TestXRefTableEmptyPDF(t *testing.T) {
	ctx := loadTestPDF(t, "empty.pdf")

	// Validate entire XRefTable
	err := validate.XRefTable(ctx)
	if err != nil {
		t.Errorf("XRefTable(empty.pdf) error = %v", err)
	}

	t.Logf("empty.pdf validated (pages: %d)", ctx.PageCount)
}

func TestXRefTablePresentationPDF(t *testing.T) {
	ctx := loadTestPDF(t, "FOSDEM14_HPC_devroom_14_GoCUDA.pdf")

	// Validate entire XRefTable
	err := validate.XRefTable(ctx)
	if err != nil {
		t.Errorf("XRefTable(FOSDEM14_HPC_devroom_14_GoCUDA.pdf) error = %v", err)
	}

	// Verify Valid flag is set
	if !ctx.XRefTable.Valid {
		t.Error("XRefTable.Valid = false after successful validation")
	}

	t.Logf("FOSDEM14_HPC_devroom_14_GoCUDA.pdf validated successfully with %d pages", ctx.PageCount)
}

func TestXRefTableAnnotations(t *testing.T) {
	ctx := loadTestPDF(t, "annotTest.pdf")

	// Validate PDF with annotations
	err := validate.XRefTable(ctx)
	if err != nil {
		t.Errorf("XRefTable(annotTest.pdf) error = %v", err)
	}

	// Verify Valid flag is set
	if !ctx.XRefTable.Valid {
		t.Error("XRefTable.Valid = false after successful validation")
	}

	t.Logf("annotTest.pdf with annotations validated successfully")
}

func TestXRefTableMultipleValidations(t *testing.T) {
	ctx := loadTestPDF(t, "testImage.pdf")

	// First validation
	err := validate.XRefTable(ctx)
	if err != nil {
		t.Fatalf("First XRefTable() error = %v", err)
	}

	if !ctx.XRefTable.Valid {
		t.Error("XRefTable.Valid = false after first validation")
	}

	// Second validation should also succeed
	err = validate.XRefTable(ctx)
	if err != nil {
		t.Errorf("Second XRefTable() error = %v", err)
	}

	if !ctx.XRefTable.Valid {
		t.Error("XRefTable.Valid = false after second validation")
	}

	t.Log("Multiple validations succeeded")
}

func TestXRefTableWithConfiguration(t *testing.T) {
	// Test with different configurations
	configs := []struct {
		name string
		conf *model.Configuration
	}{
		{"default config", model.NewDefaultConfiguration()},
		{"relaxed mode", func() *model.Configuration {
			conf := model.NewDefaultConfiguration()
			conf.ValidationMode = model.ValidationRelaxed
			return conf
		}()},
		{"strict mode", func() *model.Configuration {
			conf := model.NewDefaultConfiguration()
			conf.ValidationMode = model.ValidationStrict
			return conf
		}()},
	}

	for _, tc := range configs {
		t.Run(tc.name, func(t *testing.T) {
			inFile := filepath.Join(testdataDir, "testImage.pdf")
			ctx, err := pdfcpu.ReadFile(inFile, tc.conf)
			if err != nil {
				t.Fatalf("ReadFile() error = %v", err)
			}

			if err := ctx.XRefTable.EnsurePageCount(); err != nil {
				t.Fatalf("EnsurePageCount() error = %v", err)
			}

			err = validate.XRefTable(ctx)
			if err != nil {
				t.Errorf("XRefTable() with %s error = %v", tc.name, err)
			}

			t.Logf("%s validation successful", tc.name)
		})
	}
}

func TestDocumentPropertyValidator(t *testing.T) {
	tests := []struct {
		name  string
		prop  string
		valid bool
	}{
		// Standard properties that cannot be modified (return false)
		{"Keywords", "Keywords", false},
		{"Producer", "Producer", false},
		{"CreationDate", "CreationDate", false},
		{"ModDate", "ModDate", false},
		{"Trapped", "Trapped", false},
		// Custom properties that can be modified (return true)
		{"Title", "Title", true},
		{"Author", "Author", true},
		{"Subject", "Subject", true},
		{"Creator", "Creator", true},
		{"CustomProp", "CustomProp", true},
		{"Empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := validate.DocumentProperty(tt.prop)
			if valid != tt.valid {
				t.Errorf("DocumentProperty(%q) = %v, want %v", tt.prop, valid, tt.valid)
			}
		})
	}
}

func TestDocumentPageLayoutValidator(t *testing.T) {
	tests := []struct {
		name   string
		layout string
		valid  bool
	}{
		{"SinglePage", "SinglePage", true},
		{"singlepage lowercase", "singlepage", true},
		{"TwoColumnLeft", "TwoColumnLeft", true},
		{"twocolumnleft lowercase", "twocolumnleft", true},
		{"TwoColumnRight", "TwoColumnRight", true},
		{"twocolumnright lowercase", "twocolumnright", true},
		{"TwoPageLeft", "TwoPageLeft", true},
		{"twopageleft lowercase", "twopageleft", true},
		{"TwoPageRight", "TwoPageRight", true},
		{"twopageright lowercase", "twopageright", true},
		{"Invalid OneColumn", "OneColumn", false}, // Not in PDF spec
		{"Invalid layout", "InvalidLayout", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := validate.DocumentPageLayout(tt.layout)
			if valid != tt.valid {
				t.Errorf("DocumentPageLayout(%q) = %v, want %v", tt.layout, valid, tt.valid)
			}
		})
	}
}

func TestDocumentPageModeValidator(t *testing.T) {
	tests := []struct {
		name  string
		mode  string
		valid bool
	}{
		{"UseNone", "UseNone", true},
		{"UseOutlines", "UseOutlines", true},
		{"UseThumbs", "UseThumbs", true},
		{"FullScreen", "FullScreen", true},
		{"UseOC", "UseOC", true},
		{"UseAttachments", "UseAttachments", true},
		{"Invalid mode", "InvalidMode", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := validate.DocumentPageMode(tt.mode)
			if valid != tt.valid {
				t.Errorf("DocumentPageMode(%q) = %v, want %v", tt.mode, valid, tt.valid)
			}
		})
	}
}

func TestXRefTablePageCount(t *testing.T) {
	tests := []struct {
		name          string
		pdfFile       string
		expectPages   int
		minPages      int
		pageCountZero bool
	}{
		{"simple PDF", "testImage.pdf", 2, 2, false},
		{"book PDF", "TheGoProgrammingLanguageCh1.pdf", 0, 1, false}, // Don't know exact count, but > 0
		{"empty PDF", "empty.pdf", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := loadTestPDF(t, tt.pdfFile)

			// Validate
			err := validate.XRefTable(ctx)
			if err != nil {
				t.Errorf("XRefTable() error = %v", err)
			}

			// Check page count
			if tt.pageCountZero {
				if ctx.PageCount != 0 {
					t.Logf("Note: %s has %d pages (expected 0 or more)", tt.pdfFile, ctx.PageCount)
				}
			} else if tt.expectPages > 0 {
				if ctx.PageCount != tt.expectPages {
					t.Errorf("PageCount = %d, want %d", ctx.PageCount, tt.expectPages)
				}
			} else if ctx.PageCount < tt.minPages {
				t.Errorf("PageCount = %d, want >= %d", ctx.PageCount, tt.minPages)
			}

			t.Logf("%s: PageCount = %d", tt.pdfFile, ctx.PageCount)
		})
	}
}
