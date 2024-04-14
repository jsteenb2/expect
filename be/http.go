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

// HTTPStatus returns a matcher that checks if the response status code is equal to the given status code.
func HTTPStatus(status int) expect.Matcher[*http.Response] {
	return func(res *http.Response) expect.MatchResult {
		return expect.MatchResult{
			Description: fmt.Sprintf("have status of %d", status),
			But:         fmt.Sprintf("it was %d", res.StatusCode),
			Matches:     res.StatusCode == status,
			SubjectName: subjectNameHTTPResp,
		}
	}
}

// HTTPHeader returns a matcher that checks if the response has a header with the given name and value.
func HTTPHeader(header, value string) expect.Matcher[*http.Response] {
	return func(res *http.Response) expect.MatchResult {
		return expect.MatchResult{
			Description: fmt.Sprintf("have header %q of %q", header, value),
			Matches:     res.Header.Get(header) == value,
			But:         fmt.Sprintf("it was %q", res.Header.Get(header)),
			SubjectName: subjectNameHTTPResp,
		}
	}
}

// ContentTypeJSONHeader is a convenience matcher for HTTPHeader("content-type", "application/json").
func ContentTypeJSONHeader(res *http.Response) expect.MatchResult {
	return HTTPHeader("content-type", "application/json")(res)
}

// HTTPRespBody returns a matcher that checks if the response body meets the given matchers' criteria. Note this will read the entire body using io.ReadAll.
func HTTPRespBody(bodyMatchers expect.Matcher[io.Reader]) expect.Matcher[*http.Response] {
	return func(res *http.Response) expect.MatchResult {
		result := bodyMatchers(res.Body)
		result.SubjectName = responseBodySubjectName
		return result
	}
}
