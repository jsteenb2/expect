package spytb

import (
	"fmt"
	
	"github.com/jsteenb2/expect"
)

func Error(s string) expect.Matcher[*expect.SpyTB] {
	return func(tb *expect.SpyTB) expect.MatchResult {
		found := false
		for _, e := range tb.ErrorCalls {
			if e == s {
				found = true
			}
		}
		return expect.MatchResult{
			Description: fmt.Sprintf("have error %q", s),
			Matches:     found,
			But:         fmt.Sprintf("has %v", tb.ErrorCalls),
		}
	}
}

func NoErrors(tb *expect.SpyTB) expect.MatchResult {
	return expect.MatchResult{
		Description: "have no errors",
		Matches:     len(tb.ErrorCalls) == 0,
		But:         fmt.Sprintf("it had errors %v", tb.ErrorCalls),
	}
}
