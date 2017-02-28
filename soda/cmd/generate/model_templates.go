package generate

const modelTemplate = `package {{package_name}}

import (
	{{#each model.Imports as |i|}}"{{i}}"
	{{/each}}
)

type {{model_name}} struct {
	{{#each model.Attributes as |a|}}{{a}}
	{{/each}}
}

// String is not required by pop and may be deleted
func ({{char}} {{model_name}}) String() string {
	j{{char}}, _ := json.Marshal({{char}})
	return string(j{{char}})
}

// {{plural_model_name}} is not required by pop and may be deleted
type {{plural_model_name}} []{{model_name}}

// String is not required by pop and may be deleted
func ({{char}} {{plural_model_name}}) String() string {
	j{{char}}, _ := json.Marshal({{char}})
	return string(j{{char}})
}

// Validate gets run everytime you call a "pop.Validate" method.
// This method is not required and may be deleted.
func ({{char}} *{{model_name}}) Validate(tx *pop.Connection) (*validate.Errors, error) {
	{{#if model.ValidatableAttributes }}
	return validate.Validate(
		{{#each model.ValidatableAttributes as |a|}}
		&validators.{{capitalize a.GoType}}IsPresent{Field: {{char}}.{{a.Names.Proper}}, Name: "{{a.Names.Proper}}"},{{/each}}
	), nil
	{{ else }}
		return validate.NewErrors(), nil
	{{/if}}
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
// This method is not required and may be deleted.
func ({{char}} *{{model_name}}) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
// This method is not required and may be deleted.
func ({{char}} *{{model_name}}) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
`

const modelTestTemplate = `package {{package_name}}_test

import "testing"

func Test_{{model_name}}(t *testing.T) {
	t.Fatal("This test needs to be implemented!")
}
`
