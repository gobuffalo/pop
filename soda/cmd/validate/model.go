package validate

import (
	"regexp"
	"go/ast"
	"strings"
	"fmt"
)

type model struct{
	packages map[string]*ast.Package
	tags  map[string][]string
	rules []Checker
}

type Checker interface {
	Check(tag, structName string)[]ValidationError
}

type ValidationError struct {
	invalidSymbols string
	field string
	structName string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Invalid symbols '%v' contained in %v.%v", e.invalidSymbols, e.structName, e.field)
}

type Rule struct {
	expr *regexp.Regexp
}

var fieldsCache map[string]bool

func (r *Rule) Check(tag, structName string) []ValidationError  {
	matches := r.expr.FindAllString(tag, -1)
	errors := []ValidationError{}
	cacheKey :=  strings.Join([]string{structName, tag}, ".")

	if _, exist := fieldsCache[cacheKey]; exist {
		err := ValidationError{"duplicate entry", tag, structName}
		return append(errors, err)
	}

	fieldsCache[cacheKey] = true

	for _, match := range matches {
		if len(match) > 0 {
			err := ValidationError{match, tag, structName}
			errors = append(errors, err)
		}
	}

	return errors
}

func NewModel() model  {
	m := model{}
	m.packages = getPackages("models")
	m.tags = getTags("db", m.packages)

	return m
}

func (m *model)Validate() []ValidationError {
	errors := []ValidationError{}
	errs := []ValidationError{}


	if len(m.tags) == 0 {
		return errors
	}

	rules := createRules()

	for structName, fields := range m.tags {
		for _, field := range fields {
			for _, ch := range m.rules {
				errs = ch.Check(field, structName)
				errors = append(errors, errs...)
			}

			for _, rule := range rules {
				errs = rule.Check(field, structName)
				errors = append(errors, errs...)
			}
		}
	}

	return errors
}

func (m *model)AddRule(ch ...Checker)  {
	m.rules = append(m.rules, ch...)
}

func createRules() []Rule {
	return []Rule{
		{expr:regexp.MustCompile(`[^a-z0-9_]+`)},
	}
}
