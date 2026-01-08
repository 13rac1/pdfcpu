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

package types

import (
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/filter"
)

func TestNewStreamDict(t *testing.T) {
	d := NewDict()
	d.Insert("Type", Name("XObject"))
	d.Insert("Subtype", Name("Image"))

	streamLength := int64(100)
	streamLengthObjNr := 5

	fp := []PDFFilter{{Name: filter.Flate, DecodeParms: nil}}

	sd := NewStreamDict(d, 1234, &streamLength, &streamLengthObjNr, fp)

	if sd.StreamOffset != 1234 {
		t.Errorf("StreamOffset = %d, want 1234", sd.StreamOffset)
	}
	if sd.StreamLength == nil || *sd.StreamLength != 100 {
		t.Errorf("StreamLength = %v, want 100", sd.StreamLength)
	}
	if sd.StreamLengthObjNr == nil || *sd.StreamLengthObjNr != 5 {
		t.Errorf("StreamLengthObjNr = %v, want 5", sd.StreamLengthObjNr)
	}
	if len(sd.FilterPipeline) != 1 {
		t.Errorf("FilterPipeline length = %d, want 1", len(sd.FilterPipeline))
	}
	if sd.FilterPipeline[0].Name != filter.Flate {
		t.Errorf("FilterPipeline[0].Name = %s, want %s", sd.FilterPipeline[0].Name, filter.Flate)
	}
	if sd.IsPageContent {
		t.Error("IsPageContent should be false")
	}
	if sd.CSComponents != 0 {
		t.Errorf("CSComponents = %d, want 0", sd.CSComponents)
	}
}

func TestNewStreamDictNilValues(t *testing.T) {
	d := NewDict()
	sd := NewStreamDict(d, 0, nil, nil, nil)

	if sd.StreamLength != nil {
		t.Error("StreamLength should be nil")
	}
	if sd.StreamLengthObjNr != nil {
		t.Error("StreamLengthObjNr should be nil")
	}
	if sd.FilterPipeline != nil {
		t.Error("FilterPipeline should be nil")
	}
}

func TestStreamDictClone(t *testing.T) {
	d := NewDict()
	d.Insert("Type", Name("XObject"))
	d.Insert("Width", Integer(100))

	decodeParms := NewDict()
	decodeParms.Insert("Predictor", Integer(12))
	decodeParms.Insert("Columns", Integer(100))

	fp := []PDFFilter{
		{Name: filter.Flate, DecodeParms: decodeParms},
		{Name: filter.ASCII85, DecodeParms: nil},
	}

	streamLength := int64(500)
	sd := NewStreamDict(d, 100, &streamLength, nil, fp)
	sd.Raw = []byte{1, 2, 3}
	sd.Content = []byte{4, 5, 6}
	sd.IsPageContent = true
	sd.CSComponents = 3

	clone := sd.Clone().(StreamDict)

	// Verify clone has same values
	if clone.StreamOffset != sd.StreamOffset {
		t.Error("Clone StreamOffset mismatch")
	}
	if clone.IsPageContent != sd.IsPageContent {
		t.Error("Clone IsPageContent mismatch")
	}
	if clone.CSComponents != sd.CSComponents {
		t.Error("Clone CSComponents mismatch")
	}

	// Verify filter pipeline was cloned
	if len(clone.FilterPipeline) != 2 {
		t.Errorf("Clone FilterPipeline length = %d, want 2", len(clone.FilterPipeline))
	}
	if clone.FilterPipeline[0].Name != filter.Flate {
		t.Error("Clone FilterPipeline[0].Name mismatch")
	}

	// Verify DecodeParms was cloned
	if clone.FilterPipeline[0].DecodeParms == nil {
		t.Fatal("Clone FilterPipeline[0].DecodeParms is nil")
	}
	if pred := clone.FilterPipeline[0].DecodeParms.IntEntry("Predictor"); pred == nil || *pred != 12 {
		t.Error("Clone DecodeParms.Predictor mismatch")
	}

	// Modify original and verify clone is independent
	sd.Dict.Update("Width", Integer(200))
	if w := clone.Dict.IntEntry("Width"); w == nil || *w != 100 {
		t.Error("Clone Dict is not independent")
	}

	// Modify original DecodeParms
	sd.FilterPipeline[0].DecodeParms.Update("Predictor", Integer(99))
	if pred := clone.FilterPipeline[0].DecodeParms.IntEntry("Predictor"); pred == nil || *pred != 12 {
		t.Error("Clone DecodeParms is not independent")
	}
}

func TestStreamDictCloneNilDecodeParms(t *testing.T) {
	d := NewDict()
	fp := []PDFFilter{{Name: filter.Flate, DecodeParms: nil}}
	sd := NewStreamDict(d, 0, nil, nil, fp)

	clone := sd.Clone().(StreamDict)

	if clone.FilterPipeline[0].DecodeParms != nil {
		t.Error("Clone should have nil DecodeParms")
	}
}

func TestStreamDictHasSoleFilterNamed(t *testing.T) {
	tests := []struct {
		name       string
		filters    []PDFFilter
		filterName string
		want       bool
	}{
		{
			name:       "nil pipeline",
			filters:    nil,
			filterName: filter.Flate,
			want:       false,
		},
		{
			name:       "empty pipeline",
			filters:    []PDFFilter{},
			filterName: filter.Flate,
			want:       false,
		},
		{
			name:       "single matching filter",
			filters:    []PDFFilter{{Name: filter.Flate}},
			filterName: filter.Flate,
			want:       true,
		},
		{
			name:       "single non-matching filter",
			filters:    []PDFFilter{{Name: filter.Flate}},
			filterName: filter.LZW,
			want:       false,
		},
		{
			name:       "multiple filters",
			filters:    []PDFFilter{{Name: filter.Flate}, {Name: filter.ASCII85}},
			filterName: filter.Flate,
			want:       false,
		},
		{
			name:       "DCT filter",
			filters:    []PDFFilter{{Name: filter.DCT}},
			filterName: filter.DCT,
			want:       true,
		},
		{
			name:       "JPX filter",
			filters:    []PDFFilter{{Name: filter.JPX}},
			filterName: filter.JPX,
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sd := NewStreamDict(NewDict(), 0, nil, nil, tt.filters)
			got := sd.HasSoleFilterNamed(tt.filterName)
			if got != tt.want {
				t.Errorf("HasSoleFilterNamed(%q) = %v, want %v", tt.filterName, got, tt.want)
			}
		})
	}
}

func TestStreamDictImage(t *testing.T) {
	tests := []struct {
		name    string
		typ     string
		subtype string
		want    bool
	}{
		{"XObject/Image", "XObject", "Image", true},
		{"XObject/Form", "XObject", "Form", false},
		{"Page type", "Page", "Image", false},
		{"no type", "", "Image", false},
		{"no subtype", "XObject", "", false},
		{"both missing", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDict()
			if tt.typ != "" {
				d.Insert("Type", Name(tt.typ))
			}
			if tt.subtype != "" {
				d.Insert("Subtype", Name(tt.subtype))
			}
			sd := NewStreamDict(d, 0, nil, nil, nil)

			got := sd.Image()
			if got != tt.want {
				t.Errorf("Image() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewObjectStreamDict(t *testing.T) {
	osd := NewObjectStreamDict()

	if osd == nil {
		t.Fatal("NewObjectStreamDict() returned nil")
	}

	// Check Type
	if typ := osd.Type(); typ == nil || *typ != "ObjStm" {
		t.Errorf("Type = %v, want ObjStm", typ)
	}

	// Check Filter
	if n := osd.NameEntry("Filter"); n == nil || *n != filter.Flate {
		t.Errorf("Filter = %v, want %s", n, filter.Flate)
	}

	// Check FilterPipeline
	if len(osd.FilterPipeline) != 1 {
		t.Errorf("FilterPipeline length = %d, want 1", len(osd.FilterPipeline))
	}
	if osd.FilterPipeline[0].Name != filter.Flate {
		t.Errorf("FilterPipeline[0].Name = %s, want %s", osd.FilterPipeline[0].Name, filter.Flate)
	}
	if osd.FilterPipeline[0].DecodeParms != nil {
		t.Error("FilterPipeline[0].DecodeParms should be nil")
	}

	// Check initial state
	if osd.ObjCount != 0 {
		t.Errorf("ObjCount = %d, want 0", osd.ObjCount)
	}
	if osd.FirstObjOffset != 0 {
		t.Errorf("FirstObjOffset = %d, want 0", osd.FirstObjOffset)
	}
}

func TestObjectStreamDictIndexedObject(t *testing.T) {
	osd := NewObjectStreamDict()
	osd.ObjArray = Array{Integer(1), Integer(2), Integer(3)}

	tests := []struct {
		index   int
		want    Object
		wantErr bool
	}{
		{0, Integer(1), false},
		{1, Integer(2), false},
		{2, Integer(3), false},
		{-1, nil, true},
		{3, nil, true},
		{100, nil, true},
	}

	for _, tt := range tests {
		obj, err := osd.IndexedObject(tt.index)
		if tt.wantErr {
			if err == nil {
				t.Errorf("IndexedObject(%d) expected error", tt.index)
			}
		} else {
			if err != nil {
				t.Errorf("IndexedObject(%d) error: %v", tt.index, err)
			}
			if i, ok := obj.(Integer); !ok || i != tt.want {
				t.Errorf("IndexedObject(%d) = %v, want %v", tt.index, obj, tt.want)
			}
		}
	}
}

func TestObjectStreamDictIndexedObjectNilArray(t *testing.T) {
	osd := NewObjectStreamDict()
	osd.ObjArray = nil

	_, err := osd.IndexedObject(0)
	if err == nil {
		t.Error("IndexedObject(0) with nil array should return error")
	}
}

func TestObjectStreamDictAddObject(t *testing.T) {
	osd := NewObjectStreamDict()
	osd.Content = []byte{}

	// Add first object
	err := osd.AddObject(1, "<<>>")
	if err != nil {
		t.Fatalf("AddObject(1) error: %v", err)
	}
	if osd.ObjCount != 1 {
		t.Errorf("ObjCount = %d, want 1", osd.ObjCount)
	}
	if string(osd.Prolog) != "1 0" {
		t.Errorf("Prolog = %q, want %q", osd.Prolog, "1 0")
	}
	if string(osd.Content) != "<<>>" {
		t.Errorf("Content = %q, want %q", osd.Content, "<<>>")
	}

	// Add second object
	err = osd.AddObject(2, "[1 2 3]")
	if err != nil {
		t.Fatalf("AddObject(2) error: %v", err)
	}
	if osd.ObjCount != 2 {
		t.Errorf("ObjCount = %d, want 2", osd.ObjCount)
	}
	// Second object prolog includes offset (4 = len("<<>>"))
	if string(osd.Prolog) != "1 0 2 4" {
		t.Errorf("Prolog = %q, want %q", osd.Prolog, "1 0 2 4")
	}
	if string(osd.Content) != "<<>>[1 2 3]" {
		t.Errorf("Content = %q, want %q", osd.Content, "<<>>[1 2 3]")
	}
}

func TestObjectStreamDictFinalize(t *testing.T) {
	osd := NewObjectStreamDict()
	osd.Content = []byte{}

	osd.AddObject(1, "<<>>")
	osd.AddObject(2, "[1 2 3]")

	osd.Finalize()

	// After finalize, Content = Prolog + Content
	expectedContent := "1 0 2 4<<>>[1 2 3]"
	if string(osd.Content) != expectedContent {
		t.Errorf("Content after Finalize = %q, want %q", osd.Content, expectedContent)
	}

	// FirstObjOffset should be length of prolog
	if osd.FirstObjOffset != 7 { // len("1 0 2 4") = 7
		t.Errorf("FirstObjOffset = %d, want 7", osd.FirstObjOffset)
	}
}

func TestLazyObjectStreamObjectClone(t *testing.T) {
	osd := NewObjectStreamDict()

	l := LazyObjectStreamObject{
		osd:         osd,
		startOffset: 10,
		endOffset:   20,
	}

	clone := l.Clone().(LazyObjectStreamObject)

	if clone.osd != l.osd {
		t.Error("Clone osd pointer mismatch")
	}
	if clone.startOffset != 10 {
		t.Errorf("Clone startOffset = %d, want 10", clone.startOffset)
	}
	if clone.endOffset != 20 {
		t.Errorf("Clone endOffset = %d, want 20", clone.endOffset)
	}
}

func TestPDFFilterStruct(t *testing.T) {
	// Test PDFFilter struct
	parms := NewDict()
	parms.Insert("Predictor", Integer(12))

	f := PDFFilter{
		Name:        filter.Flate,
		DecodeParms: parms,
	}

	if f.Name != filter.Flate {
		t.Errorf("Name = %s, want %s", f.Name, filter.Flate)
	}
	if f.DecodeParms == nil {
		t.Fatal("DecodeParms is nil")
	}
	if pred := f.DecodeParms.IntEntry("Predictor"); pred == nil || *pred != 12 {
		t.Error("DecodeParms.Predictor mismatch")
	}
}

func TestStreamDictEncodeNoFilter(t *testing.T) {
	d := NewDict()
	sd := NewStreamDict(d, 0, nil, nil, nil)
	sd.Content = []byte("Hello, World!")

	err := sd.Encode()
	if err != nil {
		t.Fatalf("Encode() error: %v", err)
	}

	// Without filter, Raw should equal Content
	if string(sd.Raw) != "Hello, World!" {
		t.Errorf("Raw = %q, want %q", sd.Raw, "Hello, World!")
	}

	// StreamLength should be set
	if sd.StreamLength == nil {
		t.Fatal("StreamLength is nil")
	}
	if *sd.StreamLength != 13 {
		t.Errorf("StreamLength = %d, want 13", *sd.StreamLength)
	}

	// Dict should have Length entry
	if l := sd.IntEntry("Length"); l == nil || *l != 13 {
		t.Errorf("Dict Length = %v, want 13", l)
	}
}

func TestStreamDictEncodeAlreadyEncoded(t *testing.T) {
	d := NewDict()
	sd := NewStreamDict(d, 0, nil, nil, nil)
	sd.Raw = []byte("already encoded")
	sd.Content = nil // Content is nil, Raw is set

	err := sd.Encode()
	if err != nil {
		t.Fatalf("Encode() error: %v", err)
	}

	// Should return without changing Raw
	if string(sd.Raw) != "already encoded" {
		t.Errorf("Raw was modified unexpectedly")
	}
}

func TestStreamDictDecodeLength(t *testing.T) {
	d := NewDict()
	sd := NewStreamDict(d, 0, nil, nil, nil)
	sd.Content = []byte("Hello, World!")

	// Already decoded, should return content
	data, err := sd.DecodeLength(-1)
	if err != nil {
		t.Fatalf("DecodeLength(-1) error: %v", err)
	}
	if string(data) != "Hello, World!" {
		t.Errorf("DecodeLength(-1) = %q, want %q", data, "Hello, World!")
	}

	// With max length
	data, err = sd.DecodeLength(5)
	if err != nil {
		t.Fatalf("DecodeLength(5) error: %v", err)
	}
	if string(data) != "Hello" {
		t.Errorf("DecodeLength(5) = %q, want %q", data, "Hello")
	}
}

func TestStreamDictDecodeNoFilter(t *testing.T) {
	d := NewDict()
	sd := NewStreamDict(d, 0, nil, nil, nil)
	sd.Raw = []byte("raw data")
	sd.Content = nil

	data, err := sd.DecodeLength(-1)
	if err != nil {
		t.Fatalf("DecodeLength(-1) error: %v", err)
	}

	// Without filter, Content should equal Raw
	if string(data) != "raw data" {
		t.Errorf("DecodeLength(-1) = %q, want %q", data, "raw data")
	}
	if string(sd.Content) != "raw data" {
		t.Errorf("Content = %q, want %q", sd.Content, "raw data")
	}
}
