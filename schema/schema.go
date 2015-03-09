package schema

import "errors"

type SchemaDoc struct {
	Properties map[string]interface{} `json:"properties"`
}

type ParamDoc struct {
	Name   string
	Values map[string]interface{}
}

type ValidatorDoc struct {
	Key   string
	Value interface{}
}

type Validator interface {
	Validate(interface{}) error
	Doc() ValidatorDoc
}

type RequiredValidator struct{}

func (val RequiredValidator) Validate(v interface{}) error {
	err := errors.New("required")
	switch t := v.(type) {
	case string:
		if t != "" {
			return nil
		}
	}

	return err
}

func (v RequiredValidator) Doc() ValidatorDoc {
	return ValidatorDoc{Key: "required", Value: true}
}

type param struct {
	name       string
	typ        string
	validators []Validator
}

func (p *param) String() *param {
	p.typ = "string"
	return p
}

func (p *param) Required() *param {
	p.validators = append(p.validators, RequiredValidator{})
	return p
}

func (p *param) Validate(v interface{}) []error {
	var result []error

	for _, validator := range p.validators {
		if err := validator.Validate(v); err != nil {
			result = append(result, err)
		}
	}

	return result
}

func (p *param) Doc() ParamDoc {
	doc := ParamDoc{Name: p.name, Values: map[string]interface{}{}}

	for _, validator := range p.validators {
		valDoc := validator.Doc()
		doc.Values[valDoc.Key] = valDoc.Value
	}

	return doc
}

type schema struct {
	params []*param
}

func (s *schema) Validate(data map[string]interface{}) []error {
	var result []error

	for _, param := range s.params {
		errs := param.Validate(data[param.name])
		result = append(result, errs...)
	}

	return result
}

func (s *schema) Doc() SchemaDoc {
	doc := SchemaDoc{Properties: map[string]interface{}{}}

	for _, param := range s.params {
		paramDoc := param.Doc()
		doc.Properties[paramDoc.Name] = paramDoc.Values
	}

	return doc
}

func New(ps ...*param) *schema {
	return &schema{ps}
}

func Param(name string) *param {
	return &param{name: name}
}
