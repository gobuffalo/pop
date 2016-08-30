package generate

const modelTemplate = `package PACKAGE_NAME

import (
	IMPORTS
)

type MODEL_NAME struct {
	ATTRIBUTES
}

// String is not required by pop and may be deleted
func (CHAR MODEL_NAME) String() string {
	b, _ := json.Marshal(CHAR)
	return string(b)
}

// PLURAL_MODEL_NAME is not required by pop and may be deleted
type PLURAL_MODEL_NAME []MODEL_NAME

// String is not required by pop and may be deleted
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
