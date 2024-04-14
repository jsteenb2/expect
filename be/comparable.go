package be

import (
	"cmp"
	"fmt"
	
	"github.com/jsteenb2/expect"
)

// Eq checks if a value is equal to another value.
func Eq[T comparable](expected T) expect.Matcher[T] {
	return func(got T) expect.MatchResult {
		description := fmt.Sprintf("be equal to %+v", expected)
		but := fmt.Sprintf("it was %v", got)
		subject := ""
		
		if str, isStr := any(got).(string); isStr {
			description = fmt.Sprintf("be equal to %q", any(expected).(string))
			but = fmt.Sprintf("it was %q", str)
			subject = fmt.Sprintf("%q", str)
		}
		
		return expect.MatchResult{
			Description: description,
			Matches:     got == expected,
			But:         but,
			SubjectName: subject,
		}
	}
}

// LessThan checks if a value is less than another value.
func LessThan[T cmp.Ordered](in T) expect.Matcher[T] {
	return func(got T) expect.MatchResult {
		return expect.MatchResult{
			Description: fmt.Sprintf("be less than %v", in),
			Matches:     got < in,
			But:         fmt.Sprintf("it was %v", got),
		}
	}
}

// GreaterThan checks if a value is greater than another value.
func GreaterThan[T cmp.Ordered](in T) expect.Matcher[T] {
	return func(got T) expect.MatchResult {
		return expect.MatchResult{
			Description: fmt.Sprintf("be greater than %v", in),
			Matches:     got > in,
			But:         fmt.Sprintf("it was %v", got),
		}
	}
}
