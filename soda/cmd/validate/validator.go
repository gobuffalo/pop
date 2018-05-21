package validate

import (
	"fmt"
	"go/ast"
	"strings"
	"regexp"
	"github.com/pkg/errors"
)

const AllTags = "*"

var defaultRegexRules = []*regexp.Regexp{
	//allowed symbols in a tag
	regexp.MustCompile(`[^a-z0-9_,]+`),
	//allowed symbols of the end of a tag
	regexp.MustCompile(`[^a-z0-9]$`),
}

// Validator holds information about the parsed models
type Validator struct {
	packages map[string]*ast.Package
	tags     map[string][]*Tag
	processors map[string][]func(tag *Tag) ([]ValidationError, error)
	path string
	allowDuplicates bool
}

// ValidationError is an error created by the processors, should one of their validations fail.
type ValidationError struct {
	invalidSymbols string
	field          string
	structName     string
	duplicate      bool
	fieldName      string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Invalid symbols '%v' contained in %v.%v.%v", e.invalidSymbols, e.structName, e.fieldName, e.field)
}

// AddDefaultProcessors provides some basic processors that will validate the given model tags.
// The tags given for the processors will be the tags parsed by the validator,`*` is a reference to all tags.
// If no tags were specified all tags will be parsed and validated.
func (v *Validator) AddDefaultProcessors(tags ...string) {

	if len(tags) == 0 {
		tags = []string{AllTags}
	}

	for _, tagStr := range tags {
		v.processors[tagStr] = append(v.processors[tagStr], func(tag *Tag) ([]ValidationError, error) {
			errorss := []ValidationError{}

			for _, rexpr := range defaultRegexRules {
				match := rexpr.FindString(tagStr)

				if len(match) > 0 {
					err := ValidationError{
						match,
						tag.value,
						tag.structName,
						false,
						tag.GetName(),
					}
					errorss = append(errorss, err)
				}
			}

			return errorss, nil
		})
	}
}

// checkForDuplicates validates duplicate tag values
func checkForDuplicates(t *Tag, fieldsCache map[string]bool) []ValidationError {
	errorss := []ValidationError{}
	cacheKey := strings.Join([]string{t.structName, t.name, t.value}, ".")

	if _, exist := fieldsCache[cacheKey]; exist {
		err := ValidationError{"duplicate entry", t.GetValue(), t.GetStructName(), true, t.GetName()}
		errorss = append(errorss, err)
	}

	fieldsCache[cacheKey] = true

	return errorss
}

func (v *Validator) setPath(path string)  {
	v.path = path
}

// SetAllowDuplicates sets a flag if duplicates are allowed or not.
// This will determine whether the duplicates check will be run.
func (v *Validator) SetAllowDuplicates(allowDuplicates bool)  {
	v.allowDuplicates = allowDuplicates
}

// NewValidator creates a new validator model.
// It requires a path to the models folder.
func NewValidator(path string) Validator {
	m := Validator{}
	m.setPath(path)
	m.processors = map[string][]func(tag *Tag) ([]ValidationError, error){}
	m.allowDuplicates = false

	return m
}

// Run  will validate specified tags on all models, if none were passed.
// It returns validation errors, if any produced by the processor.
func (v *Validator) Run(models ...string)  (map[string][]ValidationError, error) {

	v.packages = getPackages(v.path, models...)
	validationErrors := map[string][]ValidationError{}

	if len(v.processors) == 0 {
		return validationErrors, errors.New( "There are no processors to run, consider adding the default ones.")
	}

	tags := []string{}

	for tag, _ := range v.processors {
		tags = append(tags, tag)
	}

	v.tags = getTags(tags, v.packages)

	return  v.validate()
}

func (v *Validator) validate() (map[string][]ValidationError, error) {
	fieldsCache := map[string]bool{}
	errorss := map[string][]ValidationError{}
	errs := []ValidationError{}
	executableProcessors := []func(tag *Tag) ([]ValidationError, error){}

	if len(v.tags) == 0 {
		return errorss, errors.New("No tags found")
	}

	for structName, fields := range v.tags {
		for _, t := range fields {

			if !v.allowDuplicates {
				duplicateErrors := checkForDuplicates(t, fieldsCache)

				if len(duplicateErrors) > 0 {
					errorss[structName] = append(errorss[structName], duplicateErrors...)
				}
			}

			processors, exists := v.processors[t.GetName()]

			if exists {
				executableProcessors = append(processors)
			}

			globalProcessors, exists := v.processors[AllTags]

			if exists {
				executableProcessors = append(executableProcessors, globalProcessors...)
			}

			for _, processor := range executableProcessors {
				errs, _ = processor(t)
				if len(errs) > 0 {
					errorss[structName] = append(errorss[structName], errs...)
				}
			}
		}
	}

	return errorss, nil
}

// AddProcessor adds a processor that will validate the given model tags
// The tags given for the processors will be the tags parsed by the validator where `*` is a reference to all tags
func (v *Validator) AddProcessor(tag string, processor func(t *Tag) ([]ValidationError, error)) {
	v.processors[tag] = append(v.processors[tag], processor)
}
