package be

import (
	"bytes"
	"fmt"
	"io"
	
	"github.com/jsteenb2/expect"
)

// ContainingByte will check if the given byte slice is contained in the byte slice.
func ContainingByte(want []byte) expect.Matcher[[]byte] {
	return func(have []byte) expect.MatchResult {
		return expect.MatchResult{
			Description: fmt.Sprintf("contain %q", want),
			Matches:     bytes.Contains(have, want),
			SubjectName: "the reader",
			But:         fmt.Sprintf("it didn't have %q", want),
		}
	}
}

// ContainingString will check if the given string is contained in the byte slice.
func ContainingString(want string) expect.Matcher[[]byte] {
	return func(have []byte) expect.MatchResult {
		return expect.MatchResult{
			Description: fmt.Sprintf("contain %q", want),
			Matches:     bytes.Contains(have, []byte(want)),
			SubjectName: "the reader",
			But:         fmt.Sprintf("it was %q", have),
		}
	}
}

// HaveData will read all the data from the io.Reader and run the given matcher on it.
func HaveData(matcher expect.Matcher[[]byte]) expect.Matcher[io.Reader] {
	return func(reader io.Reader) expect.MatchResult {
		all, err := io.ReadAll(reader)
		if err != nil {
			return expect.MatchResult{
				Description: "have data in io.Reader",
				Matches:     false,
				But:         "it could not be read",
			}
		}
		return matcher(all)
	}
}

func HaveString(matcher expect.Matcher[string]) expect.Matcher[io.Reader] {
	return func(reader io.Reader) expect.MatchResult {
		all, err := io.ReadAll(reader)
		if err != nil {
			return expect.MatchResult{
				Description: "have data in io.Reader",
				Matches:     false,
				But:         "it could not be read",
			}
		}
		return matcher(string(all))
	}
}
