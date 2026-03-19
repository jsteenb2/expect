package spytb

import (
	"fmt"
	"strings"

	"github.com/jsteenb2/expect"
)

func Error(s string) expect.Matcher[*expect.SpyTB] {
	return func(tb *expect.SpyTB) expect.MatchResult {
		found := false
		for _, e := range tb.ErrorCalls {
			if strings.Contains(stripStackTrace(e), s) {
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

func HaveMatchResult(expected expect.MatchResult) expect.Matcher[expect.MatchResult] {
	return func(got expect.MatchResult) expect.MatchResult {
		matches := got.Description == expected.Description &&
			got.Matches == expected.Matches &&
			got.But == expected.But &&
			got.SubjectName == expected.SubjectName
		return expect.MatchResult{
			Description: fmt.Sprintf("have match result with description %q", expected.Description),
			Matches:     matches,
			But:         fmt.Sprintf("got %+v", got),
		}
	}
}

func stripStackTrace(err string) string {
	idx := strings.Index(err, "Error Trace:")
	if idx == -1 {
		return err
	}
	return strings.TrimSpace(err[:idx])
}
