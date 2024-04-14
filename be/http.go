package be

import (
	"fmt"
	"io"
	"net/http"
	
	"github.com/jsteenb2/expect"
)

const (
	subjectNameHTTPResp     = "the response"
	responseBodySubjectName = subjectNameHTTPResp + " body"
)

// HaveStatus returns a matcher that checks if the response status code is equal to the given status code.
func HaveStatus(status int) expect.Matcher[*http.Response] {
	return func(res *http.Response) expect.MatchResult {
		return expect.MatchResult{
			Description: fmt.Sprintf("have status of %d", status),
			But:         fmt.Sprintf("it was %d", res.StatusCode),
			Matches:     res.StatusCode == status,
			SubjectName: subjectNameHTTPResp,
		}
	}
}

// BeOK is a convenience matcher for HaveStatus(http.StatusOK).
func BeOK(res *http.Response) expect.MatchResult {
	return HaveStatus(http.StatusOK)(res)
}

// HaveHeader returns a matcher that checks if the response has a header with the given name and value.
func HaveHeader(header, value string) expect.Matcher[*http.Response] {
	return func(res *http.Response) expect.MatchResult {
		return expect.MatchResult{
			Description: fmt.Sprintf("have header %q of %q", header, value),
			Matches:     res.Header.Get(header) == value,
			But:         fmt.Sprintf("it was %q", res.Header.Get(header)),
			SubjectName: subjectNameHTTPResp,
		}
	}
}

// HaveJSONHeader is a convenience matcher for HaveHeader("content-type", "application/json").
func HaveJSONHeader(res *http.Response) expect.MatchResult {
	return HaveHeader("content-type", "application/json")(res)
}

// HaveBody returns a matcher that checks if the response body meets the given matchers' criteria. Note this will read the entire body using io.ReadAll.
func HaveBody(bodyMatchers expect.Matcher[io.Reader]) expect.Matcher[*http.Response] {
	return func(res *http.Response) expect.MatchResult {
		result := bodyMatchers(res.Body)
		result.SubjectName = responseBodySubjectName
		return result
	}
}
