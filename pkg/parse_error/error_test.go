package parseerror

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	e := New("thing", "name")

	if e == nil {
		t.Errorf("error was not correcly created")
	}
}

func TestFromError(t *testing.T) {
	e := FromError(nil, "tea")
	if e != nil {
		t.Errorf("error should be nil, if error is nil")
	}

	e = FromError(errors.New("thing"), "123")
	if e == nil {
		t.Errorf("error was not correcly created")
	}
}
