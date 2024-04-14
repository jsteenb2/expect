package be

import (
	"fmt"
	"strings"
	
	"github.com/jsteenb2/expect"
)

// HaveLength will check a string's length meets the given matcher's criteria.
func HaveLength(matcher expect.Matcher[int]) expect.Matcher[string] {
	return func(in string) expect.MatchResult {
		result := matcher(len(in))
		result.Description = fmt.Sprintf("have length %v", result.Description)
		return result
	}
}

// HaveAllCaps will check if a string is in all caps.
func HaveAllCaps(in string) expect.MatchResult {
	return expect.MatchResult{
		Description: "in all caps",
		Matches:     strings.ToUpper(in) == in,
		But:         "it was not in all caps",
	}
}

// HaveSubstring will check if a string contains a given substring.
func HaveSubstring(substring string) expect.Matcher[string] {
	return func(in string) expect.MatchResult {
		return expect.MatchResult{
			Description: fmt.Sprintf("contain %q", substring),
			Matches:     strings.Contains(in, substring),
		}
	}
}
