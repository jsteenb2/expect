package be

import (
	"github.com/jsteenb2/expect"
)

// Not is a helper function to negate a matcher.
func Not[T any](matcher expect.Matcher[T]) expect.Matcher[T] {
	return negate(matcher)
}

func negate[T any](matcher expect.Matcher[T]) expect.Matcher[T] {
	return func(got T) expect.MatchResult {
		result := matcher(got)
		return expect.MatchResult{
			Description: "not " + result.Description,
			Matches:     !result.Matches,
		}
	}
}
