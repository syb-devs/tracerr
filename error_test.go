package tracerr_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/syb-devs/tracerr"
)

func fails() error {
	f := func() error {
		return errors.New("Whoops!")
	}
	return f()
}

func failsToo() error {
	return tracerr.WrapError(fails(), "Oh my God! This is going to explode!!")
}

func failsAgain() error {
	return tracerr.WrapError(failsToo(), "something went really wrong")
}

func middleFunc() error {
	return failsAgain()
}

func first() error {
	return tracerr.WrapError(middleFunc(), "first failed")
}

func TestWrapError(t *testing.T) {
	patterns := []string{
		"wrapper error: first failed",
		"wrapper error: something went really wrong",
		"wrapper error: Oh my God! This is going to explode!!",
		"/error_test.go:19",
		"/error_test.go:23",
		"/error_test.go:27",
		"/error_test.go:31",
		"/error_test.go:47",
		"trigger error: Whoops!",
	}

	err := first()
	trace := err.(*tracerr.WrappedError).TraceString()
	for _, pattern := range patterns {
		i := strings.Index(trace, pattern)
		if i == -1 {
			t.Error("pattern not found in stack trace string")
			t.Errorf("want:\n %s", pattern)
			t.Errorf("stack trace string:\n %s", trace)
		}
	}
}
