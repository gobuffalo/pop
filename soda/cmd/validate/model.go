package validate

import (
	"fmt"
	"go/ast"
	"regexp"
	"strings"
)

type model struct {
	packages map[string]*ast.Package
	tags     map[string][]string
	rules    []Checker
}

type Checker interface {
	Check(tag string, structName string) []ValidationError
}

type ValidationError struct {
	invalidSymbols string
	field          string
	structName     string
	duplicate      bool
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Invalid symbols '%v' contained in %v.%v", e.invalidSymbols, e.structName, e.field)
}

type Rule struct {
	expr *regexp.Regexp
}

func (r *Rule) Check(tag string, structName string, fieldsCache map[string]bool) []ValidationError  {
	match := r.expr.FindString(tag)
	errorss := []ValidationError{}
	cacheKey :=  strings.Join([]string{structName, tag}, ".")

	if _, exist := fieldsCache[cacheKey]; exist {
		err := ValidationError{"duplicate entry", tag, structName}
		errorss = append(errorss, err)
	}

	fieldsCache[cacheKey] = true

	if len(match) > 0 {
		err := ValidationError{match, tag, structName}
		errorss = append(errorss, err)
	}

	return errorss
}

func NewModel() model  {
	m := model{}
	m.packages = getPackages("models")
	m.tags = getTags("db", m.packages)

	return m
}

func (m *model)Validate() []ValidationError {
	fieldsCache := map[string]bool{}
	errorss := []ValidationError{}
	errs := []ValidationError{}

	if len(m.tags) == 0 {
		return errorss
	}

	rules := createRules()

	for structName, fields := range m.tags {
		for _, field := range fields {
			duplicateErrors := checkForDuplicates(field, structName, fieldsCache)
			errorss = append(errorss, duplicateErrors...)

			for _, ch := range m.rules {
				errs = ch.Check(field, structName)
				errorss = append(errorss, errs...)
			}

			for _, rule := range rules {
				errs = rule.Check(field, structName)
				errorss = append(errorss, errs...)
			}
		}
	}

	return errorss
}

func (m *model) AddRule(ch ...Checker) {
	m.rules = append(m.rules, ch...)
}

func createRules() []Rule {
	return []Rule{
		//tag cannot contain anything except these symbols
		{expr: regexp.MustCompile(`[^a-z0-9_,]+`)},
		//tag cannot end in anything except these symbols
		{expr: regexp.MustCompile(`[^a-z0-9]$`)},
	}
}
