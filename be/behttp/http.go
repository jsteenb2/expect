package behttp

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

// Status returns a matcher that checks if the response status code is equal to the given status code.
func Status(status int) expect.Matcher[*http.Response] {
	return func(res *http.Response) expect.MatchResult {
		return expect.MatchResult{
			Description: fmt.Sprintf("have status of %d", status),
			But:         fmt.Sprintf("it was %d", res.StatusCode),
			Matches:     res.StatusCode == status,
			SubjectName: subjectNameHTTPResp,
		}
	}
}

// StatusOK provides a helper for Status(200).
func StatusOK() expect.Matcher[*http.Response] {
	return Status(http.StatusOK)
}

// Header returns a matcher that checks if the response has a header with the given name and value.
func Header(header, value string) expect.Matcher[*http.Response] {
	return func(res *http.Response) expect.MatchResult {
		return expect.MatchResult{
			Description: fmt.Sprintf("have header %q of %q", header, value),
			Matches:     res.Header.Get(header) == value,
			But:         fmt.Sprintf("it was %q", res.Header.Get(header)),
			SubjectName: subjectNameHTTPResp,
		}
	}
}

// ContentType is a convenience matcher for Header("content-type", want).
func ContentType(want string) expect.Matcher[*http.Response] {
	return Header("Content-Type", want)
}

// ContentTypeJSON is a convenience matcher for Header("content-type", "application/json").
func ContentTypeJSON(res *http.Response) expect.MatchResult {
	return ContentType("application/json")(res)
}

// RespBody returns a matcher that checks if the response body meets the given matchers' criteria. Note this will read the entire body using io.ReadAll.
func RespBody(bodyMatchers expect.Matcher[io.Reader]) expect.Matcher[*http.Response] {
	return func(res *http.Response) expect.MatchResult {
		result := bodyMatchers(res.Body)
		result.SubjectName = responseBodySubjectName
		return result
	}
}
