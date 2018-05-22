package validate

import (
	"fmt"
	"go/ast"
	"strings"
	"regexp"
	"github.com/pkg/errors"
	vv "github.com/gobuffalo/validate"
)

const AllTags = "*"

var defaultRegexRules = map[string]*regexp.Regexp{
	//allowed symbols in a tag
	"Invalid symboles %v in %v.%v.%v": regexp.MustCompile(`[^a-z0-9_, ]+`),
	//allowed symbols of the end of a tag
	"Tag cannot end on %v in  %v.%v.%v" : regexp.MustCompile(`[^a-z0-9]$`),
}

// Validator holds information about the parsed models
type Validator struct {
	packages map[string]*ast.Package
	tags     map[string][]*Tag
	processors map[string][]func(tag *Tag, errors *vv.Errors)
	path string
	allowDuplicates bool
}

// AddDefaultProcessors provides some basic processors that will validate the given model tags.
// The tags given for the processors will be the tags parsed by the validator,`*` is a reference to all tags.
// If no tags were specified all tags will be parsed and validated.
func (v *Validator) AddDefaultProcessors(tags ...string) {

	if len(tags) == 0 {
		tags = []string{AllTags}
	}

	for _, tagStr := range tags {
		v.processors[tagStr] = append(v.processors[tagStr], func(tag *Tag, errors *vv.Errors) {
			for msg, rexpr := range defaultRegexRules {
				match := rexpr.FindString(tag.GetValue())

				if len(match) > 0 {
					errors.Add(tag.GetStructName(), fmt.Sprintf(msg, match, tag.GetStructName(), tag.GetName(), tag.GetValue()))
				}
			}
		})

		v.processors[tagStr] = append(v.processors[tagStr], func(tag *Tag, errors *vv.Errors) {
			if len(tag.GetValue()) == 0 {
				errors.Add(tag.GetStructName(), fmt.Sprintf("Tag cannot be empty %v.%v", tag.GetStructName(), tag.GetName()))
			}
		})
	}
}

// checkForDuplicates validates duplicate tag values
func checkForDuplicates(errors *vv.Errors, t *Tag, fieldsCache map[string]bool) {
	cacheKey := strings.Join([]string{t.structName, t.name, t.value}, ".")

	if _, exist := fieldsCache[cacheKey]; exist {
		errors.Add(t.GetStructName(), fmt.Sprintf("Duplicate tag value %v in %v.%v", t.GetValue(), t.GetStructName(), t.GetName()))
	}

	fieldsCache[cacheKey] = true
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
	m.processors = map[string][]func(tag *Tag, errors *vv.Errors){}
	m.allowDuplicates = false

	return m
}

// Run  will validate specified tags on all models, if none were passed.
// It returns validation errors, if any produced by the processor.
func (v *Validator) Run(models ...string)  (*vv.Errors, error) {

	v.packages = getPackages(v.path, models...)
	validationErrors := vv.NewErrors()

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

func (v *Validator) validate() (*vv.Errors, error) {
	fieldsCache := map[string]bool{}
	errorss := vv.NewErrors()

	if len(v.tags) == 0 {
		return errorss, errors.New("No tags found")
	}

	for _, fields := range v.tags {
		for _, t := range fields {
			executableProcessors := []func(tag *Tag, errors *vv.Errors) {}

			if !v.allowDuplicates {
				errorss.Append(
					vv.Validate(&TagDuplicateValidator{
					t,
					fieldsCache,
					checkForDuplicates,
				}))
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
				errorss.Append(
					vv.Validate(
						&TagValidator{
					t,
					processor,
				}))
			}
		}
	}

	return errorss, nil
}


// AddProcessor adds a processor that will validate the given model tags
// The tags given for the processors will be the tags parsed by the validator where `*` is a reference to all tags
func (v *Validator) AddProcessor(tag string, processor func(t *Tag, errors *vv.Errors)) {
	v.processors[tag] = append(v.processors[tag], processor)
}
