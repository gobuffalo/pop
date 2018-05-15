package validate

import (
	"fmt"
	"go/ast"
	"strings"
	"regexp"
	"github.com/pkg/errors"
)

type model struct {
	packages map[string]*ast.Package
	tags     map[string][]*Tag
	processors map[string][]func(tag *Tag) ([]ValidationError, error)
	path string
}

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

func (m *model) AddDefaultProcessors(tags ...string) {

	if m.processors == nil {
		m.processors =  map[string][]func(tag *Tag) ([]ValidationError, error){}
	}

	if len(tags) == 0 {
		tags = []string{"db"}
	}

	for _, tagStr := range tags {
		m.processors[tagStr] = append(m.processors[tagStr], func(tag *Tag) ([]ValidationError, error) {
			regexpr := []*regexp.Regexp{
				//allowed symbols in a tag
				regexp.MustCompile(`[^a-z0-9_,]+`),
				//allowed symbols of the end of a tag
				regexp.MustCompile(`[^a-z0-9]$`),
			}
			errorss := []ValidationError{}

			for _, rexpr := range regexpr {
				match := rexpr.FindString(tagStr)

				if len(match) > 0 {
					err := ValidationError{
						match,
						tag.value,
						tag.structName,
						false,
						tag.getName(),
					}
					errorss = append(errorss, err)
				}
			}

			return errorss, nil
		})
	}
}

func checkForDuplicates(t *Tag, fieldsCache map[string]bool) []ValidationError {
	errorss := []ValidationError{}
	cacheKey := strings.Join([]string{t.structName, t.name, t.value}, ".")

	if _, exist := fieldsCache[cacheKey]; exist {
		err := ValidationError{"duplicate entry", t.getValue(), t.structName, true, t.getName()}
		errorss = append(errorss, err)
	}

	fieldsCache[cacheKey] = true

	return errorss
}

func (m *model) setPath(path string)  {
	m.path = path
}

func NewValidator(path string) model {
	m := model{}
	m.setPath(path)

	return m
}

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

	if len(m.tags) == 0 {
		return errorss, errors.New("No tags found")
	}

	for structName, fields := range m.tags {
		for _, t := range fields {
			duplicateErrors := checkForDuplicates(t, fieldsCache)

			if len(duplicateErrors) > 0 {
				errorss[structName] = append(errorss[structName], duplicateErrors...)
			}

			for _, processorTag := range m.processors {
				for _, processor := range processorTag {
					errs, _ = processor(t)
					if len(errs) > 0 {
						errorss[structName] = append(errorss[structName], errs...)
					}
				}
			}
		}
	}

	return errorss, nil
}

func (m *model) AddProcessor(tag string, processor func(t *Tag) ([]ValidationError, error)) {
	m.processors[tag] = append(m.processors[tag], processor)
}
