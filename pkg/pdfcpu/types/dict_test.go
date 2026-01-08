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
	"strings"
	"testing"
)

func TestEncodeDict(t *testing.T) {
	dict := Dict{
		"A()": Integer(1),
	}
	expected := `<</A#28#29 1>>`
	s := dict.PDFString()
	if s != expected {
		t.Errorf("expected %s for %+v, got %s", expected, dict, s)
	}
}

func TestNewDict(t *testing.T) {
	d := NewDict()
	if d == nil {
		t.Fatal("NewDict() returned nil")
	}
	if d.Len() != 0 {
		t.Errorf("NewDict().Len() = %d, want 0", d.Len())
	}
}

func TestDictLen(t *testing.T) {
	d := NewDict()
	if d.Len() != 0 {
		t.Errorf("empty dict Len() = %d, want 0", d.Len())
	}

	d.Insert("Key1", Integer(1))
	if d.Len() != 1 {
		t.Errorf("dict with 1 entry Len() = %d, want 1", d.Len())
	}

	d.Insert("Key2", Integer(2))
	if d.Len() != 2 {
		t.Errorf("dict with 2 entries Len() = %d, want 2", d.Len())
	}
}

func TestDictInsert(t *testing.T) {
	d := NewDict()

	// Insert new key should return true
	if !d.Insert("Key", Integer(1)) {
		t.Error("Insert() returned false for new key")
	}

	// Insert same key again should return false
	if d.Insert("Key", Integer(2)) {
		t.Error("Insert() returned true for existing key")
	}

	// Value should be unchanged
	obj, found := d.Find("Key")
	if !found {
		t.Fatal("Find() returned false for inserted key")
	}
	if i, ok := obj.(Integer); !ok || int(i) != 1 {
		t.Errorf("Find() = %v, want Integer(1)", obj)
	}
}

func TestDictInsertTypes(t *testing.T) {
	d := NewDict()

	d.InsertBool("bool", true)
	if b := d.BooleanEntry("bool"); b == nil || !*b {
		t.Error("InsertBool() did not store true")
	}

	d.InsertInt("int", 42)
	if i := d.IntEntry("int"); i == nil || *i != 42 {
		t.Errorf("InsertInt() stored %v, want 42", i)
	}

	d.InsertFloat("float", 3.14)
	obj, _ := d.Find("float")
	if f, ok := obj.(Float); !ok || float32(f) != 3.14 {
		t.Errorf("InsertFloat() stored %v, want 3.14", obj)
	}

	d.InsertString("string", "hello")
	if s := d.StringEntry("string"); s == nil || *s != "hello" {
		t.Errorf("InsertString() stored %v, want 'hello'", s)
	}

	d.InsertName("name", "Type")
	if n := d.NameEntry("name"); n == nil || *n != "Type" {
		t.Errorf("InsertName() stored %v, want 'Type'", n)
	}
}

func TestDictUpdate(t *testing.T) {
	d := NewDict()
	d.Insert("Key", Integer(1))

	d.Update("Key", Integer(2))
	obj, _ := d.Find("Key")
	if i, ok := obj.(Integer); !ok || int(i) != 2 {
		t.Errorf("Update() did not change value, got %v", obj)
	}

	// Update with nil should be a no-op
	d.Update("Key", nil)
	obj, _ = d.Find("Key")
	if i, ok := obj.(Integer); !ok || int(i) != 2 {
		t.Errorf("Update(nil) changed value to %v", obj)
	}
}

func TestDictFind(t *testing.T) {
	d := NewDict()
	d.Insert("Key", Integer(1))

	// Find existing key
	obj, found := d.Find("Key")
	if !found {
		t.Error("Find() returned false for existing key")
	}
	if obj == nil {
		t.Error("Find() returned nil object for existing key")
	}

	// Find non-existing key
	_, found = d.Find("NonExistent")
	if found {
		t.Error("Find() returned true for non-existing key")
	}

	// Find encoded name
	d.Insert("A#20B", Integer(2))
	obj, found = d.Find("A B")
	if !found {
		t.Error("Find() should decode name and find 'A B'")
	}
}

func TestDictDelete(t *testing.T) {
	d := NewDict()
	d.Insert("Key", Integer(1))

	// Delete existing key
	obj := d.Delete("Key")
	if obj == nil {
		t.Error("Delete() returned nil for existing key")
	}
	if _, found := d.Find("Key"); found {
		t.Error("Delete() did not remove key")
	}

	// Delete non-existing key
	obj = d.Delete("NonExistent")
	if obj != nil {
		t.Errorf("Delete() returned %v for non-existing key", obj)
	}
}

func TestDictNewIDForPrefix(t *testing.T) {
	d := NewDict()
	d.Insert("Img0", Integer(1))
	d.Insert("Img1", Integer(2))

	id := d.NewIDForPrefix("Img", 0)
	if id != "Img2" {
		t.Errorf("NewIDForPrefix() = %q, want 'Img2'", id)
	}

	id = d.NewIDForPrefix("Font", 0)
	if id != "Font0" {
		t.Errorf("NewIDForPrefix() = %q, want 'Font0'", id)
	}
}

func TestDictEntry(t *testing.T) {
	d := NewDict()
	d.Insert("Key", Integer(1))

	// Find existing entry
	obj, found, err := d.Entry("TestDict", "Key", false)
	if err != nil {
		t.Fatalf("Entry() error: %v", err)
	}
	if !found {
		t.Error("Entry() did not find existing key")
	}
	if obj == nil {
		t.Error("Entry() returned nil object")
	}

	// Required entry missing
	_, _, err = d.Entry("TestDict", "Missing", true)
	if err == nil {
		t.Error("Entry() should return error for missing required entry")
	}

	// Optional entry missing
	obj, found, err = d.Entry("TestDict", "Missing", false)
	if err != nil {
		t.Fatalf("Entry() error: %v", err)
	}
	if found {
		t.Error("Entry() should return found=false for missing optional entry")
	}

	// Required entry with nil value
	d.Insert("NilKey", nil)
	_, _, err = d.Entry("TestDict", "NilKey", true)
	if err == nil {
		t.Error("Entry() should return error for nil required entry")
	}
}

func TestDictTypeEntries(t *testing.T) {
	d := NewDict()
	d.Insert("Boolean", Boolean(true))
	d.Insert("Integer", Integer(42))
	d.Insert("String", StringLiteral("hello"))
	d.Insert("Name", Name("Type"))
	d.Insert("Array", Array{Integer(1), Integer(2)})
	d.Insert("Dict", Dict{"Key": Integer(1)})
	d.Insert("IndRef", IndirectRef{ObjectNumber: Integer(1), GenerationNumber: Integer(0)})
	d.Insert("StringLit", StringLiteral("test"))
	d.Insert("HexLit", HexLiteral("48656C6C6F"))

	// Test each entry type
	if b := d.BooleanEntry("Boolean"); b == nil || !*b {
		t.Error("BooleanEntry() failed")
	}
	if d.BooleanEntry("NonExistent") != nil {
		t.Error("BooleanEntry() should return nil for missing key")
	}
	if d.BooleanEntry("Integer") != nil {
		t.Error("BooleanEntry() should return nil for wrong type")
	}

	if i := d.IntEntry("Integer"); i == nil || *i != 42 {
		t.Error("IntEntry() failed")
	}
	if i := d.Int64Entry("Integer"); i == nil || *i != 42 {
		t.Error("Int64Entry() failed")
	}

	if s := d.StringEntry("String"); s == nil || *s != "hello" {
		t.Error("StringEntry() failed")
	}

	if n := d.NameEntry("Name"); n == nil || *n != "Type" {
		t.Error("NameEntry() failed")
	}

	if a := d.ArrayEntry("Array"); a == nil || len(a) != 2 {
		t.Error("ArrayEntry() failed")
	}
	if d.ArrayEntry("NonExistent") != nil {
		t.Error("ArrayEntry() should return nil for missing key")
	}

	if sub := d.DictEntry("Dict"); sub == nil {
		t.Error("DictEntry() failed")
	}

	if ref := d.IndirectRefEntry("IndRef"); ref == nil {
		t.Error("IndirectRefEntry() failed")
	}

	if sl := d.StringLiteralEntry("StringLit"); sl == nil {
		t.Error("StringLiteralEntry() failed")
	}

	if hl := d.HexLiteralEntry("HexLit"); hl == nil {
		t.Error("HexLiteralEntry() failed")
	}
}

func TestDictStringOrHexLiteralEntry(t *testing.T) {
	d := NewDict()
	d.Insert("String", StringLiteral("hello"))
	d.Insert("Hex", HexLiteral("48656C6C6F"))

	s, err := d.StringOrHexLiteralEntry("String")
	if err != nil || s == nil || *s != "hello" {
		t.Error("StringOrHexLiteralEntry() failed for StringLiteral")
	}

	s, err = d.StringOrHexLiteralEntry("Hex")
	if err != nil || s == nil || *s != "Hello" {
		t.Errorf("StringOrHexLiteralEntry() failed for HexLiteral, got %v", s)
	}

	s, err = d.StringOrHexLiteralEntry("NonExistent")
	if err != nil || s != nil {
		t.Error("StringOrHexLiteralEntry() should return nil for missing key")
	}
}

func TestDictLength(t *testing.T) {
	d := NewDict()

	// No Length entry
	val, objNum := d.Length()
	if val != nil || objNum != nil {
		t.Error("Length() should return nil for missing entry")
	}

	// Direct Length
	d.Insert("Length", Integer(100))
	val, objNum = d.Length()
	if val == nil || *val != 100 || objNum != nil {
		t.Error("Length() failed for direct value")
	}

	// Indirect Length
	d2 := NewDict()
	d2.Insert("Length", IndirectRef{ObjectNumber: Integer(5), GenerationNumber: Integer(0)})
	val, objNum = d2.Length()
	if val != nil || objNum == nil || *objNum != 5 {
		t.Error("Length() failed for indirect reference")
	}
}

func TestDictSpecialEntries(t *testing.T) {
	d := NewDict()
	d.Insert("Type", Name("Page"))
	d.Insert("Subtype", Name("Form"))
	d.Insert("Size", Integer(10))
	d.Insert("Prev", Integer(12345))
	d.Insert("N", Integer(5))
	d.Insert("First", Integer(0))
	d.Insert("W", Array{Integer(1), Integer(2), Integer(3)})
	d.Insert("Index", Array{Integer(0), Integer(10)})
	d.Insert("Linearized", Integer(1))

	if t1 := d.Type(); t1 == nil || *t1 != "Page" {
		t.Error("Type() failed")
	}

	if s := d.Subtype(); s == nil || *s != "Form" {
		t.Error("Subtype() failed")
	}

	if s := d.Size(); s == nil || *s != 10 {
		t.Error("Size() failed")
	}

	if !d.IsPage() {
		t.Error("IsPage() should return true")
	}

	d2 := NewDict()
	d2.Insert("Type", Name("ObjStm"))
	if !d2.IsObjStm() {
		t.Error("IsObjStm() should return true")
	}

	if p := d.Prev(); p == nil || *p != 12345 {
		t.Error("Prev() failed")
	}

	if n := d.N(); n == nil || *n != 5 {
		t.Error("N() failed")
	}

	if f := d.First(); f == nil || *f != 0 {
		t.Error("First() failed")
	}

	if w := d.W(); w == nil || len(w) != 3 {
		t.Error("W() failed")
	}

	if idx := d.Index(); idx == nil || len(idx) != 2 {
		t.Error("Index() failed")
	}

	if !d.IsLinearizationParmDict() {
		t.Error("IsLinearizationParmDict() should return true")
	}
}

func TestDictIncrement(t *testing.T) {
	d := NewDict()
	d.Insert("Count", Integer(5))

	if err := d.Increment("Count"); err != nil {
		t.Fatalf("Increment() error: %v", err)
	}
	if i := d.IntEntry("Count"); i == nil || *i != 6 {
		t.Errorf("Increment() failed, got %v", i)
	}

	if err := d.IncrementBy("Count", 10); err != nil {
		t.Fatalf("IncrementBy() error: %v", err)
	}
	if i := d.IntEntry("Count"); i == nil || *i != 16 {
		t.Errorf("IncrementBy() failed, got %v", i)
	}

	if err := d.Increment("NonExistent"); err == nil {
		t.Error("Increment() should return error for missing key")
	}
}

func TestDictClone(t *testing.T) {
	d := NewDict()
	d.Insert("Int", Integer(1))
	d.Insert("Name", Name("Test"))
	d.Insert("Nested", Dict{"Inner": Integer(2)})
	d.Insert("Nil", nil)

	clone := d.Clone().(Dict)
	if clone.Len() != d.Len() {
		t.Error("Clone() did not preserve length")
	}

	// Modify original
	d.Update("Int", Integer(999))

	// Clone should be unaffected
	if i := clone.IntEntry("Int"); i == nil || *i != 1 {
		t.Error("Clone() did not create independent copy")
	}
}

func TestDictString(t *testing.T) {
	d := NewDict()
	d.Insert("Type", Name("Page"))
	d.Insert("Count", Integer(1))

	s := d.String()
	if !strings.Contains(s, "<<") || !strings.Contains(s, ">>") {
		t.Error("String() should contain dict delimiters")
	}
	if !strings.Contains(s, "Type") || !strings.Contains(s, "Count") {
		t.Error("String() should contain keys")
	}
}

func TestDictPDFString(t *testing.T) {
	d := NewDict()
	d.Insert("Type", Name("Page"))
	d.Insert("Count", Integer(1))
	d.Insert("Bool", Boolean(true))
	d.Insert("Float", Float(1.5))
	d.Insert("Str", StringLiteral("hello"))
	d.Insert("Hex", HexLiteral("48656C6C6F"))
	d.Insert("Null", nil)

	s := d.PDFString()
	if !strings.HasPrefix(s, "<<") || !strings.HasSuffix(s, ">>") {
		t.Error("PDFString() should have dict delimiters")
	}
	if !strings.Contains(s, "/Type/Page") {
		t.Error("PDFString() should contain /Type/Page")
	}
	if !strings.Contains(s, "/Count 1") {
		t.Error("PDFString() should contain /Count 1")
	}
	if !strings.Contains(s, "/Null null") {
		t.Error("PDFString() should contain /Null null")
	}
}

func TestDictStringEntryBytes(t *testing.T) {
	d := NewDict()
	d.Insert("String", StringLiteral("hello"))
	d.Insert("Hex", HexLiteral("48656C6C6F"))

	b, err := d.StringEntryBytes("String")
	if err != nil || string(b) != "hello" {
		t.Error("StringEntryBytes() failed for StringLiteral")
	}

	b, err = d.StringEntryBytes("Hex")
	if err != nil || string(b) != "Hello" {
		t.Error("StringEntryBytes() failed for HexLiteral")
	}

	b, err = d.StringEntryBytes("NonExistent")
	if err != nil || b != nil {
		t.Error("StringEntryBytes() should return nil for missing key")
	}
}
