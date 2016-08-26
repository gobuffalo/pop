package generate

const modelTemplate = `package PACKAGE_NAME

import (
	IMPORTS
)

type MODEL_NAME struct {
	ATTRIBUTES
}

func (CHAR MODEL_NAME) String() string {
	b, _ := json.Marshal(CHAR)
	return string(b)
}

type PLURAL_MODEL_NAME []MODEL_NAME

func (CHAR PLURAL_MODEL_NAME) String() string {
	b, _ := json.Marshal(CHAR)
	return string(b)
}
`

const modelTestTemplate = `package PACKAGE_NAME_test

import "testing"

func Test_MODEL_NAME(t *testing.T) {
	t.Fatal("This test needs to be implemented!")
}
`
