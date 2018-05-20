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

//model is a representation of a parsed pop model
//from the syntax tree of a particular go source file
type model struct {
	packages map[string]*ast.Package
	tags     map[string][]*Tag
	processors map[string][]func(tag *Tag) ([]ValidationError, error)
	path string
	allowDuplicates bool
}

//ValidationError is an error created
//by the processors, should one of their validations fail
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

//Here we provide an API with some basic processors
//that will validate the given model tags
//The tags given for the processors will be the tags parsed by the validator
// `*` is a reference to all tags
//If no tags were specified all tags will be parsed and validated
func (m *model) AddDefaultProcessors(tags ...string) {

	if len(tags) == 0 {
		tags = []string{AllTags}
	}

	for _, tagStr := range tags {
		m.processors[tagStr] = append(m.processors[tagStr], func(tag *Tag) ([]ValidationError, error) {
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

//A special case validation
//check if tags have duplicates
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

func (m *model) setPath(path string)  {
	m.path = path
}

//API to set if duplicates are allowed or not
//this will determine whether the duplicates check will be run
func (m *model) SetAllowDuplicates(allowDuplicates bool)  {
	m.allowDuplicates = allowDuplicates
}

//Creates a new validator model
//requires a path to the models folder
func NewValidator(path string) model {
	m := model{}
	m.setPath(path)
	m.processors = map[string][]func(tag *Tag) ([]ValidationError, error){}
	m.allowDuplicates = false

	return m
}

//Provides API to run validation specified tags by the user on specific models
//returns validation errors, if any produced by the processors added by the user
func (m *model) Run(models ...string)  (map[string][]ValidationError, error) {

	m.packages = getPackages(m.path, models...)
	validationErrors := map[string][]ValidationError{}

	if len(m.processors) == 0 {
		return validationErrors, errors.New( "There are no processors to run, consider adding the default ones.")
	}

	tags := []string{}

	for tag, _ := range m.processors {
		tags = append(tags, tag)
	}

	m.tags = getTags(tags, m.packages)

	return  m.validate()
}

func (m *model) validate() (map[string][]ValidationError, error) {
	fieldsCache := map[string]bool{}
	errorss := map[string][]ValidationError{}
	errs := []ValidationError{}
	executableProcessors := []func(tag *Tag) ([]ValidationError, error){}

	if len(m.tags) == 0 {
		return errorss, errors.New("No tags found")
	}

	for structName, fields := range m.tags {
		for _, t := range fields {

			if !m.allowDuplicates {
				duplicateErrors := checkForDuplicates(t, fieldsCache)

				if len(duplicateErrors) > 0 {
					errorss[structName] = append(errorss[structName], duplicateErrors...)
				}
			}

			processors, exists := m.processors[t.GetName()]

			if exists {
				executableProcessors = append(processors)
			}

			globalProcessors, exists := m.processors[AllTags]

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

//User defined processor
//that will validate the given model tags
//The tags given for the processors will be the tags parsed by the validator
// `*` is a reference to all tags
func (m *model) AddProcessor(tag string, processor func(t *Tag) ([]ValidationError, error)) {
	m.processors[tag] = append(m.processors[tag], processor)
}
