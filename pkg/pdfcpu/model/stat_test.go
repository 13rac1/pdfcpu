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

package model

import "testing"

func TestNewPDFStats(t *testing.T) {
	stats := NewPDFStats()
	if stats.rootAttrs == nil {
		t.Error("rootAttrs should not be nil")
	}
	if stats.pageAttrs == nil {
		t.Error("pageAttrs should not be nil")
	}
}

func TestPDFStatsRootAttrs(t *testing.T) {
	stats := NewPDFStats()

	// Initially should not have any attrs
	if stats.UsesRootAttr(RootVersion) {
		t.Error("RootVersion should not be set initially")
	}

	// Add an attr
	stats.AddRootAttr(RootVersion)
	if !stats.UsesRootAttr(RootVersion) {
		t.Error("RootVersion should be set after AddRootAttr")
	}

	// Other attrs should still be false
	if stats.UsesRootAttr(RootExtensions) {
		t.Error("RootExtensions should not be set")
	}

	// Add multiple attrs
	stats.AddRootAttr(RootExtensions)
	stats.AddRootAttr(RootNames)
	if !stats.UsesRootAttr(RootExtensions) {
		t.Error("RootExtensions should be set")
	}
	if !stats.UsesRootAttr(RootNames) {
		t.Error("RootNames should be set")
	}
}

func TestPDFStatsPageAttrs(t *testing.T) {
	stats := NewPDFStats()

	// Initially should not have any attrs
	if stats.UsesPageAttr(PageMediaBox) {
		t.Error("PageMediaBox should not be set initially")
	}

	// Add an attr
	stats.AddPageAttr(PageMediaBox)
	if !stats.UsesPageAttr(PageMediaBox) {
		t.Error("PageMediaBox should be set after AddPageAttr")
	}

	// Other attrs should still be false
	if stats.UsesPageAttr(PageCropBox) {
		t.Error("PageCropBox should not be set")
	}

	// Add multiple attrs
	stats.AddPageAttr(PageCropBox)
	stats.AddPageAttr(PageContents)
	if !stats.UsesPageAttr(PageCropBox) {
		t.Error("PageCropBox should be set")
	}
	if !stats.UsesPageAttr(PageContents) {
		t.Error("PageContents should be set")
	}
}

func TestRootConstants(t *testing.T) {
	// Verify a few key constants
	if RootVersion != 0 {
		t.Errorf("RootVersion = %d, want 0", RootVersion)
	}
	if RootAcroForm != 13 {
		t.Errorf("RootAcroForm = %d, want 13", RootAcroForm)
	}
}

func TestPageConstants(t *testing.T) {
	// Verify a few key constants
	if PageLastModified != 0 {
		t.Errorf("PageLastModified = %d, want 0", PageLastModified)
	}
	if PageMediaBox != 2 {
		t.Errorf("PageMediaBox = %d, want 2", PageMediaBox)
	}
	if PageContents != 8 {
		t.Errorf("PageContents = %d, want 8", PageContents)
	}
}
