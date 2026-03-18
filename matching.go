package expect

import (
	"errors"
	"fmt"
	"strings"
)

type (
	// TB is a cut-down version of testing.TB.
	TB interface {
		Error(args ...any)
		Errorf(format string, args ...any)
		Fatalf(format string, args ...any)
		Helper()
	}
	Inspector[T any] struct {
		t       TB
		Subject T
	}
)

// It is the entry point for the matcher DSL. Pass in the testing.TB and the subject you want to test.
func It[T any](t TB, subject T) Inspector[T] {
	return Inspector[T]{t, subject}
}

// NoError is a helper function that will call t.Fatalf if the error is not nil.
func NoError(t TB, err error) {
	t.Helper()
	if err == nil {
		// hooray no error
		return
	}
	format, args := "unexpected error: %v", []any{err}
	if fieldFormat, fieldArg, ok := getErrFieldsStr(err); ok {
		format += "\n" + fieldFormat
		args = append(args, fieldArg)
	}
	t.Fatalf(format, args...)
}

func getErrFieldsStr(err error) (string, string, bool) {
	fielder, ok := err.(interface{ Fields() []any })
	if !ok && fielder == nil || len(fielder.Fields()) == 0 {
		return "", "", false
	}
	
	fields := fielder.Fields()
	
	var sb strings.Builder
	for i := 0; i < len(fields); i += 2 {
		k := fields[i]
		var v any
		if vIdx := i + 1; vIdx < len(fields) {
			v = fields[vIdx]
		}
		if s, _ := k.(string); s == "stack_trace" {
			fmt.Fprintf(&sb, "\n\t%s=[", s)
			vs, ok := v.([]string)
			if ok {
				for _, vf := range vs {
					sb.WriteString("\n\t\t")
					fmt.Fprintf(&sb, "%s,", vf)
				}
			}
			if len(vs) > 0 {
				sb.WriteString("\n\t")
			}
			sb.WriteRune(']')
			continue
		}
		fmt.Fprintf(&sb, "\n\t%v=%v", k, v)
	}
	
	return "fields: %s", sb.String(), true
}

// Error is a helper function that will call t.Fatalf if the error is nil.
func Error(t TB, err error) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected an error")
	}
}

// ErrorOfType is a helper function that will call t.Fatalf if the error is not of the expected type.
func ErrorOfType(t TB, err error, expectedType error) {
	t.Helper()
	if !errors.Is(err, expectedType) {
		t.Fatalf("expected error of type %T, but got %q", expectedType, err.Error())
	}
}

// To is the method that actually runs the matchers. It will call Errorf on the testing.TB if any of the matchers fail.
func (e Inspector[T]) To(matchers ...Matcher[T]) {
	e.t.Helper()
	for _, matcher := range matchers {
		result := matcher(e.Subject)
		
		if result.SubjectName == "" {
			result.SubjectName = calculateSubjectName(e)
		}
		
		if !result.Matches {
			e.t.Error(result.Error())
		}
	}
}

func calculateSubjectName[T any](e Inspector[T]) string {
	var subjectName = fmt.Sprintf("%+v", e.Subject)
	
	if str, isStringer := any(e.Subject).(fmt.Stringer); isStringer {
		subjectName = str.String()
	}
	return subjectName
}
