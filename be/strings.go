package be

import (
	"fmt"
	"strings"
	
	"github.com/jsteenb2/expect"
)

// Len will check a string's length meets the given matcher's criteria.
// TODO: would be nice to generalize this a bit more so that it works with
//		 any slice. Finding it awkward with the existing matcher type matching.
//		 Need to noodle on it a bit.
func Len(matcher expect.Matcher[int]) expect.Matcher[string] {
	return func(in string) expect.MatchResult {
		result := matcher(len(in))
		result.Description = fmt.Sprintf("have length %v", result.Description)
		return result
	}
}

// AllCaps will check if a string is in all caps.
func AllCaps(in string) expect.MatchResult {
	return expect.MatchResult{
		Description: "in all caps",
		Matches:     strings.ToUpper(in) == in,
		But:         "it was not in all caps",
	}
}

// Substring will check if a string contains a given substring.
func Substring(substring string) expect.Matcher[string] {
	return func(in string) expect.MatchResult {
		return expect.MatchResult{
			Description: fmt.Sprintf("contain %q", substring),
			Matches:     strings.Contains(in, substring),
		}
	}
}
