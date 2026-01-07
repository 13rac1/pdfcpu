/*
Copyright 2026 The pdfcpu Authors.

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

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

// TestIndexedColorPreservation tests that indexed color images maintain their color space
// when converted to PDF image objects, rather than being expanded to DeviceRGB.
// This addresses the bug where replacing an indexed color image nearly doubles the file size.
func TestIndexedColorPreservation(t *testing.T) {
	// Create indexed color image
	palette := color.Palette{
		color.RGBA{R: 255, G: 0, B: 0, A: 255},
		color.RGBA{R: 0, G: 255, B: 0, A: 255},
		color.RGBA{R: 0, G: 0, B: 255, A: 255},
		color.RGBA{R: 255, G: 255, B: 0, A: 255},
		color.RGBA{R: 255, G: 0, B: 255, A: 255},
		color.RGBA{R: 0, G: 255, B: 255, A: 255},
		color.RGBA{R: 255, G: 255, B: 255, A: 255},
		color.RGBA{R: 0, G: 0, B: 0, A: 255},
	}

	width, height := 100, 100
	img := image.NewPaletted(image.Rect(0, 0, width, height), palette)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.SetColorIndex(x, y, uint8((x/10)%len(palette)))
		}
	}

	// Encode to PNG
	var pngBuf bytes.Buffer
	if err := png.Encode(&pngBuf, img); err != nil {
		t.Fatalf("Failed to encode PNG: %v", err)
	}

	// Create minimal XRefTable
	size := 1
	version := V17
	rootDict := types.NewDict()
	rootDict.InsertName("Type", "Catalog")

	xRefTable := &XRefTable{
		Size:          &size,
		HeaderVersion: &version,
		Table:         map[int]*XRefTableEntry{0: NewFreeHeadXRefTableEntry()},
	}
	ir, err := xRefTable.IndRefForNewObject(rootDict)
	if err != nil {
		t.Fatalf("Failed to create root indirect reference: %v", err)
	}
	xRefTable.Root = ir

	// Create image stream dict from indexed color PNG
	sd, _, _, err := CreateImageStreamDict(xRefTable, &pngBuf)
	if err != nil {
		t.Fatalf("Failed to create image stream dict: %v", err)
	}

	// Verify ColorSpace is Indexed array, not DeviceRGB
	csObj, found := sd.Find("ColorSpace")
	if !found {
		t.Fatal("ColorSpace not found in stream dict")
	}

	csArray, ok := csObj.(types.Array)
	if !ok {
		t.Fatalf("ColorSpace is %T, expected Indexed array", csObj)
	}
	if len(csArray) != 4 {
		t.Fatalf("Indexed ColorSpace has %d elements, want 4", len(csArray))
	}

	csName, ok := csArray[0].(types.Name)
	if !ok {
		t.Fatalf("ColorSpace[0] is %T, want types.Name", csArray[0])
	}
	if csName != IndexedCS {
		t.Errorf("ColorSpace is %s, want %s", csName, IndexedCS)
	}

	// Verify BitsPerComponent is 8 (palette indices)
	bpcObj, found := sd.Find("BitsPerComponent")
	if !found {
		t.Fatal("BitsPerComponent not found")
	}
	bpc, ok := bpcObj.(types.Integer)
	if !ok {
		t.Fatalf("BitsPerComponent is %T, want types.Integer", bpcObj)
	}
	if bpc != 8 {
		t.Errorf("BitsPerComponent is %d, want 8", bpc)
	}
}
