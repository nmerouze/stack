package schema_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"../schema"
)

func TestValidate(t *testing.T) {
	d := map[string]interface{}{"name": "sencha"}
	s := schema.New(
		schema.Param("name").Required(),
	)

	if err := s.Validate(d); err != nil {
		t.Fatalf("No errors expected, got: %#v", err)
	}
}

func TestValidateError(t *testing.T) {
	d := map[string]interface{}{}
	s := schema.New(
		schema.Param("name").Required(),
	)

	var errs []error
	errs = s.Validate(d)

	if len(errs) != 1 {
		t.Fatalf("1 error expected, got: %d", len(errs))
	}

	if errs[0].Error() != "required" {
		t.Fatalf(`"required" error expected, got: %d`, errs[0].Error())
	}
}

func TestDoc(t *testing.T) {
	s := schema.New(
		schema.Param("name").Required(),
	)

	result := new(bytes.Buffer)
	json.NewEncoder(result).Encode(s.Doc())

	if result.String() != `{"properties":{"name":{"required":true}}}`+"\n" {
		t.Fatalf("JSON schema not generated correctly: %#v", result.String())
	}
}

func TestRequiredValidatorError(t *testing.T) {
	v := schema.RequiredValidator{}

	if err := v.Validate(nil); err == nil {
		t.Fatalf("Validation of nil should return an error")
	}

	if err := v.Validate(""); err == nil {
		t.Fatalf("Validation of an empty string should return an error")
	}
}
