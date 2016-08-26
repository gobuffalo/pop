package generate

const modelTemplate = `package PACKAGE_NAME
import (
	IMPORTS
)

type MODEL_NAME struct {
	ATTRIBUTES
}

func (x MODEL_NAME) String() string {
	b, _ := json.Marshal(x)
	return string(b)
}

type PLURAL_MODEL_NAME []MODEL_NAME

func (x PLURAL_MODEL_NAME) String() string {
	b, _ := json.Marshal(x)
	return string(b)
}
`

const modelTestTemplate = `package PACKAGE_NAME_test

import "testing"

func Test_MODEL_NAME(t *testing.T) {
	t.Fatal("This test needs to be implemented!")
}
`
