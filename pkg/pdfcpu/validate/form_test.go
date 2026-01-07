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

package validate

import (
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

func TestFormFieldMissingFT(t *testing.T) {
	// Reproduces #1274: Widget in AcroForm Fields without FT should validate
	// (Preview and pdf.js accept this structure)

	xRefTable, err := pdfcpu.CreateXRefTableWithRootDict()
	if err != nil {
		t.Fatal(err)
	}

	// Widget without FT (field type) entry
	widgetDict := types.Dict{
		"Subtype": types.Name("Widget"),
		"Rect": types.Array{
			types.Integer(0), types.Integer(0),
			types.Integer(100), types.Integer(100),
		},
	}
	widgetRef, err := xRefTable.IndRefForNewObject(widgetDict)
	if err != nil {
		t.Fatal(err)
	}

	acroFormDict := types.Dict{
		"Fields": types.Array{*widgetRef},
	}
	acroFormRef, err := xRefTable.IndRefForNewObject(acroFormDict)
	if err != nil {
		t.Fatal(err)
	}

	rootDict, err := xRefTable.Catalog()
	if err != nil {
		t.Fatal(err)
	}
	rootDict.Insert("AcroForm", *acroFormRef)

	// Should validate successfully (matches Preview/pdf.js behavior)
	if err := validateForm(xRefTable, rootDict, false, model.V10); err != nil {
		t.Errorf("Widget without FT should validate: %v", err)
	}
}
